package bipf

import (
	"io"
	"sync"
)

var (
	streamPool   = newSyncStreamPool()
	iteratorPool = newSyncIteratorPool()
)

type syncStreamPool struct {
	pool *sync.Pool
}

func newSyncStreamPool() *syncStreamPool {
	return &syncStreamPool{
		pool: &sync.Pool{
			New: func() any {
				return newStream(nil, 512)
			},
		},
	}
}

func (p *syncStreamPool) BorrowStream(writer io.Writer) *stream {
	stream := p.pool.Get().(*stream)
	stream.Reset(writer)
	return stream
}

func (p *syncStreamPool) ReturnStream(stream *stream) {
	stream.out = nil
	p.pool.Put(stream)
}

type syncIteratorPool struct {
	pool *sync.Pool
}

func newSyncIteratorPool() *syncIteratorPool {
	return &syncIteratorPool{
		pool: &sync.Pool{
			New: func() any {
				return newIterator()
			},
		},
	}
}

func (p *syncIteratorPool) BorrowIterator(data []byte) *iterator {
	iter := p.pool.Get().(*iterator)
	iter.ResetBytes(data)
	return iter
}

func (p *syncIteratorPool) ReturnIterator(iter *iterator) {
	p.pool.Put(iter)
}
