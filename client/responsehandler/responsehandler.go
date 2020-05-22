package responsehandler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

// EncodeResponse func
func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if err, ok := response.(error); ok && err != nil {
		EncodeError(ctx, err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// EncodeError func
func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case errors.New("Not found"):
		return http.StatusNotFound
	case errors.New("Inconsistent IDs"), errors.New("Missing parameter"):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
