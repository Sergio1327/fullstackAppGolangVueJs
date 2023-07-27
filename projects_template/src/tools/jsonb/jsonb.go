package jsonb

import (
	"encoding/json"

	"github.com/jmoiron/sqlx/types"
)

func Marshal(obj any) (j types.JSONText, err error) {
	var b []byte

	if b, err = json.Marshal(obj); err != nil {
		return
	}

	if err = j.Scan(b); err != nil {
		return
	}

	return
}

func Unmarshal[T any](j types.JSONText) (t T, err error) {
	err = j.Unmarshal(&t)
	return
}
