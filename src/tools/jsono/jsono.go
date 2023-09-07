package jsono

import (
	"encoding/json"
)

func UpdateJSON[T any](data string, f func(T) (t T, err error)) (newData string, err error) {
	var tData T

	tData, err = Unmarshal[T](data)
	if err != nil {
		return
	}

	tData, err = f(tData)
	if err != nil {
		return
	}

	var b []byte

	b, err = json.Marshal(tData)
	if err != nil {
		return
	}

	newData = string(b)

	return
}

func Marshal(data any) (dataStr string, err error) {
	var b []byte
	b, err = json.Marshal(data)
	if err != nil {
		return
	}
	dataStr = string(b)
	return
}

func Unmarshal[T any](data string) (t T, err error) {
	err = json.Unmarshal([]byte(data), &t)
	return
}
