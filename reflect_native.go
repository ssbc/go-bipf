package bipf

import (
	"errors"
	"github.com/modern-go/reflect2"
	"reflect"
	"strconv"
	"unsafe"
)

const ptrSize = 32 << uintptr(^uintptr(0)>>63)

func createEncoderOfNative(ctx *ctx, typ reflect2.Type) (valEncoder, error) {
	kind := typ.Kind()

	if kind == reflect.Slice && typ.(reflect2.SliceType).Elem().Kind() == reflect.Uint8 {
		return &bytesCodec{}, nil
	}

	typeName := typ.String()
	switch kind {
	case reflect.String:
		if typeName != "string" {
			return encoderOfType(ctx, reflect2.TypeOfPtr((*string)(nil)).Elem())
		}
		return &stringCodec{}, nil
	case reflect.Int:
		if typeName != "int" {
			return encoderOfType(ctx, reflect2.TypeOfPtr((*int)(nil)).Elem())
		}
		if strconv.IntSize == 32 {
			return &int32Codec{}, nil
		}
		return &int64Codec{}, nil
	case reflect.Int8:
		if typeName != "int8" {
			return encoderOfType(ctx, reflect2.TypeOfPtr((*int8)(nil)).Elem())
		}
		return &int8Codec{}, nil
	case reflect.Int16:
		if typeName != "int16" {
			return encoderOfType(ctx, reflect2.TypeOfPtr((*int16)(nil)).Elem())
		}
		return &int16Codec{}, nil
	case reflect.Int32:
		if typeName != "int32" {
			return encoderOfType(ctx, reflect2.TypeOfPtr((*int32)(nil)).Elem())
		}
		return &int32Codec{}, nil
	case reflect.Int64:
		if typeName != "int64" {
			return encoderOfType(ctx, reflect2.TypeOfPtr((*int64)(nil)).Elem())
		}
		return &int64Codec{}, nil
	case reflect.Uint:
		if typeName != "uint" {
			return encoderOfType(ctx, reflect2.TypeOfPtr((*uint)(nil)).Elem())
		}
		if strconv.IntSize == 32 {
			return &uint32Codec{}, nil
		}
		return &uint64Codec{}, nil
	case reflect.Uint8:
		if typeName != "uint8" {
			return encoderOfType(ctx, reflect2.TypeOfPtr((*uint8)(nil)).Elem())
		}
		return &uint8Codec{}, nil
	case reflect.Uint16:
		if typeName != "uint16" {
			return encoderOfType(ctx, reflect2.TypeOfPtr((*uint16)(nil)).Elem())
		}
		return &uint16Codec{}, nil
	case reflect.Uint32:
		if typeName != "uint32" {
			return encoderOfType(ctx, reflect2.TypeOfPtr((*uint32)(nil)).Elem())
		}
		return &uint32Codec{}, nil
	case reflect.Uintptr:
		if typeName != "uintptr" {
			return encoderOfType(ctx, reflect2.TypeOfPtr((*uintptr)(nil)).Elem())
		}
		if ptrSize == 32 {
			return &uint32Codec{}, nil
		}
		return &uint64Codec{}, nil
	case reflect.Uint64:
		if typeName != "uint64" {
			return encoderOfType(ctx, reflect2.TypeOfPtr((*uint64)(nil)).Elem())
		}
		return &uint64Codec{}, nil
	case reflect.Float32:
		if typeName != "float32" {
			return encoderOfType(ctx, reflect2.TypeOfPtr((*float32)(nil)).Elem())
		}
		return &float32Codec{}, nil
	case reflect.Float64:
		if typeName != "float64" {
			return encoderOfType(ctx, reflect2.TypeOfPtr((*float64)(nil)).Elem())
		}
		return &float64Codec{}, nil
	case reflect.Bool:
		if typeName != "bool" {
			return encoderOfType(ctx, reflect2.TypeOfPtr((*bool)(nil)).Elem())
		}
		return &boolCodec{}, nil
	default:
		return nil, errors.New("native codec not found")
	}
}

