package bipf

import (
	"encoding"
	"errors"
	"unsafe"

	"github.com/modern-go/reflect2"
)

type Marshaler interface {
	MarshalBIPF() ([]byte, error)
}

type Unmarshaler interface {
	UnmarshalBIPF([]byte) error
}

var marshalerType = reflect2.TypeOfPtr((*Marshaler)(nil)).Elem()
var unmarshalerType = reflect2.TypeOfPtr((*Unmarshaler)(nil)).Elem()
var binaryMarshalerType = reflect2.TypeOfPtr((*encoding.BinaryMarshaler)(nil)).Elem()
var binaryUnmarshalerType = reflect2.TypeOfPtr((*encoding.BinaryUnmarshaler)(nil)).Elem()

func createDecoderOfMarshaler(typ reflect2.Type) valDecoder {
	ptrType := reflect2.PtrTo(typ)
	if ptrType.Implements(unmarshalerType) {
		return &referenceDecoder{
			&unmarshalerDecoder{ptrType},
		}
	}
	if ptrType.Implements(binaryUnmarshalerType) {
		return &referenceDecoder{
			&binaryUnmarshalerDecoder{ptrType},
		}
	}
	return nil
}

func createEncoderOfMarshaler(ctx *ctx, typ reflect2.Type) (valEncoder, error) {
	if typ == marshalerType {
		checkIsEmpty, err := createCheckIsEmpty(ctx, typ)
		if err != nil {
			return nil, err
		}
		var encoder valEncoder = &directMarshalerEncoder{
			checkIsEmpty: checkIsEmpty,
		}
		return encoder, nil
	}
	if typ.Implements(marshalerType) {
		checkIsEmpty, err := createCheckIsEmpty(ctx, typ)
		if err != nil {
			return nil, err
		}
		var encoder valEncoder = &marshalerEncoder{
			valType:      typ,
			checkIsEmpty: checkIsEmpty,
		}
		return encoder, nil
	}
	ptrType := reflect2.PtrTo(typ)
	if ctx.prefix != "" && ptrType.Implements(marshalerType) {
		checkIsEmpty, err := createCheckIsEmpty(ctx, ptrType)
		if err != nil {
			return nil, err
		}
		var encoder valEncoder = &marshalerEncoder{
			valType:      ptrType,
			checkIsEmpty: checkIsEmpty,
		}
		return &referenceEncoder{encoder}, nil
	}
	if typ == binaryMarshalerType {
		checkIsEmpty, err := createCheckIsEmpty(ctx, typ)
		if err != nil {
			return nil, err
		}
		enc, err := encoderOf(reflect2.TypeOf([]byte{}))
		if err != nil {
			return nil, err
		}
		var encoder valEncoder = &directBinaryMarshalerEncoder{
			checkIsEmpty:  checkIsEmpty,
			stringEncoder: enc,
		}
		return encoder, nil
	}
	if typ.Implements(binaryMarshalerType) {
		checkIsEmpty, err := createCheckIsEmpty(ctx, typ)
		if err != nil {
			return nil, err
		}
		enc, err := encoderOf(reflect2.TypeOf([]byte{}))
		if err != nil {
			return nil, err
		}
		var encoder valEncoder = &binaryMarshalerEncoder{
			valType:      typ,
			bytesEncoder: enc,
			checkIsEmpty: checkIsEmpty,
		}
		return encoder, nil
	}
	// if prefix is empty, the type is the root type
	if ctx.prefix != "" && ptrType.Implements(binaryMarshalerType) {
		checkIsEmpty, err := createCheckIsEmpty(ctx, ptrType)
		if err != nil {
			return nil, err
		}
		enc, err := encoderOf(reflect2.TypeOf([]byte{}))
		if err != nil {
			return nil, err
		}
		var encoder valEncoder = &binaryMarshalerEncoder{
			valType:      ptrType,
			bytesEncoder: enc,
			checkIsEmpty: checkIsEmpty,
		}
		return &referenceEncoder{encoder}, nil
	}
	return nil, errors.New("encoder of marshaler not found")
}

type marshalerEncoder struct {
	checkIsEmpty checkIsEmpty
	valType      reflect2.Type
}

