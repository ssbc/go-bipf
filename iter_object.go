package bipf

import (
	"errors"
)

func (iter *iterator) ReadObjectCB(callback func(*iterator) error) error {
	v, l, err := iter.readTag()
	if err != nil {
		return err
	}

	if v != valueTypeObject {
		return iter.annotateError(errors.New("expected an object"))
	}

	if l == 0 {
		return nil
	}

	if err := iter.incrementDepth(); err != nil {
		return err
	}

	start := iter.numRead()

	for iter.numRead()-start < l {
		if err := callback(iter); err != nil {
			return err
		}

		if iter.numRead()-start > l {
			return iter.annotateError(errors.New("out of bounds"))
		}
	}

	if err := iter.decrementDepth(); err != nil {
		return err
	}

	return nil
}
