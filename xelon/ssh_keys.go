package xelon

import (
	"context"
	"fmt"
	"net/http"
)

const sshBasePath = "vmlist/ssh"

// SSHKeysService handles communication with the ssh keys related methods of the Xelon API.
type SSHKeysService service

// SSHKey represents a Xelon ssh key.
type SSHKey struct {
	CreatedAt   string `json:"created_at,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	PublicKey   string `json:"ssh_key,omitempty"`
}

// SSHKeyCreateRequest represents a request to create a ssh key.
type SSHKeyCreateRequest struct {
	*SSHKey
}

func (v SSHKey) String() string {
	return Stringify(v)
}

// List provides a list of all added SSH keys.
func (s *SSHKeysService) List(ctx context.Context) ([]SSHKey, *Response, error) {
	path := "sshKeys/"

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var sshKeys []SSHKey
	resp, err := s.client.Do(ctx, req, &sshKeys)
	if err != nil {
		return nil, resp, err
	}

	return sshKeys, resp, nil
}

// Create makes a new ssh key with given payload.
func (s *SSHKeysService) Create(ctx context.Context, createRequest *SSHKeyCreateRequest) (*SSHKey, *Response, error) {
	if createRequest == nil {
		return nil, nil, ErrEmptyPayloadNotAllowed
	}

	path := fmt.Sprintf("%v/add", sshBasePath)

	req, err := s.client.NewRequest(http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	sshKey := new(SSHKey)
	resp, err := s.client.Do(ctx, req, sshKey)
	if err != nil {
		return nil, resp, err
	}

	return sshKey, resp, nil
}

// Delete removes a ssh key identified by id.
func (s *SSHKeysService) Delete(ctx context.Context, sshKeyID int) (*Response, error) {
	path := fmt.Sprintf("%v/%v/delete", sshBasePath, sshKeyID)
	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
