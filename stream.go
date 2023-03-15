package bipf

import (
	"encoding/binary"
	"io"
	"math"
)

type stream struct {
	out io.Writer
	buf []byte
}

func newStream(out io.Writer, bufSize int) *stream {
	return &stream{
		out: out,
		buf: make([]byte, 0, bufSize),
	}
}

func (stream *stream) Reset(out io.Writer) {
	stream.out = out
	stream.buf = stream.buf[:0]
}

func (stream *stream) Buffered() int {
	return len(stream.buf)
}

func (stream *stream) Buffer() []byte {
	return stream.buf
}

func (stream *stream) Write(p []byte) (nn int, err error) {
	stream.buf = append(stream.buf, p...)
	if stream.out != nil {
		nn, err = stream.out.Write(stream.buf)
		stream.buf = stream.buf[nn:]
		return
	}
	return len(p), nil
}

func (stream *stream) Flush() error {
	if stream.out == nil {
		return nil
	}
	_, err := stream.out.Write(stream.buf)
	if err != nil {
		return err
	}
	stream.buf = stream.buf[:0]
	return nil
}

func (stream *stream) WriteNil() {
	stream.WriteTag(0, valueTypeBoolNull)
}

func (stream *stream) WriteBool(val bool) error {
	stream.WriteTag(1, valueTypeBoolNull)
	if val {
		stream.buf = append(stream.buf, 1)
	} else {
		stream.buf = append(stream.buf, 0)
	}
	return nil
}

func (stream *stream) WriteEmptyObject() {
	stream.WriteTag(0, valueTypeObject)
}

func (stream *stream) WriteEmptyArray() {
	stream.WriteTag(0, valueTypeArray)
}

func (stream *stream) WriteString(s string) error {
	stream.WriteTag(uint64(len(s)), valueTypeString)
	stream.buf = append(stream.buf, s...)
	return nil
}

func (stream *stream) WriteFloat32(val float32) error {
	return stream.WriteFloat64(float64(val))
}

func (stream *stream) WriteFloat64(val float64) error {
	bits := math.Float64bits(val)
	stream.WriteTag(8, valueTypeDouble)
	stream.buf = binary.LittleEndian.AppendUint64(stream.buf, bits)
	return nil
}

func (stream *stream) WriteBuffer(b []byte) error {
	stream.WriteTag(uint64(len(b)), valueTypeBuffer)
	stream.buf = append(stream.buf, b...)
	return nil
}

func (stream *stream) WriteTag(length uint64, typ valueType) {
	v := length<<3 | uint64(typ)
	stream.buf = binary.AppendUvarint(stream.buf, v)
}
