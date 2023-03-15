package bipf

import (
	"errors"
	"reflect"
	"unsafe"

	"github.com/modern-go/reflect2"
)

func encoderOfStruct(ctx *ctx, typ reflect2.Type) (valEncoder, error) {
	type bindingTo struct {
		binding *binding
		toName  string
		ignored bool
	}
	var orderedBindings []*bindingTo
	structDescriptor, err := describeStruct(ctx, typ)
	if err != nil {
		return nil, err
	}
	for _, binding := range structDescriptor.Fields {
		for _, toName := range binding.ToNames {
			newBinding := &bindingTo{
				binding: binding,
				toName:  toName,
			}
			for _, oldBinding := range orderedBindings {
				if oldBinding.toName != toName {
					continue
				}
				oldBinding.ignored, newBinding.ignored = resolveConflictBinding(oldBinding.binding, newBinding.binding)
			}
			orderedBindings = append(orderedBindings, newBinding)
		}
	}
	if len(orderedBindings) == 0 {
		return &emptyStructEncoder{}, nil
	}
	var finalOrderedFields []structFieldTo
	for _, bindingTo := range orderedBindings {
		if !bindingTo.ignored {
			finalOrderedFields = append(finalOrderedFields, structFieldTo{
				encoder: bindingTo.binding.Encoder.(*structFieldEncoder),
				toName:  bindingTo.toName,
			})
		}
	}
	return &structEncoder{typ, finalOrderedFields}, nil
}

func createCheckIsEmpty(ctx *ctx, typ reflect2.Type) (checkIsEmpty, error) {
	encoder, err := createEncoderOfNative(ctx, typ)
	if err == nil {
		return encoder, nil
	}
	kind := typ.Kind()
	switch kind {
	case reflect.Interface:
		return &dynamicEncoder{typ}, nil
	case reflect.Struct:
		return &structEncoder{typ: typ}, nil
	case reflect.Array:
		return &arrayEncoder{}, nil
	case reflect.Slice:
		return &sliceEncoder{}, nil
	case reflect.Map:
		return encoderOfMap(ctx, typ)
	case reflect.Ptr:
		return &optionalEncoder{}, nil
	default:
		return nil, errors.New("could not find an empty check for this type")
	}
}

func resolveConflictBinding(old, new *binding) (ignoreOld, ignoreNew bool) {
	newTagged := new.Field.Tag().Get(tagKey) != ""
	oldTagged := old.Field.Tag().Get(tagKey) != ""
	if newTagged {
		if oldTagged {
			if len(old.levels) > len(new.levels) {
				return true, false
			} else if len(new.levels) > len(old.levels) {
				return false, true
			} else {
				return true, true
			}
		} else {
			return true, false
		}
	} else {
		if oldTagged {
			return true, false
		}
		if len(old.levels) > len(new.levels) {
			return true, false
		} else if len(new.levels) > len(old.levels) {
			return false, true
		} else {
			return true, true
		}
	}
}

type structFieldEncoder struct {
	field        reflect2.StructField
	fieldEncoder valEncoder
	omitempty    bool
}

func (encoder *structFieldEncoder) Encode(ptr unsafe.Pointer, stream *stream) error {
	fieldPtr := encoder.field.UnsafeGet(ptr)
	err := encoder.fieldEncoder.Encode(fieldPtr, stream)
	if err != nil {
		return wrapf(err, "field name '%s'", encoder.field.Name())
	}
	return nil
}

func (encoder *structFieldEncoder) IsEmpty(ptr unsafe.Pointer) (bool, error) {
	fieldPtr := encoder.field.UnsafeGet(ptr)
	return encoder.fieldEncoder.IsEmpty(fieldPtr)
}

func (encoder *structFieldEncoder) IsEmbeddedPtrNil(ptr unsafe.Pointer) bool {
	isEmbeddedPtrNil, converted := encoder.fieldEncoder.(isEmbeddedPtrNil)
	if !converted {
		return false
	}
	fieldPtr := encoder.field.UnsafeGet(ptr)
	return isEmbeddedPtrNil.IsEmbeddedPtrNil(fieldPtr)
}

type isEmbeddedPtrNil interface {
	IsEmbeddedPtrNil(ptr unsafe.Pointer) bool
}

type structEncoder struct {
	typ    reflect2.Type
	fields []structFieldTo
}

type structFieldTo struct {
	encoder *structFieldEncoder
	toName  string
}

func (encoder *structEncoder) Encode(ptr unsafe.Pointer, stream *stream) error {
	tmpStream := streamPool.BorrowStream(nil)
	defer streamPool.ReturnStream(tmpStream)

	for _, field := range encoder.fields {
		isEmpty, err := field.encoder.IsEmpty(ptr)
		if err != nil {
			return err
		}
		if field.encoder.omitempty && isEmpty {
			continue
		}
		if field.encoder.IsEmbeddedPtrNil(ptr) {
			continue
		}

		err = tmpStream.WriteString(field.toName)
		if err != nil {
			return err
		}

		err = field.encoder.Encode(ptr, tmpStream)
		if err != nil {
			return err
		}
	}

	stream.WriteTag(uint64(tmpStream.Buffered()), valueTypeObject)
	_, err := stream.Write(tmpStream.Buffer())
	if err != nil {
		return err
	}

	return nil
}

func (encoder *structEncoder) IsEmpty(ptr unsafe.Pointer) (bool, error) {
	return false, nil
}

type emptyStructEncoder struct {
}

func (encoder *emptyStructEncoder) Encode(ptr unsafe.Pointer, stream *stream) error {
	stream.WriteEmptyObject()
	return nil
}

func (encoder *emptyStructEncoder) IsEmpty(ptr unsafe.Pointer) (bool, error) {
	return false, nil
}
