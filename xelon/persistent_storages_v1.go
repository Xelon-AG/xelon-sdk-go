package xelon

import (
	"context"
	"fmt"
	"net/http"
)

const persistentStorageBasePath = "persistentStorage"

// PersistentStoragesServiceV1 handles communication with the persistent storage related methods of the Xelon API.
// Deprecated.
type PersistentStoragesServiceV1 service

type PersistentStorage struct {
	AssignedServers []DeviceLocalVMDetails `json:"assigned_servers,omitempty"`
	Capacity        int                    `json:"capacity,omitempty"`
	Formatted       int                    `json:"formatted,omitempty"`
	ID              int                    `json:"id,omitempty"`
	LocalID         string                 `json:"local_id,omitempty"`
	Name            string                 `json:"name,omitempty"`
	Type            int                    `json:"type,omitempty"`
	UUID            string                 `json:"uuid,omitempty"`
}

type PersistentStorageCreateRequest struct {
	*PersistentStorage
	CloudID string `json:"cloudId,omitempty"`
	Size    int    `json:"size,omitempty"`
}

type PersistentStorageExtendRequest struct {
	Size int `json:"size"`
}

type PersistentStorageAttachDetachRequest struct {
	ServerID []string `json:"server_id"`
}

func (s *PersistentStoragesServiceV1) List(ctx context.Context, tenantID string) ([]PersistentStorage, *Response, error) {
	if tenantID == "" {
		return nil, nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v", tenantID, persistentStorageBasePath)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var persistentStorages []PersistentStorage
	resp, err := s.client.Do(ctx, req, &persistentStorages)
	if err != nil {
		return nil, resp, err
	}

	return persistentStorages, resp, nil
}

func (s *PersistentStoragesServiceV1) Get(ctx context.Context, tenantID, localID string) (*PersistentStorage, *Response, error) {
	if tenantID == "" || localID == "" {
		return nil, nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v/%v", tenantID, persistentStorageBasePath, localID)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	persistentStorage := new(PersistentStorage)
	resp, err := s.client.Do(ctx, req, persistentStorage)
	if err != nil {
		return nil, resp, err
	}

	return persistentStorage, resp, nil
}

func (s *PersistentStoragesServiceV1) GetByName(ctx context.Context, tenantID, name string) (*PersistentStorage, *Response, error) {
	if tenantID == "" || name == "" {
		return nil, nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v/query?name=%v", tenantID, persistentStorageBasePath, name)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	persistentStorage := new(PersistentStorage)
	resp, err := s.client.Do(ctx, req, persistentStorage)
	if err != nil {
		return nil, resp, err
	}

	return persistentStorage, resp, nil
}

func (s *PersistentStoragesServiceV1) Create(ctx context.Context, tenantID string, createRequest *PersistentStorageCreateRequest) (*APIResponse, *Response, error) {
	if tenantID == "" {
		return nil, nil, ErrEmptyArgument
	}
	if createRequest == nil {
		return nil, nil, ErrEmptyPayloadNotAllowed
	}

	path := fmt.Sprintf("%v/%v", tenantID, persistentStorageBasePath)
	req, err := s.client.NewRequest(http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	apiResponse := new(APIResponse)
	resp, err := s.client.Do(ctx, req, apiResponse)
	if err != nil {
		return nil, resp, err
	}

	return apiResponse, resp, nil
}

func (s *PersistentStoragesServiceV1) Extend(ctx context.Context, localID string, extendRequest *PersistentStorageExtendRequest) (*APIResponse, *Response, error) {
	if localID == "" {
		return nil, nil, ErrEmptyArgument
	}
	if extendRequest == nil {
		return nil, nil, ErrEmptyPayloadNotAllowed
	}

	path := fmt.Sprintf("%v/%v", persistentStorageBasePath, localID)
	req, err := s.client.NewRequest(http.MethodPost, path, extendRequest)
	if err != nil {
		return nil, nil, err
	}

	apiResponse := new(APIResponse)
	resp, err := s.client.Do(ctx, req, apiResponse)
	if err != nil {
		return nil, resp, err
	}

	return apiResponse, resp, nil
}

func (s *PersistentStoragesServiceV1) Delete(ctx context.Context, tenantID, localID string) (*Response, error) {
	if tenantID == "" || localID == "" {
		return nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v/%v", tenantID, persistentStorageBasePath, localID)
	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *PersistentStoragesServiceV1) AttachToDevice(ctx context.Context, tenantID, localID string, attachRequest *PersistentStorageAttachDetachRequest) (*APIResponse, *Response, error) {
	if tenantID == "" || localID == "" {
		return nil, nil, ErrEmptyArgument
	}
	if attachRequest == nil {
		return nil, nil, ErrEmptyPayloadNotAllowed
	}

	path := fmt.Sprintf("%v/%v/%v/addToVirtualMachine", tenantID, persistentStorageBasePath, localID)
	req, err := s.client.NewRequest(http.MethodPost, path, attachRequest)
	if err != nil {
		return nil, nil, err
	}

	apiResponse := new(APIResponse)
	resp, err := s.client.Do(ctx, req, apiResponse)
	if err != nil {
		return nil, resp, err
	}
	return apiResponse, resp, nil
}

func (s *PersistentStoragesServiceV1) DetachFromDevice(ctx context.Context, tenantID, localID string, detachRequest *PersistentStorageAttachDetachRequest) (*APIResponse, *Response, error) {
	if tenantID == "" || localID == "" {
		return nil, nil, ErrEmptyArgument
	}
	if detachRequest == nil {
		return nil, nil, ErrEmptyPayloadNotAllowed
	}

	path := fmt.Sprintf("%v/%v/%v/removeFromVirtualMachine", tenantID, persistentStorageBasePath, localID)
	req, err := s.client.NewRequest(http.MethodPost, path, detachRequest)
	if err != nil {
		return nil, nil, err
	}

	apiResponse := new(APIResponse)
	resp, err := s.client.Do(ctx, req, apiResponse)
	if err != nil {
		return nil, resp, err
	}
	return apiResponse, resp, nil
}
