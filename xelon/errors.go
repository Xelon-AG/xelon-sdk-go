package xelon

import (
	"encoding/json"
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
	Code  int          `json:"code,omitempty"`
	Error ErrorWrapper `json:"error,omitempty"`
}

// ErrorWrapper wraps a Error object with additional functionality.
type ErrorWrapper struct {
	Error
	Partial bool `json:"-"`
}

// Error represents an error object in Xelon business logic. It can be a string
// or array of validation messages.
type Error struct {
	Message     string `json:"-"`
	Validations map[string]interface{}
}

func (e ErrorElement) String() string {
	if e.Error.Partial {
		return fmt.Sprintf("(code: %v, error: %v)", e.Code, e.Error.Message)
	} else {
		var v []string
		for s, i := range e.Error.Validations {
			v = append(v, fmt.Sprintf("%v - %v", s, i))
		}
		validations := strings.Join(v, ", ")
		return fmt.Sprintf("(code: %v, validations: (%v))", e.Code, validations)
	}
}

func (r *ErrorResponse) Error() string {
	if r.Response.StackifyID != "" {
		return fmt.Sprintf("%v %v: %d (stackify id %v) %+v",
			r.Response.Request.Method, sanitizeURL(r.Response.Request.URL), r.Response.StatusCode, r.Response.StackifyID, r.ErrorElement)
	} else {
		return fmt.Sprintf("%v %v: %d %+v",
			r.Response.Request.Method, sanitizeURL(r.Response.Request.URL), r.Response.StatusCode, r.ErrorElement)
	}
}

func (w *ErrorWrapper) UnmarshalJSON(data []byte) error {
	s := string(data)
	if strings.HasPrefix(s, "{") {
		var validations map[string]interface{}
		err := json.Unmarshal(data, &validations)
		if err != nil {
			return err
		}
		w.Validations = validations
		return nil
	}
	s = strings.Trim(s, "\"")
	w.Message = s
	w.Partial = true
	return nil
}
