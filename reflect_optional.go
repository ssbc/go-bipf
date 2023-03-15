package bipf

import (
	"github.com/modern-go/reflect2"
	"unsafe"
)

func decoderOfOptional(ctx *ctx, typ reflect2.Type) (valDecoder, error) {
	ptrType := typ.(*reflect2.UnsafePtrType)
	elemType := ptrType.Elem()
	decoder, err := decoderOfType(ctx, elemType)
	if err != nil {
		return nil, err
	}
	return &optionalDecoder{elemType, decoder}, nil
}

func encoderOfOptional(ctx *ctx, typ reflect2.Type) (valEncoder, error) {
	ptrType := typ.(*reflect2.UnsafePtrType)
	elemType := ptrType.Elem()
	elemEncoder, err := encoderOfType(ctx, elemType)
	if err != nil {
		return nil, err
	}
	encoder := &optionalEncoder{elemEncoder}
	return encoder, nil
}

type optionalDecoder struct {
	ValueType    reflect2.Type
	ValueDecoder valDecoder
}

func (decoder *optionalDecoder) Decode(ptr unsafe.Pointer, iter *iterator) error {
	ok, err := iter.CheckNilIsNext()
	if err != nil {
		return err
	}

	if ok {
		*((*unsafe.Pointer)(ptr)) = nil
	} else {
		if *((*unsafe.Pointer)(ptr)) == nil {
			//pointer to null, we have to allocate memory to hold the value
			newPtr := decoder.ValueType.UnsafeNew()
			err := decoder.ValueDecoder.Decode(newPtr, iter)
			if err != nil {
				return err
			}
			*((*unsafe.Pointer)(ptr)) = newPtr
		} else {
			//reuse existing instance
			err := decoder.ValueDecoder.Decode(*((*unsafe.Pointer)(ptr)), iter)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type dereferenceDecoder struct {
	// only to deference a pointer
	valueType    reflect2.Type
	valueDecoder valDecoder
}

func (decoder *dereferenceDecoder) Decode(ptr unsafe.Pointer, iter *iterator) error {
	if *((*unsafe.Pointer)(ptr)) == nil {
		//pointer to null, we have to allocate memory to hold the value
		newPtr := decoder.valueType.UnsafeNew()
		if err := decoder.valueDecoder.Decode(newPtr, iter); err != nil {
			return err
		}
		*((*unsafe.Pointer)(ptr)) = newPtr
	} else {
		//reuse existing instance
		if err := decoder.valueDecoder.Decode(*((*unsafe.Pointer)(ptr)), iter); err != nil {
			return err
		}
	}

	return nil
}

type optionalEncoder struct {
	ValueEncoder valEncoder
}

func (encoder *optionalEncoder) Encode(ptr unsafe.Pointer, stream *stream) error {
	if *((*unsafe.Pointer)(ptr)) == nil {
		stream.WriteNil()
		return nil
	} else {
		return encoder.ValueEncoder.Encode(*((*unsafe.Pointer)(ptr)), stream)
	}
}

func (encoder *optionalEncoder) IsEmpty(ptr unsafe.Pointer) (bool, error) {
	return *((*unsafe.Pointer)(ptr)) == nil, nil
}

type dereferenceEncoder struct {
	ValueEncoder valEncoder
}

func (encoder *dereferenceEncoder) Encode(ptr unsafe.Pointer, stream *stream) error {
	if *((*unsafe.Pointer)(ptr)) == nil {
		stream.WriteNil()
		return nil
	} else {
		return encoder.ValueEncoder.Encode(*((*unsafe.Pointer)(ptr)), stream)
	}
}

func (encoder *dereferenceEncoder) IsEmpty(ptr unsafe.Pointer) (bool, error) {
	dePtr := *((*unsafe.Pointer)(ptr))
	if dePtr == nil {
		return true, nil
	}
	return encoder.ValueEncoder.IsEmpty(dePtr)
}

func (encoder *dereferenceEncoder) IsEmbeddedPtrNil(ptr unsafe.Pointer) bool {
	deReferenced := *((*unsafe.Pointer)(ptr))
	if deReferenced == nil {
		return true
	}
	isEmbeddedPtrNil, converted := encoder.ValueEncoder.(isEmbeddedPtrNil)
	if !converted {
		return false
	}
	fieldPtr := unsafe.Pointer(deReferenced)
	return isEmbeddedPtrNil.IsEmbeddedPtrNil(fieldPtr)
}

type referenceEncoder struct {
	encoder valEncoder
}

func (encoder *referenceEncoder) Encode(ptr unsafe.Pointer, stream *stream) error {
	return encoder.encoder.Encode(unsafe.Pointer(&ptr), stream)
}

func (encoder *referenceEncoder) IsEmpty(ptr unsafe.Pointer) (bool, error) {
	return encoder.encoder.IsEmpty(unsafe.Pointer(&ptr))
}

type referenceDecoder struct {
	decoder valDecoder
}

func (decoder *referenceDecoder) Decode(ptr unsafe.Pointer, iter *iterator) error {
	return decoder.decoder.Decode(unsafe.Pointer(&ptr), iter)
}
