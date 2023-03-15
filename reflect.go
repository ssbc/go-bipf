package bipf

import (
	"errors"
	"github.com/modern-go/reflect2"
	"reflect"
	"unsafe"
)

type valDecoder interface {
	Decode(ptr unsafe.Pointer, iter *iterator) error
}

type valEncoder interface {
	IsEmpty(ptr unsafe.Pointer) (bool, error)
	Encode(ptr unsafe.Pointer, stream *stream) error
}

func (stream *stream) WriteVal(val any) error {
	if nil == val {
		stream.WriteNil()
		return nil
	}
	cacheKey := reflect2.RTypeOf(val)
	encoder := encCache.getEncoderFromCache(cacheKey)
	if encoder == nil {
		typ := reflect2.TypeOf(val)
		var err error
		encoder, err = encoderOf(typ)
		if err != nil {
			return err
		}
	}
	return encoder.Encode(reflect2.PtrOf(val), stream)
}

type checkIsEmpty interface {
	IsEmpty(ptr unsafe.Pointer) (bool, error)
}

func encoderOf(typ reflect2.Type) (valEncoder, error) {
	cacheKey := typ.RType()
	encoder := encCache.getEncoderFromCache(cacheKey)
	if encoder != nil {
		return encoder, nil
	}
	ctx := &ctx{
		prefix:   "",
		decoders: map[reflect2.Type]valDecoder{},
		encoders: map[reflect2.Type]valEncoder{},
	}
	encoder, err := encoderOfType(ctx, typ)
	if err != nil {
		return nil, err
	}
	if typ.LikePtr() {
		encoder = &onePtrEncoder{encoder}
	}
	encCache.addEncoderToCache(cacheKey, encoder)
	return encoder, nil
}

func decoderOf(typ reflect2.Type) (valDecoder, error) {
	cacheKey := typ.RType()
	decoder := decCache.getDecoderFromCache(cacheKey)
	if decoder != nil {
		return decoder, nil
	}
	ctx := &ctx{
		prefix:   "",
		decoders: map[reflect2.Type]valDecoder{},
		encoders: map[reflect2.Type]valEncoder{},
	}
	ptrType := typ.(*reflect2.UnsafePtrType)
	decoder, err := decoderOfType(ctx, ptrType.Elem())
	if err != nil {
		return nil, err
	}
	decCache.addDecoderToCache(cacheKey, decoder)
	return decoder, nil
}

func decoderOfType(ctx *ctx, typ reflect2.Type) (valDecoder, error) {
	decoder := ctx.decoders[typ]
	if decoder != nil {
		return decoder, nil
	}
	placeholder := &placeholderDecoder{}
	ctx.decoders[typ] = placeholder
	decoder, err := createDecoderOfType(ctx, typ)
	if err != nil {
		return nil, err
	}
	placeholder.decoder = decoder
	return decoder, nil
}

func createDecoderOfType(ctx *ctx, typ reflect2.Type) (valDecoder, error) {
	decoder := createDecoderOfMarshaler(typ)
	if decoder != nil {
		return decoder, nil
	}
	decoder, err := createDecoderOfNative(ctx, typ)
	if err == nil {
		return decoder, nil
	}
	switch typ.Kind() {
	case reflect.Interface:
		ifaceType, isIFace := typ.(*reflect2.UnsafeIFaceType)
		if isIFace {
			return &ifaceDecoder{valType: ifaceType}, nil
		}
		return &efaceDecoder{}, nil
	case reflect.Struct:
		return decoderOfStruct(ctx, typ)
	case reflect.Array:
		return decoderOfArray(ctx, typ)
	case reflect.Slice:
		return decoderOfSlice(ctx, typ)
	case reflect.Map:
		return decoderOfMap(ctx, typ)
	case reflect.Ptr:
		return decoderOfOptional(ctx, typ)
	}
	return nil, errors.New("decoder not found")
}

type ctx struct {
	prefix   string
	encoders map[reflect2.Type]valEncoder
	decoders map[reflect2.Type]valDecoder
}

func (b *ctx) append(prefix string) *ctx {
	return &ctx{
		prefix:   b.prefix + " " + prefix,
		encoders: b.encoders,
		decoders: b.decoders,
	}
}

func encoderOfType(ctx *ctx, typ reflect2.Type) (valEncoder, error) {
	encoder := ctx.encoders[typ]
	if encoder != nil {
		return encoder, nil
	}
	placeholder := &placeholderEncoder{}
	ctx.encoders[typ] = placeholder
	encoder, err := createEncoderOfType(ctx, typ)
	if err != nil {
		return nil, err
	}
	placeholder.encoder = encoder
	return encoder, nil
}

func createEncoderOfType(ctx *ctx, typ reflect2.Type) (valEncoder, error) {
	encoder, err := createEncoderOfMarshaler(ctx, typ)
	if err == nil {
		return encoder, nil
	}
	encoder, err = createEncoderOfNative(ctx, typ)
	if err == nil {
		return encoder, nil
	}
	kind := typ.Kind()
	switch kind {
	case reflect.Interface:
		return &dynamicEncoder{typ}, nil
	case reflect.Struct:
		return encoderOfStruct(ctx, typ)
	case reflect.Array:
		return encoderOfArray(ctx, typ)
	case reflect.Slice:
		return encoderOfSlice(ctx, typ)
	case reflect.Map:
		return encoderOfMap(ctx, typ)
	case reflect.Ptr:
		return encoderOfOptional(ctx, typ)
	}
	return nil, errors.New("encoder not found")
}

type placeholderDecoder struct {
	decoder valDecoder
}

func (decoder *placeholderDecoder) Decode(ptr unsafe.Pointer, iter *iterator) error {
	return decoder.decoder.Decode(ptr, iter)
}

type placeholderEncoder struct {
	encoder valEncoder
}

func (encoder *placeholderEncoder) Encode(ptr unsafe.Pointer, stream *stream) error {
	return encoder.encoder.Encode(ptr, stream)
}

func (encoder *placeholderEncoder) IsEmpty(ptr unsafe.Pointer) (bool, error) {
	return encoder.encoder.IsEmpty(ptr)
}

type onePtrEncoder struct {
	encoder valEncoder
}

func (encoder *onePtrEncoder) IsEmpty(ptr unsafe.Pointer) (bool, error) {
	return encoder.encoder.IsEmpty(unsafe.Pointer(&ptr))
}

func (encoder *onePtrEncoder) Encode(ptr unsafe.Pointer, stream *stream) error {
	return encoder.encoder.Encode(unsafe.Pointer(&ptr), stream)
}

func (iter *iterator) ReadVal(obj any) error {
	depth := iter.depth
	cacheKey := reflect2.RTypeOf(obj)
	decoder := decCache.getDecoderFromCache(cacheKey)
	if decoder == nil {
		typ := reflect2.TypeOf(obj)
		if typ == nil || typ.Kind() != reflect.Ptr {
			return errors.New("can only unmarshal into pointer")
		}
		var err error
		decoder, err = decoderOf(typ)
		if err != nil {
			return err
		}
	}
	ptr := reflect2.PtrOf(obj)
	if ptr == nil {
		return errors.New("can not read into nil pointer")
	}
	if err := decoder.Decode(ptr, iter); err != nil {
		return err
	}
	if iter.depth != depth {
		return errors.New("unexpected mismatched nesting")
	}
	return nil
}
