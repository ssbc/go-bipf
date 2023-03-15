package bipf

import (
	"encoding/binary"
	"errors"
	"math"
)

func (iter *iterator) ReadFloat32() (float32, error) {
	v, err := iter.ReadFloat64()
	if err != nil {
		return 0, err
	}
	if v > math.MaxFloat32 {
		return 0, errors.New("overflow")
	}
	if v < math.SmallestNonzeroFloat32 {
		return 0, errors.New("underflow")
	}
	return float32(v), nil
}

func (iter *iterator) ReadFloat64() (float64, error) {
	v, l, err := iter.readTag()
	if err != nil {
		return 0, err
	}

	if l != 8 {
		return 0, errors.New("invalid length")
	}

	if v != valueTypeDouble {
		return 0, errors.New("invalid type")
	}

	buf := make([]byte, 8)
	_, err = iter.Read(buf)
	if err != nil {
		return 0, err
	}

	u := binary.LittleEndian.Uint64(buf)
	return math.Float64frombits(u), nil
}
