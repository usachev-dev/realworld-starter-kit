package api_errors

import (
	"encoding/json"
	"net/http"
)

type E struct {
	code   int
	errors map[string][]string
}

func (e E) Error() string {
	result, _ := json.Marshal(e.errors)
	return string(result)
}

func (e E) Send(w http.ResponseWriter) {
	body, _ := json.Marshal(e.errors)
	w.WriteHeader(e.code)
	w.Write(body)
}

var Ok = E{
	http.StatusOK,
	map[string][]string{},
}

func (e E) IsOk() bool {
	return e.code == http.StatusOK
}

func NewError(code int) *E {
	return &E{
		code,
		map[string][]string{},
	}
}

func (e *E) Add(key string, error string) *E {
	if e.errors[key] != nil {
		e.errors[key] = append(e.errors[key], error)
	} else {
		e.errors[key] = []string{error}
	}
	return e
}

func (e *E) SetCode(code int) *E {
	e.code = code
	return e
}
