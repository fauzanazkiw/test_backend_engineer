package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const maxRequestBodyBytes = 1 << 20 // 1 MB

// decodeJSON membaca satu object JSON dari request body.
func decodeJSON(r *http.Request, dst any) error {
	if r.Body == nil {
		return errors.New("request body is empty")
	}

	decoder := json.NewDecoder(io.LimitReader(r.Body, maxRequestBodyBytes))
	if err := decoder.Decode(dst); err != nil {
		return err
	}

	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain a single JSON object")
	}

	return nil
}