func (encoder *marshalerEncoder) Encode(ptr unsafe.Pointer, stream *stream) error {
	obj := encoder.valType.UnsafeIndirect(ptr)
	if encoder.valType.IsNullable() && reflect2.IsNil(obj) {
		stream.WriteNil()
		return nil
	}
	marshaler := obj.(Marshaler)
	bytes, err := marshaler.MarshalBIPF()
	if err != nil {
		return err
	} else {
		_, err := stream.Write(bytes)
		return err
	}
}

func (encoder *marshalerEncoder) IsEmpty(ptr unsafe.Pointer) (bool, error) {
	return encoder.checkIsEmpty.IsEmpty(ptr)
}

type directMarshalerEncoder struct {
	checkIsEmpty checkIsEmpty
}

func (encoder *directMarshalerEncoder) Encode(ptr unsafe.Pointer, stream *stream) error {
	marshaler := *(*Marshaler)(ptr)
	if marshaler == nil {
		stream.WriteNil()
		return nil
	}
	bytes, err := marshaler.MarshalBIPF()
	if err != nil {
		return err
	} else {
		_, err := stream.Write(bytes)
		return err
	}
}

func (encoder *directMarshalerEncoder) IsEmpty(ptr unsafe.Pointer) (bool, error) {
	return encoder.checkIsEmpty.IsEmpty(ptr)
}

type binaryMarshalerEncoder struct {
	valType      reflect2.Type
	bytesEncoder valEncoder
	checkIsEmpty checkIsEmpty
}

func (encoder *binaryMarshalerEncoder) Encode(ptr unsafe.Pointer, stream *stream) error {
	obj := encoder.valType.UnsafeIndirect(ptr)
	if encoder.valType.IsNullable() && reflect2.IsNil(obj) {
		stream.WriteNil()
		return nil
	}
	marshaler := (obj).(encoding.BinaryMarshaler)
	bytes, err := marshaler.MarshalBinary()
	if err != nil {
		return err
	} else {
		err := encoder.bytesEncoder.Encode(unsafe.Pointer(&bytes), stream)
		return err
	}
}

func (encoder *binaryMarshalerEncoder) IsEmpty(ptr unsafe.Pointer) (bool, error) {
	return encoder.checkIsEmpty.IsEmpty(ptr)
}

type directBinaryMarshalerEncoder struct {
	stringEncoder valEncoder
	checkIsEmpty  checkIsEmpty
}

func (encoder *directBinaryMarshalerEncoder) Encode(ptr unsafe.Pointer, stream *stream) error {
	marshaler := *(*encoding.BinaryMarshaler)(ptr)
	if marshaler == nil {
		stream.WriteNil()
		return nil
	}
	bytes, err := marshaler.MarshalBinary()
	if err != nil {
		return err
	} else {
		err := encoder.stringEncoder.Encode(unsafe.Pointer(&bytes), stream)
		if err != nil {
			return err
		}
	}
	return nil
}

func (encoder *directBinaryMarshalerEncoder) IsEmpty(ptr unsafe.Pointer) (bool, error) {
	return encoder.checkIsEmpty.IsEmpty(ptr)
}

type unmarshalerDecoder struct {
	valType reflect2.Type
}

func (decoder *unmarshalerDecoder) Decode(ptr unsafe.Pointer, iter *iterator) error {
	valType := decoder.valType
	obj := valType.UnsafeIndirect(ptr)
	unmarshaler := obj.(Unmarshaler)

	_, l, err := iter.readTag()
	if err != nil {
		return err
	}
	buf := make([]byte, l)
	_, err = iter.Read(buf)
	if err != nil {
		return err
	}

	err = unmarshaler.UnmarshalBIPF(buf)
	if err != nil {
		return err
	}

	return nil
}

type binaryUnmarshalerDecoder struct {
	valType reflect2.Type
}

func (decoder *binaryUnmarshalerDecoder) Decode(ptr unsafe.Pointer, iter *iterator) error {
	valType := decoder.valType
	obj := valType.UnsafeIndirect(ptr)
	if reflect2.IsNil(obj) {
		ptrType := valType.(*reflect2.UnsafePtrType)
		elemType := ptrType.Elem()
		elem := elemType.UnsafeNew()
		ptrType.UnsafeSet(ptr, unsafe.Pointer(&elem))
		obj = valType.UnsafeIndirect(ptr)
	}
	unmarshaler := (obj).(encoding.BinaryUnmarshaler)
	bytes, err := iter.ReadBuffer()
	if err != nil {
		return err
	}
	err = unmarshaler.UnmarshalBinary(bytes)
	if err != nil {
		return err
	}
	return nil
}
