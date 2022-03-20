package xelon

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSSHKeys_List(t *testing.T) {
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

func TestSSHKeys_Create(t *testing.T) {
	setup()
	defer teardown()

	input := &SSHKeyCreateRequest{
		SSHKey: &SSHKey{
			ID:        1,
			Name:      "test key",
			PublicKey: "public-key",
		}}
	mux.HandleFunc("/vmlist/ssh/add", func(w http.ResponseWriter, r *http.Request) {
		v := new(SSHKeyCreateRequest)
		_ = json.NewDecoder(r.Body).Decode(v)
		assert.Equal(t, input, v)
		assert.Equal(t, http.MethodPost, r.Method)
		_, _ = fmt.Fprint(w, `{"id":1,"name":"test key"}`)
	})
	expected := &SSHKey{
		ID:   1,
		Name: "test key",
	}

	sshKey, _, err := client.SSHKeys.Create(ctx, input)

	assert.NoError(t, err)
	assert.Equal(t, expected, sshKey)
}

func TestSSHKeys_Create_emptyPayload(t *testing.T) {
	_, _, err := client.SSHKeys.Create(ctx, nil)

	assert.Error(t, err)
	assert.Equal(t, ErrEmptyPayloadNotAllowed, err)
}

func TestSSHKeys_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/vmlist/ssh/0/delete", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
	})

	_, err := client.SSHKeys.Delete(ctx, 0)

	assert.NoError(t, err)
}