func createDecoderOfNative(ctx *ctx, typ reflect2.Type) (valDecoder, error) {
	if typ.Kind() == reflect.Slice && typ.(reflect2.SliceType).Elem().Kind() == reflect.Uint8 {
		return &bytesCodec{}, nil
	}
	typeName := typ.String()
	switch typ.Kind() {
	case reflect.String:
		if typeName != "string" {
			return decoderOfType(ctx, reflect2.TypeOfPtr((*string)(nil)).Elem())
		}
		return &stringCodec{}, nil
	case reflect.Int:
		if typeName != "int" {
			return decoderOfType(ctx, reflect2.TypeOfPtr((*int)(nil)).Elem())
		}
		if strconv.IntSize == 32 {
			return &int32Codec{}, nil
		}
		return &int64Codec{}, nil
	case reflect.Int8:
		if typeName != "int8" {
			return decoderOfType(ctx, reflect2.TypeOfPtr((*int8)(nil)).Elem())
		}
		return &int8Codec{}, nil
	case reflect.Int16:
		if typeName != "int16" {
			return decoderOfType(ctx, reflect2.TypeOfPtr((*int16)(nil)).Elem())
		}
		return &int16Codec{}, nil
	case reflect.Int32:
		if typeName != "int32" {
			return decoderOfType(ctx, reflect2.TypeOfPtr((*int32)(nil)).Elem())
		}
		return &int32Codec{}, nil
	case reflect.Int64:
		if typeName != "int64" {
			return decoderOfType(ctx, reflect2.TypeOfPtr((*int64)(nil)).Elem())
		}
		return &int64Codec{}, nil
	case reflect.Uint:
		if typeName != "uint" {
			return decoderOfType(ctx, reflect2.TypeOfPtr((*uint)(nil)).Elem())
		}
		if strconv.IntSize == 32 {
			return &uint32Codec{}, nil
		}
		return &uint64Codec{}, nil
	case reflect.Uint8:
		if typeName != "uint8" {
			return decoderOfType(ctx, reflect2.TypeOfPtr((*uint8)(nil)).Elem())
		}
		return &uint8Codec{}, nil
	case reflect.Uint16:
		if typeName != "uint16" {
			return decoderOfType(ctx, reflect2.TypeOfPtr((*uint16)(nil)).Elem())
		}
		return &uint16Codec{}, nil
	case reflect.Uint32:
		if typeName != "uint32" {
			return decoderOfType(ctx, reflect2.TypeOfPtr((*uint32)(nil)).Elem())
		}
		return &uint32Codec{}, nil
	case reflect.Uintptr:
		if typeName != "uintptr" {
			return decoderOfType(ctx, reflect2.TypeOfPtr((*uintptr)(nil)).Elem())
		}
		if ptrSize == 32 {
			return &uint32Codec{}, nil
		}
		return &uint64Codec{}, nil
	case reflect.Uint64:
		if typeName != "uint64" {
			return decoderOfType(ctx, reflect2.TypeOfPtr((*uint64)(nil)).Elem())
		}
		return &uint64Codec{}, nil
	case reflect.Float32:
		if typeName != "float32" {
			return decoderOfType(ctx, reflect2.TypeOfPtr((*float32)(nil)).Elem())
		}
		return &float32Codec{}, nil
	case reflect.Float64:
		if typeName != "float64" {
			return decoderOfType(ctx, reflect2.TypeOfPtr((*float64)(nil)).Elem())
		}
		return &float64Codec{}, nil
	case reflect.Bool:
		if typeName != "bool" {
			return decoderOfType(ctx, reflect2.TypeOfPtr((*bool)(nil)).Elem())
		}
		return &boolCodec{}, nil
	}
	return nil, errors.New("decoder of native not found")
}

type stringCodec struct {
}

func (codec *stringCodec) Decode(ptr unsafe.Pointer, iter *iterator) error {
	ok, err := iter.CheckNilIsNext()
	if err != nil {
		return err
	}

	if !ok {
		*((*string)(ptr)), err = iter.ReadString()
		return err
	}

	return nil
}

