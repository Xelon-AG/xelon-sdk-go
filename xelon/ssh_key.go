package xelon

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const sshBasePath = "ssh-keys"

// SSHKeysService handles communication with the SSH keys related methods of the Xelon REST API.
type SSHKeysService service

// SSHKey represents a Xelon SSH key.
type SSHKey struct {
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	ID        string     `json:"identifier,omitempty"`
	Name      string     `json:"name"`
	PublicKey string     `json:"sshKey"`
}

type SSHKeyCreateRequest struct {
	SSHKey
}

type SSHKeyUpdateRequest struct {
	SSHKey
}

// SSHKeyListOptions specifies the optional parameters to the SSHKeysService.List.
type SSHKeyListOptions struct {
	Sort   string `url:"sort,omitempty"`
	Search string `url:"search,omitempty"`

	ListOptions
}

type sshKeyRoot struct {
	SSHKey  *SSHKey `json:"data,omitempty"`
	Message string  `json:"message,omitempty"`
}

type sshKeysRoot struct {
	SSHKeys []SSHKey `json:"data"`
	Meta    *Meta    `json:"meta,omitempty"`
}

func (v SSHKey) String() string { return Stringify(v) }

// List provides a list of all SSH keys.
func (s *SSHKeysService) List(ctx context.Context, opts *SSHKeyListOptions) ([]SSHKey, *Response, error) {
	path, err := addOptions(sshBasePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(sshKeysRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root.SSHKeys, resp, nil
}

// Get provides detailed information for SSH key identified by id.
func (s *SSHKeysService) Get(ctx context.Context, sshKeyID string) (*SSHKey, *Response, error) {
	if sshKeyID == "" {
		return nil, nil, errors.New("failed to get ssh key: id must be supplied")
	}

	path := fmt.Sprintf("%v/%v", sshBasePath, sshKeyID)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	sshKey := new(SSHKey)
	resp, err := s.client.Do(ctx, req, sshKey)
	if err != nil {
		return nil, resp, err
	}

	return sshKey, resp, err
}

// Create makes a new SSH key with given payload.
func (s *SSHKeysService) Create(ctx context.Context, createRequest *SSHKeyCreateRequest) (*SSHKey, *Response, error) {
	if createRequest == nil {
		return nil, nil, errors.New("failed to create ssh key: payload must be supplied")
	}

	req, err := s.client.NewRequest(http.MethodPost, sshBasePath, createRequest)
	if err != nil {
		return nil, nil, err
	}

	sshKeyRoot := new(sshKeyRoot)
	resp, err := s.client.Do(ctx, req, sshKeyRoot)
	if err != nil {
		return nil, resp, err
	}

	return sshKeyRoot.SSHKey, resp, nil
}

// Update changes SSH key identified by id.
func (s *SSHKeysService) Update(ctx context.Context, sshKeyID string, updateRequest *SSHKeyUpdateRequest) (*SSHKey, *Response, error) {
	if sshKeyID == "" {
		return nil, nil, errors.New("failed to update ssh key: id must be supplied")
	}
	if updateRequest == nil {
		return nil, nil, errors.New("failed to update ssh key: payload must be supplied")
	}

	path := fmt.Sprintf("%v/%v", sshBasePath, sshKeyID)
	req, err := s.client.NewRequest(http.MethodPatch, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	sshKeyRoot := new(sshKeyRoot)
	resp, err := s.client.Do(ctx, req, sshKeyRoot)
	if err != nil {
		return nil, resp, err
	}

	return sshKeyRoot.SSHKey, resp, nil
}

// Delete removes SSH key identified by id.
func (s *SSHKeysService) Delete(ctx context.Context, sshKeyID string) (*Response, error) {
	if sshKeyID == "" {
		return nil, errors.New("failed to delete ssh key: id must be supplied")
	}

	path := fmt.Sprintf("%v/%v", sshBasePath, sshKeyID)
	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
