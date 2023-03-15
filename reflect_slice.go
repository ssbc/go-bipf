package bipf

import (
	"errors"
	"fmt"
	"github.com/modern-go/reflect2"
	"unsafe"
)

func decoderOfSlice(ctx *ctx, typ reflect2.Type) (valDecoder, error) {
	sliceType := typ.(*reflect2.UnsafeSliceType)
	decoder, err := decoderOfType(ctx.append("[sliceElem]"), sliceType.Elem())
	if err != nil {
		return nil, err
	}
	return &sliceDecoder{sliceType, decoder}, nil
}

func encoderOfSlice(ctx *ctx, typ reflect2.Type) (valEncoder, error) {
	sliceType := typ.(*reflect2.UnsafeSliceType)
	encoder, err := encoderOfType(ctx.append("[sliceElem]"), sliceType.Elem())
	if err != nil {
		return nil, err
	}
	return &sliceEncoder{sliceType, encoder}, nil
}

type sliceEncoder struct {
	sliceType   *reflect2.UnsafeSliceType
	elemEncoder valEncoder
}

func (encoder *sliceEncoder) Encode(ptr unsafe.Pointer, stream *stream) error {
	if encoder.sliceType.UnsafeIsNil(ptr) {
		stream.WriteNil()
		return nil
	}
	length := encoder.sliceType.UnsafeLengthOf(ptr)
	if length == 0 {
		stream.WriteEmptyArray()
		return nil
	}

	tmpStream := streamPool.BorrowStream(nil)
	defer streamPool.ReturnStream(tmpStream)

	for i := 0; i < length; i++ {
		elemPtr := encoder.sliceType.UnsafeGetIndex(ptr, i)
		err := encoder.elemEncoder.Encode(elemPtr, tmpStream)
		if err != nil {
			return wrapf(err, "type '%v'", encoder.sliceType)
		}
	}

	stream.WriteTag(uint64(tmpStream.Buffered()), valueTypeArray)
	if _, err := stream.Write(tmpStream.Buffer()); err != nil {
		return err
	}

	return nil
}

func (encoder *sliceEncoder) IsEmpty(ptr unsafe.Pointer) (bool, error) {
	return encoder.sliceType.UnsafeLengthOf(ptr) == 0, nil
}

type sliceDecoder struct {
	sliceType   *reflect2.UnsafeSliceType
	elemDecoder valDecoder
}

func (decoder *sliceDecoder) Decode(ptr unsafe.Pointer, iter *iterator) error {
	typ, l, err := iter.readTag()
	if err != nil {
		return nil
	}

	sliceType := decoder.sliceType

	if typ == valueTypeBoolNull && l == 0 {
		sliceType.UnsafeSetNil(ptr)
		return nil
	}

	if typ != valueTypeArray {
		return fmt.Errorf("incorrect type: %v", typ)
	}

	start := iter.numRead()
	i := 0

	for iter.numRead()-start < l {
		sliceType.UnsafeGrow(ptr, i+1)

		elemPtr := sliceType.UnsafeGetIndex(ptr, i)
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
