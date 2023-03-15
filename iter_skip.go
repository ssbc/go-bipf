package bipf

import (
	"errors"
)

func (iter *iterator) CheckNilIsNext() (bool, error) {
	b, err := iter.ReadByte()
	if err != nil {
		return false, err
	}

	if b == byte(valueTypeBoolNull) && b>>3 == 0 {
		return true, nil
	}

	iter.unreadByte()
	return false, nil
}

func (iter *iterator) ReadNil() error {
	b, err := iter.ReadByte()
	if err != nil {
		return err
	}

	if b == byte(valueTypeBoolNull) && b>>3 == 0 {
		return nil
	}

	return errors.New("this value isn't a nil")
}

func (iter *iterator) ReadBool() (bool, error) {
	v, l, err := iter.readTag()
	if err != nil {
		return false, err
	}

	if v != valueTypeBoolNull {
		return false, errors.New("invalid type")
	}

	if l != 1 {
		return false, errors.New("invalid length")
	}

	b, err := iter.ReadByte()
	if err != nil {
		return false, err
	}

	switch b {
	case 0x00:
		return false, nil
	case 0x01:
		return true, nil
	default:
		return false, errors.New("invalid bool value")
	}
}
