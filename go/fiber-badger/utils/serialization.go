package utils

import (
	"bytes"
	"encoding/gob"
)

func MarshalStruct(s interface{}) ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(s)
	return b.Bytes(), err
}

func UnmarshalStruct[T interface{}](data []byte) (T, error) {
	var result T
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	err := dec.Decode(&result)

	return result, err
}