func (codec *stringCodec) Encode(ptr unsafe.Pointer, stream *stream) error {
	str := *((*string)(ptr))
	return stream.WriteString(str)
}

func (codec *stringCodec) IsEmpty(ptr unsafe.Pointer) (bool, error) {
	return *((*string)(ptr)) == "", nil
}

type int8Codec struct {
}

func (codec *int8Codec) Decode(ptr unsafe.Pointer, iter *iterator) error {
	ok, err := iter.CheckNilIsNext()
	if err != nil {
		return err
	}

	if !ok {
		*((*int8)(ptr)), err = iter.ReadInt8()
		return err
	}

	return nil
}

func (codec *int8Codec) Encode(ptr unsafe.Pointer, stream *stream) error {
	return stream.WriteInt8(*((*int8)(ptr)))
}

func (codec *int8Codec) IsEmpty(ptr unsafe.Pointer) (bool, error) {
	return *((*int8)(ptr)) == 0, nil
}

type int16Codec struct {
}

func (codec *int16Codec) Decode(ptr unsafe.Pointer, iter *iterator) error {
	ok, err := iter.CheckNilIsNext()
	if err != nil {
		return err
	}

	if !ok {
		*((*int16)(ptr)), err = iter.ReadInt16()
		return err
	}

	return nil
}

func (codec *int16Codec) Encode(ptr unsafe.Pointer, stream *stream) error {
	return stream.WriteInt16(*((*int16)(ptr)))
}

func (codec *int16Codec) IsEmpty(ptr unsafe.Pointer) (bool, error) {
	return *((*int16)(ptr)) == 0, nil
}

type int32Codec struct {
}

func (codec *int32Codec) Decode(ptr unsafe.Pointer, iter *iterator) error {
	ok, err := iter.CheckNilIsNext()
	if err != nil {
		return err
	}

	if !ok {
		*((*int32)(ptr)), err = iter.ReadInt32()
		return err
	}

	return nil
}

func (codec *int32Codec) Encode(ptr unsafe.Pointer, stream *stream) error {
	return stream.WriteInt32(*((*int32)(ptr)))
}

func (codec *int32Codec) IsEmpty(ptr unsafe.Pointer) (bool, error) {
	return *((*int32)(ptr)) == 0, nil
}

type int64Codec struct {
}

func (codec *int64Codec) Decode(ptr unsafe.Pointer, iter *iterator) error {
	ok, err := iter.CheckNilIsNext()
	if err != nil {
		return err
	}

	if !ok {
		*((*int64)(ptr)), err = iter.ReadInt64()
		return err
	}

	return nil
}

func (codec *int64Codec) Encode(ptr unsafe.Pointer, stream *stream) error {
	return stream.WriteInt64(*((*int64)(ptr)))
}

func (codec *int64Codec) IsEmpty(ptr unsafe.Pointer) (bool, error) {
	return *((*int64)(ptr)) == 0, nil
}

type uint8Codec struct {
}

func (codec *uint8Codec) Decode(ptr unsafe.Pointer, iter *iterator) error {
	ok, err := iter.CheckNilIsNext()
	if err != nil {
		return err
	}

	if !ok {
		*((*uint8)(ptr)), err = iter.ReadUint8()
		return err
	}

	return nil
}

func (codec *uint8Codec) Encode(ptr unsafe.Pointer, stream *stream) error {
	return stream.WriteUint8(*((*uint8)(ptr)))
}

func (codec *uint8Codec) IsEmpty(ptr unsafe.Pointer) (bool, error) {
	return *((*uint8)(ptr)) == 0, nil
}

type uint16Codec struct {
}

func (codec *uint16Codec) Decode(ptr unsafe.Pointer, iter *iterator) error {
	ok, err := iter.CheckNilIsNext()
	if err != nil {
		return err
	}

	if !ok {
		*((*uint16)(ptr)), err = iter.ReadUint16()
		return err
	}

	return nil
}

func (codec *uint16Codec) Encode(ptr unsafe.Pointer, stream *stream) error {
	return stream.WriteUint16(*((*uint16)(ptr)))
}

