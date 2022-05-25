package xelon

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrors_UnmarshalJSON_errorAsString(t *testing.T) {
	validationResponse := []byte(`
	{
		"error": "UNAUTHORIZED_REQUEST"
	}`)

	var errorElement ErrorElement
	err := json.Unmarshal(validationResponse, &errorElement)

	assert.NoError(t, err)
	assert.True(t, errorElement.Error.Partial)
	assert.Equal(t, "UNAUTHORIZED_REQUEST", errorElement.Error.Message)
}

func TestErrors_UnmarshalJSON_errorsAsValidationArray(t *testing.T) {
	validationResponse := []byte(`
	{
	  "error": {
	    "name": [
	      "The name field is required."
	    ],
	    "ssh_key": [
	      "SSH key is not valid"
	    ]
	  }
	}`)

	var errorElement ErrorElement
	err := json.Unmarshal(validationResponse, &errorElement)

	assert.NoError(t, err)
	assert.False(t, errorElement.Error.Partial)
	assert.Equal(t, 2, len(errorElement.Error.Validations))
}
