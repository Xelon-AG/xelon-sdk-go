package xelon

import (
	"encoding/json"
	"io"
	"net/http"
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDevices_DeviceNetworkIPAddresses_UnmarshalJSON(t *testing.T) {
	type testCase struct {
		input     string
		expect    []netip.Addr
		expectErr bool
	}
	tests := map[string]testCase{
		"single string": {
			input:     `"10.0.0.1"`,
			expect:    []netip.Addr{netip.MustParseAddr("10.0.0.1")},
			expectErr: false,
		},
		"array of strings": {
			input: `["10.0.0.1", "2001:db8::1"]`,
			expect: []netip.Addr{
				netip.MustParseAddr("10.0.0.1"),
				netip.MustParseAddr("2001:db8::1"),
			},
			expectErr: false,
		},
		"empty array": {
			input:     `[]`,
			expect:    []netip.Addr{},
			expectErr: false,
		},
		"null becomes empty": {
			input:     `null`,
			expect:    []netip.Addr{},
			expectErr: false,
		},
		"invalid ip rejected": {
			input:     `"not-an-ip"`,
			expect:    nil,
			expectErr: true,
		},
		"one bad entry rejects whole list": {
			input:     `["10.0.0.1", "not-an-ip"]`,
			expect:    nil,
			expectErr: true,
		},
		"wrong type rejected": {
			input:     `42`,
			expect:    nil,
			expectErr: true,
		},
		"unspecified rejected": {
			input:     `"0.0.0.0"`,
			expect:    nil,
			expectErr: true,
		},
		"duplicates rejected": {
			input:     `["10.0.0.1", ""10.0.0.1"]`,
			expect:    nil,
			expectErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var actual DeviceNetworkIPAddresses
			err := json.Unmarshal([]byte(test.input), &actual)

			if test.expectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.ElementsMatch(t, test.expect, actual)
		})
	}
}

func TestDevices_ListSSHKeys(t *testing.T) {
	setup(t)
	defer teardown()

	mux.HandleFunc("GET /devices/device-1/ssh-key", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/devices/device-1/ssh-key", r.URL.Path)

		fixture := loadFixture(t, "devices_list_ssh_keys.json")
		_, _ = w.Write(fixture)
	})

	sshKeys, resp, err := client.Devices.ListSSHKeys(ctx, "device-1")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Nil(t, resp.Meta)
	assert.Equal(t, []SSHKey{{
		CreatedAt: mustTime(t, "2026-09-14T10:00:00Z"),
		ID:        "ssh-key-1",
		Name:      "first key",
		PublicKey: "ssh-ed25519 AAAA-first",
	}, {
		CreatedAt: mustTime(t, "2026-09-15T11:00:00Z"),
		ID:        "ssh-key-2",
		Name:      "second key",
		PublicKey: "ssh-ed25519 AAAA-second",
	}}, sshKeys)
}

func TestDevices_AddSSHKey(t *testing.T) {
	setup(t)
	defer teardown()

	mux.HandleFunc("POST /devices/device-1/ssh-key", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/devices/device-1/ssh-key", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.Equal(t, "{\"sshKeyId\":\"ssh-key-1\"}\n", string(body))

		_, _ = io.WriteString(w, `{"data":{"identifier":"device-1","displayName":"device"},"message":"Key successfully added."}`)
	})

	resp, err := client.Devices.AddSSHKey(ctx, "device-1", "ssh-key-1")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDevices_DeviceCreateRequest_SSHKeyIDs(t *testing.T) {
	withKeys, err := json.Marshal(&DeviceCreateRequest{SSHKeyIDs: []string{"ssh-key-1", "ssh-key-2"}})
	assert.NoError(t, err)
	assert.Contains(t, string(withKeys), `"sshKeyIds":["ssh-key-1","ssh-key-2"]`)

	withoutKeys, err := json.Marshal(&DeviceCreateRequest{})
	assert.NoError(t, err)
	assert.NotContains(t, string(withoutKeys), "sshKeyIds")
}

func TestDevices_RemoveSSHKey(t *testing.T) {
	setup(t)
	defer teardown()

	mux.HandleFunc("DELETE /devices/device-1/ssh-key/ssh-key-1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/devices/device-1/ssh-key/ssh-key-1", r.URL.Path)

		_, _ = io.WriteString(w, `{"message":"SSH key removal has been started."}`)
	})

	resp, err := client.Devices.RemoveSSHKey(ctx, "device-1", "ssh-key-1")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDevices_SSHKeyValidationErrors(t *testing.T) {
	c := newTestClient(t)

	tests := map[string]struct {
		err    error
		target error
	}{
		"list empty device id": {
			err: func() error {
				_, _, err := c.Devices.ListSSHKeys(ctx, "")
				return err
			}(),
			target: ErrEmptyArgument,
		},
		"add empty device id": {
			err: func() error {
				_, err := c.Devices.AddSSHKey(ctx, "", "ssh-key-1")
				return err
			}(),
			target: ErrEmptyArgument,
		},
		"add empty ssh key id": {
			err: func() error {
				_, err := c.Devices.AddSSHKey(ctx, "device-1", "")
				return err
			}(),
			target: ErrEmptyArgument,
		},
		"remove empty device id": {
			err: func() error {
				_, err := c.Devices.RemoveSSHKey(ctx, "", "ssh-key-1")
				return err
			}(),
			target: ErrEmptyArgument,
		},
		"remove empty ssh key id": {
			err: func() error {
				_, err := c.Devices.RemoveSSHKey(ctx, "device-1", "")
				return err
			}(),
			target: ErrEmptyArgument,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Error(t, test.err)
			assert.ErrorIs(t, test.err, test.target)
		})
	}
}