func (codec *uint16Codec) IsEmpty(ptr unsafe.Pointer) (bool, error) {
	return *((*uint16)(ptr)) == 0, nil
}

type uint32Codec struct {
}

func (codec *uint32Codec) Decode(ptr unsafe.Pointer, iter *iterator) error {
	ok, err := iter.CheckNilIsNext()
	if err != nil {
		return err
	}

	if !ok {
		*((*uint32)(ptr)), err = iter.ReadUint32()
		return err
	}

	return nil
}

func (codec *uint32Codec) Encode(ptr unsafe.Pointer, stream *stream) error {
	return stream.WriteUint32(*((*uint32)(ptr)))
}

func (codec *uint32Codec) IsEmpty(ptr unsafe.Pointer) (bool, error) {
	return *((*uint32)(ptr)) == 0, nil
}

type uint64Codec struct {
}

func (codec *uint64Codec) Decode(ptr unsafe.Pointer, iter *iterator) error {
	ok, err := iter.CheckNilIsNext()
	if err != nil {
		return err
	}

	if !ok {
		*((*uint64)(ptr)), err = iter.ReadUint64()
		return err
	}

	return nil
}

func (codec *uint64Codec) Encode(ptr unsafe.Pointer, stream *stream) error {
	return stream.WriteUint64(*((*uint64)(ptr)))
}

func (codec *uint64Codec) IsEmpty(ptr unsafe.Pointer) (bool, error) {
	return *((*uint64)(ptr)) == 0, nil
}

type float32Codec struct {
}

func (codec *float32Codec) Decode(ptr unsafe.Pointer, iter *iterator) error {
	ok, err := iter.CheckNilIsNext()
	if err != nil {
		return err
	}

	if !ok {
		*((*float32)(ptr)), err = iter.ReadFloat32()
		return err
	}

	return nil
}

func (codec *float32Codec) Encode(ptr unsafe.Pointer, stream *stream) error {
	return stream.WriteFloat32(*((*float32)(ptr)))
}

func (codec *float32Codec) IsEmpty(ptr unsafe.Pointer) (bool, error) {
	return *((*float32)(ptr)) == 0, nil
}

type float64Codec struct {
}

func (codec *float64Codec) Decode(ptr unsafe.Pointer, iter *iterator) error {
	ok, err := iter.CheckNilIsNext()
	if err != nil {
		return err
	}

	if !ok {
		*((*float64)(ptr)), err = iter.ReadFloat64()
		return err
	}

	return nil
}

func (codec *float64Codec) Encode(ptr unsafe.Pointer, stream *stream) error {
	return stream.WriteFloat64(*((*float64)(ptr)))
}

func (codec *float64Codec) IsEmpty(ptr unsafe.Pointer) (bool, error) {
	return *((*float64)(ptr)) == 0, nil
}

type boolCodec struct {
}

func (codec *boolCodec) Decode(ptr unsafe.Pointer, iter *iterator) error {
	ok, err := iter.CheckNilIsNext()
	if err != nil {
		return err
	}

	if !ok {
		*((*bool)(ptr)), err = iter.ReadBool()
		return err
	}

	return nil
}

func (codec *boolCodec) Encode(ptr unsafe.Pointer, stream *stream) error {
	return stream.WriteBool(*((*bool)(ptr)))
}

func (codec *boolCodec) IsEmpty(ptr unsafe.Pointer) (bool, error) {
	return !(*((*bool)(ptr))), nil
}

type bytesCodec struct {
}

func (b bytesCodec) Decode(ptr unsafe.Pointer, iter *iterator) error {
	var err error
	*((*[]byte)(ptr)), err = iter.ReadBuffer()
	return err
}

func (b bytesCodec) Encode(ptr unsafe.Pointer, stream *stream) error {
	return stream.WriteBuffer(*((*[]byte)(ptr)))
}

func (b bytesCodec) IsEmpty(ptr unsafe.Pointer) (bool, error) {
	return len(*(*[]byte)(ptr)) == 0, nil
}
