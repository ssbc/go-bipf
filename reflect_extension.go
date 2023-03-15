package bipf

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"unicode"

	"github.com/modern-go/reflect2"
)

var fieldDecoders = map[string]valDecoder{}
var fieldEncoders = map[string]valEncoder{}

type structDescriptor struct {
	Type   reflect2.Type
	Fields []*binding
}

// binding describe how should we encode/decode the struct field
type binding struct {
	levels    []int
	Field     reflect2.StructField
	FromNames []string
	ToNames   []string
	Encoder   valEncoder
	Decoder   valDecoder
}

func describeStruct(ctx *ctx, typ reflect2.Type) (*structDescriptor, error) {
	structType := typ.(*reflect2.UnsafeStructType)
	embeddedBindings := []*binding{}
	bindings := []*binding{}
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		tag, hastag := field.Tag().Lookup(tagKey)
		if hastag && (tag == "-" || field.Name() == "_") {
			continue
		}
		tagParts := strings.Split(tag, ",")
		if field.Anonymous() && (tag == "" || tagParts[0] == "") {
			if field.Type().Kind() == reflect.Struct {
				structDescriptor, err := describeStruct(ctx, field.Type())
				if err != nil {
					return nil, err
				}
				for _, binding := range structDescriptor.Fields {
					binding.levels = append([]int{i}, binding.levels...)
					omitempty := binding.Encoder.(*structFieldEncoder).omitempty
					binding.Encoder = &structFieldEncoder{field, binding.Encoder, omitempty}
					binding.Decoder = &structFieldDecoder{field, binding.Decoder}
					embeddedBindings = append(embeddedBindings, binding)
				}
				continue
			} else if field.Type().Kind() == reflect.Ptr {
				ptrType := field.Type().(*reflect2.UnsafePtrType)
				if ptrType.Elem().Kind() == reflect.Struct {
					structDescriptor, err := describeStruct(ctx, ptrType.Elem())
					if err != nil {
						return nil, err
					}
					for _, binding := range structDescriptor.Fields {
						binding.levels = append([]int{i}, binding.levels...)
						omitempty := binding.Encoder.(*structFieldEncoder).omitempty
						binding.Encoder = &dereferenceEncoder{binding.Encoder}
						binding.Encoder = &structFieldEncoder{field, binding.Encoder, omitempty}
						binding.Decoder = &dereferenceDecoder{ptrType.Elem(), binding.Decoder}
						binding.Decoder = &structFieldDecoder{field, binding.Decoder}
						embeddedBindings = append(embeddedBindings, binding)
					}
					continue
				}
			}
		}
		fieldNames := calcFieldNames(field.Name(), tagParts[0], tag)
		fieldCacheKey := fmt.Sprintf("%s/%s", typ.String(), field.Name())
		decoder := fieldDecoders[fieldCacheKey]
		if decoder == nil {
			var err error
			decoder, err = decoderOfType(ctx.append(field.Name()), field.Type())
			if err != nil {
				return nil, err
			}
		}
		encoder := fieldEncoders[fieldCacheKey]
		if encoder == nil {
			var err error
			encoder, err = encoderOfType(ctx.append(field.Name()), field.Type())
			if err != nil {
				return nil, err
			}
		}
		binding := &binding{
			Field:     field,
			FromNames: fieldNames,
			ToNames:   fieldNames,
			Decoder:   decoder,
			Encoder:   encoder,
		}
		binding.levels = []int{i}
		bindings = append(bindings, binding)
	}
	return createStructDescriptor(typ, bindings, embeddedBindings), nil
}
func createStructDescriptor(typ reflect2.Type, bindings []*binding, embeddedBindings []*binding) *structDescriptor {
	structDescriptor := &structDescriptor{
		Type:   typ,
		Fields: bindings,
	}
	processTags(structDescriptor)
	// merge normal & embedded bindings & sort with original order
	allBindings := sortableBindings(append(embeddedBindings, structDescriptor.Fields...))
	sort.Sort(allBindings)
	structDescriptor.Fields = allBindings
	return structDescriptor
}

type sortableBindings []*binding

func (bindings sortableBindings) Len() int {
	return len(bindings)
}

func (bindings sortableBindings) Less(i, j int) bool {
	left := bindings[i].levels
	right := bindings[j].levels
	k := 0
	for {
		if left[k] < right[k] {
			return true
		} else if left[k] > right[k] {
			return false
		}
		k++
	}
}

func (bindings sortableBindings) Swap(i, j int) {
	bindings[i], bindings[j] = bindings[j], bindings[i]
}

func processTags(structDescriptor *structDescriptor) {
	for _, binding := range structDescriptor.Fields {
		shouldOmitEmpty := false
		tagParts := strings.Split(binding.Field.Tag().Get(tagKey), ",")
		for _, tagPart := range tagParts[1:] {
			if tagPart == "omitempty" {
				shouldOmitEmpty = true
			}
		}
		binding.Decoder = &structFieldDecoder{binding.Field, binding.Decoder}
		binding.Encoder = &structFieldEncoder{binding.Field, binding.Encoder, shouldOmitEmpty}
	}
}

func calcFieldNames(originalFieldName string, tagProvidedFieldName string, wholeTag string) []string {
	// ignore?
	if wholeTag == "-" {
		return []string{}
	}
	// rename?
	var fieldNames []string
	if tagProvidedFieldName == "" {
		fieldNames = []string{originalFieldName}
	} else {
		fieldNames = []string{tagProvidedFieldName}
	}
	// private?
	isNotExported := unicode.IsLower(rune(originalFieldName[0])) || originalFieldName[0] == '_'
	if isNotExported {
		fieldNames = []string{}
	}
	return fieldNames
}
