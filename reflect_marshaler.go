package bipf

import (
	"encoding"
	"encoding/json"
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
var unmarshalerType = reflect2.TypeOfPtr((*json.Unmarshaler)(nil)).Elem()
var textMarshalerType = reflect2.TypeOfPtr((*encoding.TextMarshaler)(nil)).Elem()
var textUnmarshalerType = reflect2.TypeOfPtr((*encoding.TextUnmarshaler)(nil)).Elem()

func createDecoderOfMarshaler(typ reflect2.Type) valDecoder {
	ptrType := reflect2.PtrTo(typ)
	if ptrType.Implements(unmarshalerType) {
		return &referenceDecoder{
			&unmarshalerDecoder{ptrType},
		}
	}
	if ptrType.Implements(textUnmarshalerType) {
		return &referenceDecoder{
			&textUnmarshalerDecoder{ptrType},
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
	if typ == textMarshalerType {
		checkIsEmpty, err := createCheckIsEmpty(ctx, typ)
		if err != nil {
			return nil, err
		}
		enc, err := encoderOf(reflect2.TypeOf(""))
		if err != nil {
			return nil, err
		}
		var encoder valEncoder = &directTextMarshalerEncoder{
			checkIsEmpty:  checkIsEmpty,
			stringEncoder: enc,
		}
		return encoder, nil
	}
	if typ.Implements(textMarshalerType) {
		checkIsEmpty, err := createCheckIsEmpty(ctx, typ)
		if err != nil {
			return nil, err
		}
		enc, err := encoderOf(reflect2.TypeOf(""))
		if err != nil {
			return nil, err
		}
		var encoder valEncoder = &textMarshalerEncoder{
			valType:       typ,
			stringEncoder: enc,
			checkIsEmpty:  checkIsEmpty,
		}
		return encoder, nil
	}
	// if prefix is empty, the type is the root type
	if ctx.prefix != "" && ptrType.Implements(textMarshalerType) {
		checkIsEmpty, err := createCheckIsEmpty(ctx, ptrType)
		if err != nil {
			return nil, err
		}
		enc, err := encoderOf(reflect2.TypeOf(""))
		if err != nil {
			return nil, err
		}
		var encoder valEncoder = &textMarshalerEncoder{
			valType:       ptrType,
			stringEncoder: enc,
			checkIsEmpty:  checkIsEmpty,
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

type textMarshalerEncoder struct {
	valType       reflect2.Type
	stringEncoder valEncoder
	checkIsEmpty  checkIsEmpty
}

func (encoder *textMarshalerEncoder) Encode(ptr unsafe.Pointer, stream *stream) error {
	obj := encoder.valType.UnsafeIndirect(ptr)
	if encoder.valType.IsNullable() && reflect2.IsNil(obj) {
		stream.WriteNil()
		return nil
	}
	marshaler := (obj).(encoding.TextMarshaler)
	bytes, err := marshaler.MarshalText()
	if err != nil {
		return err
	} else {
		str := string(bytes)
		err := encoder.stringEncoder.Encode(unsafe.Pointer(&str), stream)
		return err
	}
}

func (encoder *textMarshalerEncoder) IsEmpty(ptr unsafe.Pointer) (bool, error) {
	return encoder.checkIsEmpty.IsEmpty(ptr)
}

type directTextMarshalerEncoder struct {
	stringEncoder valEncoder
	checkIsEmpty  checkIsEmpty
}

func (encoder *directTextMarshalerEncoder) Encode(ptr unsafe.Pointer, stream *stream) error {
	marshaler := *(*encoding.TextMarshaler)(ptr)
	if marshaler == nil {
		stream.WriteNil()
		return nil
	}
	bytes, err := marshaler.MarshalText()
	if err != nil {
		return err
	} else {
		str := string(bytes)
		err := encoder.stringEncoder.Encode(unsafe.Pointer(&str), stream)
		if err != nil {
			return err
		}
	}
	return nil
}

func (encoder *directTextMarshalerEncoder) IsEmpty(ptr unsafe.Pointer) (bool, error) {
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

type textUnmarshalerDecoder struct {
	valType reflect2.Type
}

func (decoder *textUnmarshalerDecoder) Decode(ptr unsafe.Pointer, iter *iterator) error {
	valType := decoder.valType
	obj := valType.UnsafeIndirect(ptr)
	if reflect2.IsNil(obj) {
		ptrType := valType.(*reflect2.UnsafePtrType)
		elemType := ptrType.Elem()
		elem := elemType.UnsafeNew()
		ptrType.UnsafeSet(ptr, unsafe.Pointer(&elem))
		obj = valType.UnsafeIndirect(ptr)
	}
	unmarshaler := (obj).(encoding.TextUnmarshaler)
	str, err := iter.ReadString()
	if err != nil {
		return err
	}
	err = unmarshaler.UnmarshalText([]byte(str))
	if err != nil {
		return err
	}
	return nil
}
