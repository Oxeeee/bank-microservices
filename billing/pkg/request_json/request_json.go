package jsonwrap

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type PointerToStruct[T any] interface{ *T }

func Unwrap[T any, P PointerToStruct[T]](dest P, r *http.Request) error {
	if r == nil || r.Body == nil {
		return errors.New("invalid request")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return errors.New("invalid body")
	}

	defer r.Body.Close()

	err = json.Unmarshal(body, &dest)
	if err != nil {
		return errors.New("cannot unmarshal json")
	}

	return nil
}
