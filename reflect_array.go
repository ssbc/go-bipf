package bipf

import (
	"errors"
	"fmt"
	"github.com/modern-go/reflect2"
	"unsafe"
)

func decoderOfArray(ctx *ctx, typ reflect2.Type) (valDecoder, error) {
	arrayType := typ.(*reflect2.UnsafeArrayType)
	decoder, err := decoderOfType(ctx.append("[arrayElem]"), arrayType.Elem())
	if err != nil {
		return nil, err
	}
	return &arrayDecoder{arrayType, decoder}, nil
}

func encoderOfArray(ctx *ctx, typ reflect2.Type) (valEncoder, error) {
	arrayType := typ.(*reflect2.UnsafeArrayType)
	if arrayType.Len() == 0 {
		return emptyArrayEncoder{}, nil
	}
	encoder, err := encoderOfType(ctx.append("[arrayElem]"), arrayType.Elem())
	if err != nil {
		return nil, err
	}
	return &arrayEncoder{arrayType, encoder}, nil
}

type emptyArrayEncoder struct{}

func (encoder emptyArrayEncoder) Encode(ptr unsafe.Pointer, stream *stream) error {
	stream.WriteEmptyArray()
	return nil
}

func (encoder emptyArrayEncoder) IsEmpty(ptr unsafe.Pointer) (bool, error) {
	return true, nil
}

type arrayEncoder struct {
	arrayType   *reflect2.UnsafeArrayType
	elemEncoder valEncoder
}

func (encoder *arrayEncoder) Encode(ptr unsafe.Pointer, stream *stream) error {
	tmpStream := streamPool.BorrowStream(nil)
	defer streamPool.ReturnStream(tmpStream)

	elemPtr := unsafe.Pointer(ptr)
	for i := 0; i < encoder.arrayType.Len(); i++ {
		elemPtr = encoder.arrayType.UnsafeGetIndex(ptr, i)
		err := encoder.elemEncoder.Encode(elemPtr, tmpStream)
		if err != nil {
			return wrapf(err, "type '%v'", encoder.arrayType)
		}
	}

	stream.WriteTag(uint64(tmpStream.Buffered()), valueTypeArray)
	if _, err := stream.Write(tmpStream.Buffer()); err != nil {
		return err
	}

	return nil
}

func (encoder *arrayEncoder) IsEmpty(ptr unsafe.Pointer) (bool, error) {
	return false, nil
}

type arrayDecoder struct {
	arrayType   *reflect2.UnsafeArrayType
	elemDecoder valDecoder
}

func (decoder *arrayDecoder) Decode(ptr unsafe.Pointer, iter *iterator) error {
	typ, l, err := iter.readTag()
	if err != nil {
		return nil
	}

	arrayType := decoder.arrayType

	if typ == valueTypeBoolNull && l == 0 {
		return nil
	}

	if typ != valueTypeArray {
		return fmt.Errorf("incorrect type: %v", typ)
	}

	start := iter.numRead()
	i := 0

	for iter.numRead()-start < l {
		if i >= arrayType.Len() {
			return errors.New("provided array is too short")
		}
		elemPtr := arrayType.UnsafeGetIndex(ptr, i)
		err := decoder.elemDecoder.Decode(elemPtr, iter)
		if err != nil {
			return err
		}

		if iter.numRead()-start > l {
			return errors.New("out of bounds")
		}

		i++
	}

	return nil
}
