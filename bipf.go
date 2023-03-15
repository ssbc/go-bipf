package bipf

import (
	"errors"
	"io"
)

const tagKey = "bipf"

// Marshal returns the BIPF encoding of v.
//
// Marshal traverses the value v recursively. If an encountered value implements
// the Marshaler interface and is not a nil pointer, Marshal calls its
// MarshalBIPF method to produce BIPF. If no MarshalBIPF method is present but
// the value implements encoding.BinaryMarshaler instead, Marshal calls its
// MarshalBinary method and encodes the result as a BIPF BUFFER. The nil pointer
// exception is not strictly necessary but mimics a similar, necessary exception
// in the behavior of UnmarshalBIPF.
//
// Otherwise, Marshal uses the following type-dependent default encodings:
//
// Boolean values encode as BIPF BULLNULL.
//
// Floating point numbers encode as BIPF DOUBLE.
//
// Integer numbers encode as BIPF INT. BIPF only supports 32-bit integers. The
// value of the encoded integer must fit in 32 bytes or an error will be
// returned.
//
// String values directly encode as BIPF STRING.
//
// Array and slice values encode as BIPF ARRAY, except that []byte encodes as a
// BIPF BUFFER, and a nil slice encodes as the BIPF BOOLNULL.
//
// Struct values encode as BIPF OBJECT. Each exported struct field becomes a
// member of the object, using the field name as the object key, unless the
// field is omitted for one of the reasons given below.
//
// The encoding of each struct field can be customized by the format string
// stored under the "bipf" key in the struct field's tag. The format string
// gives the name of the field, possibly followed by a comma-separated list of
// options. The name may be empty in order to specify options without overriding
// the default field name.
//
// The "omitempty" option specifies that the field should be omitted from the
// encoding if the field has an empty value, defined as false, 0, a nil pointer,
// a nil interface value, and any empty array, slice, map, or string.
//
// As a special case, if the field tag is "-", the field is always omitted. Note
// that a field with name "-" can still be generated using the tag "-,".
//
// Examples of struct field tags and their meanings:
//
//	// Field appears in BIPF as key "myName".
//	Field int `bipf:"myName"`
//
//	// Field appears in BIPF as key "myName" and
//	// the field is omitted from the object if its value is empty,
//	// as defined above.
//	Field int `bipf:"myName,omitempty"`
//
//	// Field appears in BIPF as key "Field" (the default), but
//	// the field is skipped if empty.
//	// Note the leading comma.
//	Field int `bipf:",omitempty"`
//
//	// Field is ignored by this package.
//	Field int `bipf:"-"`
//
//	// Field appears in BIPF as key "-".
//	Field int `bipf:"-,"`
//
// Anonymous struct fields are usually marshaled as if their inner exported
// fields were fields in the outer struct, subject to the usual Go visibility
// rules amended as described in the next paragraph. An anonymous struct field
// with a name given in its BIPF tag is treated as having that name, rather than
// being anonymous. An anonymous struct field of interface type is treated the
// same as having that type as its name, rather than being anonymous.
//
// The Go visibility rules for struct fields are amended for BIPF when deciding
// which field to marshal or unmarshal. If there are multiple fields at the same
// level, and that level is the least nested (and would therefore be the nesting
// level selected by the usual Go rules), the following extra rules apply:
//
// 1) Of those fields, if any are BIPF-tagged, only tagged fields are considered,
// even if there are multiple untagged fields that would otherwise conflict.
//
// 2) If there is exactly one field (tagged or not according to the first rule), that is selected.
//
// 3) Otherwise there are multiple fields, and all are ignored; no error occurs.
//
// Map values encode as BIPF objects.
//
// Pointer values encode as the value pointed to.
// A nil pointer encodes as the BIPF BOOLNULL value.
//
// Interface values encode as the value contained in the interface.
// A nil interface value encodes as the BIPF BOOLNULL value.
//
// Channel, complex, and function values cannot be encoded in BIPF.
func Marshal(v any) ([]byte, error) {
	stream := streamPool.BorrowStream(nil)
	defer streamPool.ReturnStream(stream)
	if err := stream.WriteVal(v); err != nil {
		return nil, err
	}
	result := stream.Buffer()
	copied := make([]byte, len(result))
	copy(copied, result)
	return copied, nil
}

// Unmarshal parses the BIPF-encoded data and stores the result
// in the value pointed to by v. If v is nil or not a pointer,
// Unmarshal returns an error.
//
// Unmarshal uses the inverse of the encodings that
// Marshal uses, allocating maps, slices, and pointers as necessary.
//
// To unmarshal BIPF into a value implementing the Unmarshaler interface,
// Unmarshal calls that value's UnmarshalBIPF method, including when the input
// is a BIPF BOOLNULL. Otherwise, if the value implements
// encoding.BinaryUnmarshaler and the input is a BIPF BUFFER, Unmarshal calls
// that value's UnmarshalBinary method with the contents of the BUFFER.
//
// To unmarshal BIPF into a struct, Unmarshal matches incoming object
// keys to the keys used by Marshal (either the struct field name or its tag),
// preferring an exact match but also accepting a case-insensitive match. By
// default, object keys which don't have a corresponding struct field are
// ignored.
//
// To unmarshal BIPF into an interface value,
// Unmarshal stores one of these in the interface value:
//
//	string, for BIPF STRING
//	[]byte for BIPF BUFFER
//	int32 for BIPF INT
//	float64, for BIPF DOUBLE
//	[]any, for BIPF ARRAY
//	map[any]any, for BIPF OBJECT
//	bool, for BIPF BOOLNULL not set to null
//	nil for BIPF BULLNULL set to null
//
// To unmarshal a BIPF array into a slice, Unmarshal resets the slice length
// to zero and then appends each element to the slice.
//
// If the BIPF-encoded data contain an error, Unmarshal returns an error.
//
// If a BIPF value is not appropriate for a given target type, or if a BIPF
// number overflows or underflows the target type, Unmarshal returns an error.
//
// When unmarshaling BIPF STRING, invalid UTF-8 is not treated as an error.
func Unmarshal(data []byte, v any) error {
	iter := iteratorPool.BorrowIterator(data)
	defer iteratorPool.ReturnIterator(iter)
	if err := iter.ReadVal(v); err != nil {
		return err
	}
	_, err := iter.ReadByte()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}
	return errors.New("there are bytes left after unmarshal")
}

type valueType byte

const (
	valueTypeString   valueType = 0b000
	valueTypeBuffer   valueType = 0b001
	valueTypeInt      valueType = 0b010
	valueTypeDouble   valueType = 0b011
	valueTypeArray    valueType = 0b100
	valueTypeObject   valueType = 0b101
	valueTypeBoolNull valueType = 0b110
	valueTypeExtended valueType = 0b111
)
