package xelon

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyArgument          = errors.New("argument cannot be empty")
	ErrEmptyPayloadNotAllowed = errors.New("empty payload is not allowed")
)

type ErrorResponse struct {
	Response     *Response
	ErrorElement ErrorElement
}

type ErrorElement struct {
	Error string `json:"error,omitempty"`
	Code  int    `json:"code,omitempty"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %+v",
		r.Response.Request.Method, sanitizeURL(r.Response.Request.URL), r.Response.StatusCode, r.ErrorElement)
}
