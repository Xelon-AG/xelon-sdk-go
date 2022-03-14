package xelon

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSSHKeys_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/sshKeys/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = fmt.Fprint(w, `[{"id":1,"name":"test key"}]`)
	})
	expected := []SSHKey{
		{
			ID:   1,
			Name: "test key",
		},
	}

	sshKeys, _, err := client.SSHKeys.List(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expected, sshKeys)
}
