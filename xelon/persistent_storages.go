package xelon

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

const persistentStorageBasePath = "persistent-storages"

// PersistentStoragesService handles communication with the persistent storage related methods of the Xelon REST API.
type PersistentStoragesService service

type PersistentStorage struct {
	AttachedDevices []PersistentStorageAttachedDevice `json:"attachedDevices,omitempty"`
	Capacity        int                               `json:"capacity,omitempty"`
	Cloud           *Cloud                            `json:"cloud,omitempty"`
	Formatted       bool                              `json:"formatted,omitempty"`
	ID              string                            `json:"identifier,omitempty"`
	Name            string                            `json:"name,omitempty"`
	Tenant          *Tenant                           `json:"tenant,omitempty"`
	Type            int                               `json:"type,omitempty"`
	UUID            string                            `json:"uuid,omitempty"`
}

type PersistentStorageAttachedDevice struct {
	ID   string `json:"identifier,omitempty"`
	Name string `json:"name,omitempty"`
}

type PersistentStorageCreateRequest struct {
	CloudID  string `json:"cloudIdentifier,omitempty"`
	DeviceID string `json:"deviceIdentifier,omitempty"`
	Name     string `json:"name"`
	Size     int    `json:"storageSize"`
	TenantID string `json:"tenantIdentifier,omitempty"`
	Type     int    `json:"type"`
}

type persistentStorageAttachDetachDeviceRequest struct {
	DeviceID string `json:"deviceIdentifier"`
}

type persistentStorageExtendRequest struct {
	Capacity int `json:"diskSize"`
}

// PersistentStorageListOptions specifies the optional parameters to the PersistentStoragesService.List.
type PersistentStorageListOptions struct {
	Sort   string `url:"sort,omitempty"`
	Search string `url:"search,omitempty"`

	ListOptions
}

type persistentStorageRoot struct {
	PersistentStorage *PersistentStorage `json:"data,omitempty"`
	Message           string             `json:"message,omitempty"`
}

type persistentStoragesRoot struct {
	PersistentStorages []PersistentStorage `json:"data"`
	Meta               *Meta               `json:"meta,omitempty"`
}

func (v PersistentStorage) String() string { return Stringify(v) }

// List provides a list of all persistent storages.
func (s *PersistentStoragesService) List(ctx context.Context, opts *PersistentStorageListOptions) ([]PersistentStorage, *Response, error) {
	path, err := addOptions(persistentStorageBasePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(persistentStoragesRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root.PersistentStorages, resp, nil
}

// Get provides detailed information for persistent storage identified by id.
func (s *PersistentStoragesService) Get(ctx context.Context, persistentStorageID string) (*PersistentStorage, *Response, error) {
	if persistentStorageID == "" {
		return nil, nil, errors.New("failed to get persistent storage: id must be supplied")
	}

	path := fmt.Sprintf("%v/%v", persistentStorageBasePath, persistentStorageID)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	persistentStorage := new(PersistentStorage)
	resp, err := s.client.Do(ctx, req, persistentStorage)
	if err != nil {
		return nil, resp, err
	}

	return persistentStorage, resp, err
}

// Create makes a new persistent storage with given payload.
func (s *PersistentStoragesService) Create(ctx context.Context, createRequest *PersistentStorageCreateRequest) (*PersistentStorage, *Response, error) {
	if createRequest == nil {
		return nil, nil, errors.New("failed to create persistent storage: payload must be supplied")
	}

	req, err := s.client.NewRequest(http.MethodPost, persistentStorageBasePath, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(persistentStorageRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.PersistentStorage, resp, nil
}

// Delete removes persistent storage identified by id.
func (s *PersistentStoragesService) Delete(ctx context.Context, persistentStorageID string) (*Response, error) {
	if persistentStorageID == "" {
		return nil, errors.New("failed to delete persistent storage: id must be supplied")
	}

	path := fmt.Sprintf("%v/%v", persistentStorageBasePath, persistentStorageID)
	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// AttachToDevice connects persistent storage to specified device.
func (s *PersistentStoragesService) AttachToDevice(ctx context.Context, persistentStorageID, deviceID string) (*Response, error) {
	if persistentStorageID == "" {
		return nil, errors.New("failed to attach persistent storage: id must be supplied")
	}
	if deviceID == "" {
		return nil, errors.New("failed to attach persistent storage: device id must be supplied")
	}
	attachDeviceRequest := &persistentStorageAttachDetachDeviceRequest{DeviceID: deviceID}

	path := fmt.Sprintf("%v/%v/attach-device", persistentStorageBasePath, persistentStorageID)
	req, err := s.client.NewRequest(http.MethodPost, path, attachDeviceRequest)
	if err != nil {
		return nil, err
	}

	root := new(persistentStorageRoot)
	return s.client.Do(ctx, req, root)
}

// DetachFromDevice disconnects persistent storage from specified device.
func (s *PersistentStoragesService) DetachFromDevice(ctx context.Context, persistentStorageID, deviceID string) (*Response, error) {
	if persistentStorageID == "" {
		return nil, errors.New("failed to detach persistent storage: id must be supplied")
	}
	if deviceID == "" {
		return nil, errors.New("failed to detach persistent storage: device id must be supplied")
	}
	detachDeviceRequest := &persistentStorageAttachDetachDeviceRequest{DeviceID: deviceID}

	path := fmt.Sprintf("%v/%v/detach-device", persistentStorageBasePath, persistentStorageID)
	req, err := s.client.NewRequest(http.MethodPost, path, detachDeviceRequest)
	if err != nil {
		return nil, err
	}

	root := new(persistentStorageRoot)
	return s.client.Do(ctx, req, root)
}

// Extend increases persistent storage capacity.
func (s *PersistentStoragesService) Extend(ctx context.Context, persistentStorageID string, capacity int) (*Response, error) {
	if persistentStorageID == "" {
		return nil, errors.New("failed to extend persistent storage: id must be supplied")
	}
	extendRequest := &persistentStorageExtendRequest{Capacity: capacity}

	path := fmt.Sprintf("%v/%v/extend", persistentStorageBasePath, persistentStorageID)
	req, err := s.client.NewRequest(http.MethodPost, path, extendRequest)
	if err != nil {
		return nil, err
	}

	root := new(persistentStorageRoot)
	return s.client.Do(ctx, req, root)
}
