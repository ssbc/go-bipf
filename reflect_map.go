package bipf

import (
	"errors"
	"io"
	"reflect"
	"unsafe"

	"github.com/modern-go/reflect2"
)

func decoderOfMap(ctx *ctx, typ reflect2.Type) (valDecoder, error) {
	mapType := typ.(*reflect2.UnsafeMapType)
	keyDecoder, err := decoderOfMapKey(ctx.append("[mapKey]"), mapType.Key())
	if err != nil {
		return nil, err
	}
	elemDecoder, err := decoderOfType(ctx.append("[mapElem]"), mapType.Elem())
	if err != nil {
		return nil, err
	}
	return &mapDecoder{
		mapType:     mapType,
		keyType:     mapType.Key(),
		elemType:    mapType.Elem(),
		keyDecoder:  keyDecoder,
		elemDecoder: elemDecoder,
	}, nil
}

func encoderOfMap(ctx *ctx, typ reflect2.Type) (valEncoder, error) {
	mapType := typ.(*reflect2.UnsafeMapType)
	keyEncoder, err := encoderOfMapKey(ctx.append("[mapKey]"), mapType.Key())
	if err != nil {
		return nil, err
	}
	elemEncoder, err := encoderOfType(ctx.append("[mapElem]"), mapType.Elem())
	if err != nil {
		return nil, err
	}
	return &mapEncoder{
		mapType:     mapType,
		keyEncoder:  keyEncoder,
		elemEncoder: elemEncoder,
	}, nil
}

func decoderOfMapKey(ctx *ctx, typ reflect2.Type) (valDecoder, error) {
	ptrType := reflect2.PtrTo(typ)
	if ptrType.Implements(unmarshalerType) {
		return &referenceDecoder{
			&unmarshalerDecoder{
				valType: ptrType,
			},
		}, nil
	}
	if typ.Implements(unmarshalerType) {
		return &unmarshalerDecoder{
			valType: typ,
		}, nil
	}
	if ptrType.Implements(binaryUnmarshalerType) {
		return &referenceDecoder{
			&binaryUnmarshalerDecoder{
				valType: ptrType,
			},
		}, nil
	}
	if typ.Implements(binaryUnmarshalerType) {
		return &binaryUnmarshalerDecoder{
			valType: typ,
		}, nil
	}

	switch typ.Kind() {
	case reflect.String, reflect.Bool,
		reflect.Uint8, reflect.Int8,
		reflect.Uint16, reflect.Int16,
		reflect.Uint32, reflect.Int32,
		reflect.Uint64, reflect.Int64,
		reflect.Uint, reflect.Int,
		reflect.Float32, reflect.Float64,
		reflect.Uintptr:
		return decoderOfType(ctx, reflect2.DefaultTypeOfKind(typ.Kind()))
	default:
		return nil, errors.New("decoder of map key not found")
	}
}

func encoderOfMapKey(ctx *ctx, typ reflect2.Type) (valEncoder, error) {
	if typ.Kind() != reflect.String {
		if typ == binaryMarshalerType {
			enc, err := encoderOf(reflect2.TypeOf([]byte{}))
			if err != nil {
				return nil, err
			}
			return &directBinaryMarshalerEncoder{
				stringEncoder: enc,
			}, nil
		}
		if typ.Implements(binaryMarshalerType) {
			enc, err := encoderOf(reflect2.TypeOf([]byte{}))
			if err != nil {
				return nil, err
			}
			return &binaryMarshalerEncoder{
				valType:      typ,
				bytesEncoder: enc,
			}, nil
		}
	}

	switch typ.Kind() {
	case reflect.String, reflect.Bool,
		reflect.Uint8, reflect.Int8,
		reflect.Uint16, reflect.Int16,
		reflect.Uint32, reflect.Int32,
		reflect.Uint64, reflect.Int64,
		reflect.Uint, reflect.Int,
		reflect.Float32, reflect.Float64,
		reflect.Uintptr:
		typ = reflect2.DefaultTypeOfKind(typ.Kind())
		return encoderOfType(ctx, typ)
	default:
		if typ.Kind() == reflect.Interface {
			return &dynamicMapKeyEncoder{ctx, typ}, nil
		}
		return nil, errors.New("encoder of map key not found")
	}
}

