package bipf

import (
	"errors"
)

func (iter *iterator) ReadString() (string, error) {
	typ, length, err := iter.readTag()
	if err != nil {
		return "", err
	}
	if typ != valueTypeString {
		return "", errors.New("expected a string")
	}

	str := make([]byte, length)
	_, err = iter.Read(str)
	if err != nil {
		return "", err
	}
	return string(str), nil
}
