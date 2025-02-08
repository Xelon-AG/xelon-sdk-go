package xelon

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrors_String_errorsWithErrorOnly(t *testing.T) {
	validationResponse := []byte(`
{
  "error": "Virtual machine is not found"
}`)

	var errorElement ErrorElement
	err := json.Unmarshal(validationResponse, &errorElement)

	assert.NoError(t, err)
	assert.Equal(t, "Virtual machine is not found", errorElement.Error)
	assert.Empty(t, errorElement.Message)
	assert.Equal(t, 0, len(errorElement.Validations))
}

func TestErrors_String_errorsWithMessage(t *testing.T) {
	validationResponse := []byte(`
	{
		"message": "UNAUTHORIZED_REQUEST"
	}`)

	var errorElement ErrorElement
	err := json.Unmarshal(validationResponse, &errorElement)

	assert.NoError(t, err)
	assert.Equal(t, "UNAUTHORIZED_REQUEST", errorElement.Message)
	assert.Empty(t, errorElement.Error)
	assert.Equal(t, 0, len(errorElement.Validations))
}

func TestErrors_String_errorsWithValidations(t *testing.T) {
	validationResponse := []byte(`
{
  "errors": {
    "isLocked": [
      "Your network is currently locked. Please contact our support to remove it"
    ],
    "name": [
      "The field cannot be modified"
    ]
  }
}`)

	var errorElement ErrorElement
	err := json.Unmarshal(validationResponse, &errorElement)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(errorElement.Validations))
	assert.Empty(t, errorElement.Error)
	assert.Empty(t, errorElement.Message)
}

func TestErrors_String_errorsWithFullPayload(t *testing.T) {
	validationResponse := []byte(`
{
  "error": "Virtual machine is not found",
  "message": "Server Error",
  "errors": {
    "isLocked": [
      "Your network is currently locked. Please contact our support to remove it."
    ],
    "name": [
      "The field cannot be modified."
    ]
  }
}`)

	var errorsElement ErrorElement
	err := json.Unmarshal(validationResponse, &errorsElement)

	assert.NoError(t, err)
	assert.Equal(t, "Virtual machine is not found", errorsElement.Error)
	assert.Equal(t, "Server Error", errorsElement.Message)
	assert.Equal(t, 2, len(errorsElement.Validations))
}

func TestErrors_Error(t *testing.T) {
	errorResponse := &ErrorResponse{
		ErrorElement: ErrorElement{
			Error: "Virtual machine is not found",
		},
		Response: &Response{
			Response: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    &url.URL{Scheme: "http", Host: "localhost", Path: "api/testing"},
				},
			},
		},
	}

	actualErrorResponse := errorResponse.Error()
	expectedErrorResponse := "GET http://localhost/api/testing: 500 (error: Virtual machine is not found)"

	assert.Equal(t, expectedErrorResponse, actualErrorResponse)
}
