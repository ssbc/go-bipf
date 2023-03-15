package bipf

import (
	"encoding/binary"
	"errors"
	"math"
)

func (iter *iterator) ReadInt8() (int8, error) {
	val, err := iter.ReadInt32()
	if err != nil {
		return 0, err
	}
	if val > math.MaxInt8 {
		return 0, errors.New("overflow")
	}
	if val < math.MinInt8 {
		return 0, errors.New("underflow")
	}
	return int8(val), nil
}

func (iter *iterator) ReadUint8() (ret uint8, err error) {
	val, err := iter.ReadInt32()
	if err != nil {
		return 0, err
	}
	if val > math.MaxUint8 {
		return 0, errors.New("overflow")
	}
	if val < 0 {
		return 0, errors.New("underflow")
	}
	return uint8(val), nil
}

func (iter *iterator) ReadInt16() (int16, error) {
	val, err := iter.ReadInt32()
	if err != nil {
		return 0, err
	}
	if val > math.MaxInt16 {
		return 0, errors.New("overflow")
	}
	if val < math.MinInt16 {
		return 0, errors.New("underflow")
	}
	return int16(val), nil
}

func (iter *iterator) ReadUint16() (uint16, error) {
	val, err := iter.ReadInt32()
	if err != nil {
		return 0, err
	}
	if val > math.MaxUint16 {
		return 0, errors.New("overflow")
	}
	if val < 0 {
		return 0, errors.New("underflow")
	}
	return uint16(val), nil
}

func (iter *iterator) ReadInt32() (ret int32, err error) {
	v, l, err := iter.readTag()
	if err != nil {
		return 0, err
	}

	if l != 4 {
		return 0, errors.New("invalid length")
	}

	if v != valueTypeInt {
		return 0, errors.New("invalid type")
	}

	buf := make([]byte, 4)
	_, err = iter.Read(buf)
	if err != nil {
		return 0, err
	}

	u := binary.LittleEndian.Uint32(buf)
	return int32(u), nil
}

func (iter *iterator) ReadUint32() (uint32, error) {
	val, err := iter.ReadInt32()
	if err != nil {
		return 0, err
	}
	if val < 0 {
		return 0, errors.New("underflow")
	}
	return uint32(val), nil
}

func (iter *iterator) ReadInt64() (int64, error) {
	val, err := iter.ReadInt32()
	return int64(val), err
}

func (iter *iterator) ReadUint64() (uint64, error) {
	val, err := iter.ReadInt32()
	return uint64(val), err
}
