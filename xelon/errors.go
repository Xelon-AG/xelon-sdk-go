package xelon

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrEmptyArgument          = errors.New("xelon-sdk-go: argument cannot be empty")
	ErrEmptyPayloadNotAllowed = errors.New("xelon-sdk-go: empty payload is not allowed")
)

type ErrorResponse struct {
	Response     *Response // HTTP response that caused this error.
	ErrorElement ErrorElement
}

type ErrorElement struct {
	Error       string         `json:"error,omitempty"`
	Message     string         `json:"message,omitempty"`
	Validations map[string]any `json:"errors,omitempty"`
}

func (e ErrorElement) String() string {
	var elements []string

	if e.Error != "" {
		elements = append(elements, fmt.Sprintf("error: %v", e.Error))
	}

	if e.Message != "" {
		elements = append(elements, fmt.Sprintf("details: %v", e.Message))
	}

	var validations []string
	for k, v := range e.Validations {
		validations = append(validations, fmt.Sprintf("%v - %v", k, v))
	}
	if len(validations) > 0 {
		elements = append(elements, fmt.Sprintf("validations: (%v)", strings.Join(validations, ", ")))
	}

	return strings.Join(elements, ", ")
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d (%+v)",
		r.Response.Request.Method, sanitizeURL(r.Response.Request.URL), r.Response.StatusCode, r.ErrorElement)
}
