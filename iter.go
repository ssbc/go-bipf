package bipf

import (
	"encoding/binary"
	"errors"
	"io"
)

const maxDepth = 10000

type iterator struct {
	reader           io.Reader
	numOfReadBytes   int
	buf              []byte
	head             int
	tail             int
	depth            int
	captureStartedAt int
	captured         []byte
}

func newIterator() *iterator {
	return &iterator{
		reader: nil,
		buf:    nil,
		head:   0,
		tail:   0,
		depth:  0,
	}
}

func (iter *iterator) Reset(reader io.Reader) *iterator {
	iter.reader = reader
	iter.head = 0
	iter.tail = 0
	iter.depth = 0
	iter.numOfReadBytes = 0
	return iter
}

func (iter *iterator) ResetBytes(input []byte) *iterator {
	iter.reader = nil
	iter.buf = input
	iter.head = 0
	iter.tail = len(input)
	iter.depth = 0
	iter.numOfReadBytes = 0
	return iter
}

func (iter *iterator) whatIsNext() (valueType, error) {
	b, err := iter.ReadByte()
	if err != nil {
		return 0, err
	}
	typ := b & 0x07
	iter.unreadByte()
	return valueType(typ), nil
}

func (iter *iterator) readTag() (valueType, uint64, error) {
	v, err := binary.ReadUvarint(iter)
	if err != nil {
		return 0, 0, iter.annotateError(wrapf(err, "error reading uvarint"))
	}
	typ := byte(v) & 0x07
	length := v >> 3
	return valueType(typ), length, nil
}

func (iter *iterator) Read(b []byte) (n int, err error) {
	for i := 0; i < len(b); i++ {
		rb, err := iter.ReadByte()
		if err != nil {
			return n, err
		}
		b[i] = rb
		n++
	}
	return n, nil
}

func (iter *iterator) ReadByte() (ret byte, err error) {
	if iter.head == iter.tail {
		if err := iter.loadMore(); err != nil {
			return 0, err
		}
		ret = iter.buf[iter.head]
		iter.head++
		iter.numOfReadBytes++
		return ret, nil
	}
	ret = iter.buf[iter.head]
	iter.head++
	iter.numOfReadBytes++
	return ret, nil
}

func (iter *iterator) annotateError(err error) error {
	peekStart := iter.head - 10
	if peekStart < 0 {
		peekStart = 0
	}
	peekEnd := iter.head + 10
	if peekEnd > iter.tail {
		peekEnd = iter.tail
	}
	parsing := string(iter.buf[peekStart:peekEnd])
	contextStart := iter.head - 50
	if contextStart < 0 {
		contextStart = 0
	}
	contextEnd := iter.head + 50
	if contextEnd > iter.tail {
		contextEnd = iter.tail
	}
	context := string(iter.buf[contextStart:contextEnd])
	return wrapf(err, "error found in #%v byte of ...|%x|..., bigger context ...|%x|...", iter.head-peekStart, parsing, context)
}

func (iter *iterator) loadMore() error {
	if iter.reader == nil {
		iter.head = iter.tail
		return io.EOF
	}
	if iter.captured != nil {
		iter.captured = append(iter.captured,
			iter.buf[iter.captureStartedAt:iter.tail]...)
		iter.captureStartedAt = 0
	}
	for {
		n, err := iter.reader.Read(iter.buf)
		if n == 0 {
			if err != nil {
				return err
			}
		} else {
			iter.head = 0
			iter.tail = n
			return nil
		}
	}
}

func (iter *iterator) unreadByte() {
	iter.numOfReadBytes--
	iter.head--
	return
}

func (iter *iterator) ReadAny() (any, error) {
	valueType, err := iter.whatIsNext()
	if err != nil {
		return nil, err
	}

	switch valueType {
	case valueTypeString:
		return iter.ReadString()
	case valueTypeBuffer:
		return iter.ReadBuffer()
	case valueTypeInt:
		return iter.ReadInt32()
	case valueTypeDouble:
		return iter.ReadFloat64()
	case valueTypeArray:
		var arr []any
		if err := iter.ReadArrayCB(func(iter *iterator) error {
			var elem any
			if err := iter.ReadVal(&elem); err != nil {
				return err
			}
			arr = append(arr, elem)
			return nil
		}); err != nil {
			return nil, err
		}
		return arr, nil
	case valueTypeObject:
		obj := make(map[any]any)
		if err := iter.ReadObjectCB(func(Iter *iterator) error {
			var key any
			if err := iter.ReadVal(&key); err != nil {
				return err
			}

			var value any
			if err := iter.ReadVal(&value); err != nil {
				return err
			}

			obj[key] = value
			return nil
		}); err != nil {
			return nil, err
		}
		return obj, nil
	case valueTypeBoolNull:
		ok, err := iter.CheckNilIsNext()
		if err != nil {
			return nil, err
		}
		if ok {
			return nil, iter.ReadNil()
		}
		return iter.ReadBool()
	default:
		return nil, iter.annotateError(errors.New("unsupported type"))
	}
}

func (iter *iterator) incrementDepth() error {
	iter.depth++
	if iter.depth <= maxDepth {
		return nil
	}
	return errors.New("exceeded max depth")
}

func (iter *iterator) decrementDepth() error {
	iter.depth--
	if iter.depth >= 0 {
		return nil
	}
	return errors.New("unexpected negative nesting")
}

func (iter *iterator) numRead() uint64 {
	return uint64(iter.numOfReadBytes)
}

func (iter *iterator) skip() error {
	_, l, err := iter.readTag()
	if err != nil {
		return err
	}

	for i := 0; i < int(l); i++ {
		_, err := iter.ReadByte()
		if err != nil {
			return err
		}
	}

	return nil
}

func (iter *iterator) ReadBuffer() ([]byte, error) {
	v, l, err := iter.readTag()
	if err != nil {
		return nil, err
	}

	if v != valueTypeBuffer {
		return nil, errors.New("unexpected type")
	}

	b := make([]byte, int(l))
	_, err = iter.Read(b)
	return b, err
}
