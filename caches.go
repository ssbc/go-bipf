package bipf

import "github.com/modern-go/concurrent"

var (
	encCache = newEncoderCache()
	decCache = newDecoderCache()
)

type encoderCache struct {
	encoderCache *concurrent.Map
}

func newEncoderCache() *encoderCache {
	return &encoderCache{
		encoderCache: concurrent.NewMap(),
	}
}
func (c *encoderCache) addEncoderToCache(cacheKey uintptr, encoder valEncoder) {
	c.encoderCache.Store(cacheKey, encoder)
}

func (c *encoderCache) getEncoderFromCache(cacheKey uintptr) valEncoder {
	encoder, found := c.encoderCache.Load(cacheKey)
	if found {
		return encoder.(valEncoder)
	}
	return nil
}

type decoderCache struct {
	decoderCache *concurrent.Map
}

func newDecoderCache() *decoderCache {
	return &decoderCache{
		decoderCache: concurrent.NewMap(),
	}
}

func (c *decoderCache) addDecoderToCache(cacheKey uintptr, decoder valDecoder) {
	c.decoderCache.Store(cacheKey, decoder)
}

func (c *decoderCache) getDecoderFromCache(cacheKey uintptr) valDecoder {
	decoder, found := c.decoderCache.Load(cacheKey)
	if found {
		return decoder.(valDecoder)
	}
	return nil
}
