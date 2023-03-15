package bipf

import (
	"errors"
	"reflect"
	"unsafe"

	"github.com/modern-go/reflect2"
)

type dynamicEncoder struct {
	valType reflect2.Type
}

func (encoder *dynamicEncoder) Encode(ptr unsafe.Pointer, stream *stream) error {
	obj := encoder.valType.UnsafeIndirect(ptr)
	return stream.WriteVal(obj)
}

func (encoder *dynamicEncoder) IsEmpty(ptr unsafe.Pointer) (bool, error) {
	return encoder.valType.UnsafeIndirect(ptr) == nil, nil
}

type efaceDecoder struct {
}

func (decoder *efaceDecoder) Decode(ptr unsafe.Pointer, iter *iterator) error {
	pObj := (*any)(ptr)
	obj := *pObj
	if obj == nil {
		var err error
		*pObj, err = iter.ReadAny()
		if err != nil {
			return err
		}
		return nil
	}
	typ := reflect2.TypeOf(obj)
	if typ.Kind() != reflect.Ptr {
		var err error
		*pObj, err = iter.ReadAny()
		if err != nil {
			return err
		}
		return nil
	}
	ptrType := typ.(*reflect2.UnsafePtrType)
	ptrElemType := ptrType.Elem()
	nilIsNext, err := iter.CheckNilIsNext()
	if err != nil {
		return err
	}
	if nilIsNext {
		if ptrElemType.Kind() != reflect.Ptr {
			if err := iter.ReadNil(); err != nil {
				return err
			}
			*pObj = nil
			return nil
		}
	}
	if reflect2.IsNil(obj) {
		obj := ptrElemType.New()
		err := iter.ReadVal(obj)
		if err != nil {
			return err
		}
		*pObj = obj
		return nil
	}
	return iter.ReadVal(obj)
}

type ifaceDecoder struct {
	valType *reflect2.UnsafeIFaceType
}

func (decoder *ifaceDecoder) Decode(ptr unsafe.Pointer, iter *iterator) error {
	nilIsNext, err := iter.CheckNilIsNext()
	if err != nil {
		return err
	}
	if nilIsNext {
		if err := iter.ReadNil(); err != nil {
			return err
		}
		decoder.valType.UnsafeSet(ptr, decoder.valType.UnsafeNew())
		return nil
	}
	obj := decoder.valType.UnsafeIndirect(ptr)
	if reflect2.IsNil(obj) {
		return errors.New("decode non empty interface: can not unmarshal into nil")
	}
	return iter.ReadVal(obj)
}
