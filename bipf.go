package bipf

import (
	"errors"
	"io"
)

const tagKey = "bipf"

func Marshal(v any) ([]byte, error) {
	stream := streamPool.BorrowStream(nil)
	defer streamPool.ReturnStream(stream)
	if err := stream.WriteVal(v); err != nil {
		return nil, err
	}
	result := stream.Buffer()
	copied := make([]byte, len(result))
	copy(copied, result)
	return copied, nil
}

func Unmarshal(data []byte, v any) error {
	iter := iteratorPool.BorrowIterator(data)
	defer iteratorPool.ReturnIterator(iter)
	if err := iter.ReadVal(v); err != nil {
		return err
	}
	_, err := iter.ReadByte()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}
	return errors.New("there are bytes left after unmarshal")
}

type valueType byte

const (
	valueTypeString   valueType = 0b000
	valueTypeBuffer   valueType = 0b001
	valueTypeInt      valueType = 0b010
	valueTypeDouble   valueType = 0b011
	valueTypeArray    valueType = 0b100
	valueTypeObject   valueType = 0b101
	valueTypeBoolNull valueType = 0b110
	valueTypeExtended valueType = 0b111
)
