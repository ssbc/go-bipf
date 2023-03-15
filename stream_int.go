package bipf

import (
	"encoding/binary"
	"errors"
	"math"
)

func (stream *stream) WriteUint8(v uint8) error {
	return stream.WriteInt32(int32(v))
}

func (stream *stream) WriteInt8(v int8) error {
	return stream.WriteInt32(int32(v))
}

func (stream *stream) WriteUint16(v uint16) error {
	return stream.WriteInt32(int32(v))
}

func (stream *stream) WriteInt16(v int16) error {
	return stream.WriteInt32(int32(v))
}

func (stream *stream) WriteUint32(v uint32) error {
	return stream.WriteInt32(int32(v))
}

func (stream *stream) WriteInt32(v int32) error {
	stream.WriteTag(4, valueTypeInt)
	stream.buf = binary.LittleEndian.AppendUint32(stream.buf, uint32(v))
	return nil
}

func (stream *stream) WriteUint64(v uint64) error {
	if v > math.MaxUint32 {
		return errors.New("value > MaxUint32")
	}
	return stream.WriteInt32(int32(v))
}

func (stream *stream) WriteInt64(v int64) error {
	if v > math.MaxInt32 {
		return errors.New("value > MaxInt32")
	}
	if v < math.MinInt32 {
		return errors.New("value < MinInt32")
	}
	return stream.WriteInt32(int32(v))
}
