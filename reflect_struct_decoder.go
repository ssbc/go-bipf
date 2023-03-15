package bipf

import (
	"errors"
	"strings"
	"unsafe"

	"github.com/modern-go/reflect2"
)

func decoderOfStruct(ctx *ctx, typ reflect2.Type) (valDecoder, error) {
	bindings := map[string]*binding{}
	structDescriptor, err := describeStruct(ctx, typ)
	if err != nil {
		return nil, err
	}
	for _, binding := range structDescriptor.Fields {
		for _, fromName := range binding.FromNames {
			old := bindings[fromName]
			if old == nil {
				bindings[fromName] = binding
				continue
			}
			ignoreOld, ignoreNew := resolveConflictBinding(old, binding)
			if ignoreOld {
				delete(bindings, fromName)
			}
			if !ignoreNew {
				bindings[fromName] = binding
			}
		}
	}
	fields := map[string]*structFieldDecoder{}
	for k, binding := range bindings {
		fields[k] = binding.Decoder.(*structFieldDecoder)
	}

	for k, binding := range bindings {
		if _, found := fields[strings.ToLower(k)]; !found {
			fields[strings.ToLower(k)] = binding.Decoder.(*structFieldDecoder)
		}
	}

	return &generalStructDecoder{typ, fields, false}, nil
}

type generalStructDecoder struct {
	typ                   reflect2.Type
	fields                map[string]*structFieldDecoder
	disallowUnknownFields bool
}

func (decoder *generalStructDecoder) Decode(ptr unsafe.Pointer, iter *iterator) error {
	typ, l, err := iter.readTag()
	if err != nil {
		return err
	}

	if typ == valueTypeBoolNull && l == 0 {
		return nil
	}

	if typ != valueTypeObject {
		return errors.New("unexpected type")
	}

	start := iter.numRead()

	if err := iter.incrementDepth(); err != nil {
		return err
	}

	for iter.numRead()-start < l {
		if err := decoder.decodeOneField(ptr, iter); err != nil {
			return err
		}
	}

	if err := iter.decrementDepth(); err != nil {
		return err
	}

	return nil
}

func (decoder *generalStructDecoder) decodeOneField(ptr unsafe.Pointer, iter *iterator) error {
	var fieldDecoder *structFieldDecoder

	field, err := iter.ReadString()
	if err != nil {
		return err
	}

	fieldDecoder = decoder.fields[field]
	if fieldDecoder == nil {
		fieldDecoder = decoder.fields[strings.ToLower(field)]
	}

	if fieldDecoder == nil {
		if err := iter.skip(); err != nil {
			return err
		}
		return nil
	}

	return fieldDecoder.Decode(ptr, iter)
}

type structFieldDecoder struct {
	field        reflect2.StructField
	fieldDecoder valDecoder
}

func (decoder *structFieldDecoder) Decode(ptr unsafe.Pointer, iter *iterator) error {
	fieldPtr := decoder.field.UnsafeGet(ptr)
	if err := decoder.fieldDecoder.Decode(fieldPtr, iter); err != nil {
		return wrap(err, decoder.field.Name())
	}
	return nil
}
