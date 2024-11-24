package xelon

import (
	"context"
	"fmt"
	"net/http"
)

const sshBasePathV1 = "vmlist/ssh"

// SSHKeysServiceV1 handles communication with the ssh keys related methods of the Xelon API.
// Deprecated.
type SSHKeysServiceV1 service

// SSHKeyV1 represents a Xelon ssh key.
// Deprecated.
type SSHKeyV1 struct {
	CreatedAt   string `json:"created_at,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	PublicKey   string `json:"ssh_key,omitempty"`
}

// SSHKeyCreateRequestV1 represents a request to create an ssh key.
// Deprecated.
type SSHKeyCreateRequestV1 struct {
	*SSHKeyV1
}

func (v SSHKeyV1) String() string {
	return Stringify(v)
}

// List provides a list of all added SSH keys.
func (s *SSHKeysServiceV1) List(ctx context.Context) ([]SSHKeyV1, *Response, error) {
	path := "sshKeys/"

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var sshKeys []SSHKeyV1
	resp, err := s.client.Do(ctx, req, &sshKeys)
	if err != nil {
		return nil, resp, err
	}

	return sshKeys, resp, nil
}

// Create makes a new ssh key with given payload.
func (s *SSHKeysServiceV1) Create(ctx context.Context, createRequest *SSHKeyCreateRequestV1) (*SSHKeyV1, *Response, error) {
	if createRequest == nil {
		return nil, nil, ErrEmptyPayloadNotAllowed
	}

	path := fmt.Sprintf("%v/add", sshBasePathV1)

	req, err := s.client.NewRequest(http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	sshKey := new(SSHKeyV1)
	resp, err := s.client.Do(ctx, req, sshKey)
	if err != nil {
		return nil, resp, err
	}

	return sshKey, resp, nil
}

// Delete removes a ssh key identified by id.
func (s *SSHKeysServiceV1) Delete(ctx context.Context, sshKeyID int) (*Response, error) {
	path := fmt.Sprintf("%v/%v/delete", sshBasePathV1, sshKeyID)
	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