type mapDecoder struct {
	mapType     *reflect2.UnsafeMapType
	keyType     reflect2.Type
	elemType    reflect2.Type
	keyDecoder  valDecoder
	elemDecoder valDecoder
}

func (decoder *mapDecoder) Decode(ptr unsafe.Pointer, iter *iterator) error {
	mapType := decoder.mapType
	ok, err := iter.CheckNilIsNext()
	if err != nil {
		return err
	}
	if ok {
		mapType.UnsafeSet(ptr, mapType.UnsafeNew())
		return nil
	}

	typ, l, err := iter.readTag()
	if err != nil {
		return err
	}

	if typ != valueTypeObject {
		return errors.New("invalid type")
	}

	if mapType.UnsafeIsNil(ptr) {
		mapType.UnsafeSet(ptr, mapType.UnsafeMakeMap(0))
	}

	start := iter.numRead()

	for iter.numRead()-start < l {
		key := decoder.keyType.UnsafeNew()
		err := decoder.keyDecoder.Decode(key, iter)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}

		elem := decoder.elemType.UnsafeNew()
		err = decoder.elemDecoder.Decode(elem, iter)
		if err != nil {
			return err
		}

		decoder.mapType.UnsafeSetIndex(ptr, key, elem)

		if iter.numRead()-start > l {
			return errors.New("out of bounds")
		}
	}

	return nil
}

type dynamicMapKeyEncoder struct {
	ctx     *ctx
	valType reflect2.Type
}

func (encoder *dynamicMapKeyEncoder) Encode(ptr unsafe.Pointer, stream *stream) error {
	obj := encoder.valType.UnsafeIndirect(ptr)
	enc, err := encoderOfMapKey(encoder.ctx, reflect2.TypeOf(obj))
	if err != nil {
		return err
	}
	return enc.Encode(reflect2.PtrOf(obj), stream)
}

func (encoder *dynamicMapKeyEncoder) IsEmpty(ptr unsafe.Pointer) (bool, error) {
	obj := encoder.valType.UnsafeIndirect(ptr)
	enc, err := encoderOfMapKey(encoder.ctx, reflect2.TypeOf(obj))
	if err != nil {
		return false, err
	}
	return enc.IsEmpty(reflect2.PtrOf(obj))
}

type mapEncoder struct {
	mapType     *reflect2.UnsafeMapType
	keyEncoder  valEncoder
	elemEncoder valEncoder
}

func (encoder *mapEncoder) Encode(ptr unsafe.Pointer, stream *stream) error {
	if *(*unsafe.Pointer)(ptr) == nil {
		stream.WriteNil()
		return nil
	}

	tmpStream := streamPool.BorrowStream(nil)
	defer streamPool.ReturnStream(tmpStream)

	iter := encoder.mapType.UnsafeIterate(ptr)
	for i := 0; iter.HasNext(); i++ {
		key, elem := iter.UnsafeNext()
		err := encoder.keyEncoder.Encode(key, tmpStream)
		if err != nil {
			return err
		}

		err = encoder.elemEncoder.Encode(elem, tmpStream)
		if err != nil {
			return err
		}
	}

	stream.WriteTag(uint64(tmpStream.Buffered()), valueTypeObject)
	if _, err := stream.Write(tmpStream.Buffer()); err != nil {
		return err
	}

	return nil
}

func (encoder *mapEncoder) IsEmpty(ptr unsafe.Pointer) (bool, error) {
	iter := encoder.mapType.UnsafeIterate(ptr)
	return !iter.HasNext(), nil
}
