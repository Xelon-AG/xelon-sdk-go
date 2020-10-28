package xelon

import (
	"context"
	"fmt"
	"net/http"
)

const persistentStorageBasePath = "persistentStorage"

type PersistentStorageService service

type PersistentStorage struct {
	Capacity int    `json:"capacity,omitempty"`
	ID       int    `json:"id,omitempty"`
	LocalID  string `json:"local_id,omitempty"`
	Name     string `json:"name,omitempty"`
	Type     int    `json:"type,omitempty"`
}

type PersistentStorageCreateRequest struct {
	*PersistentStorage
	Size int `json:"size,omitempty"`
}

func (s *PersistentStorageService) List(ctx context.Context, tenantID string) ([]PersistentStorage, *http.Response, error) {
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

func (s *PersistentStorageService) Get(ctx context.Context, tenantID, localID string) (*PersistentStorage, *http.Response, error) {
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

func (s *PersistentStorageService) Create(ctx context.Context, tenantID string, createRequst *PersistentStorageCreateRequest) (*APIResponse, *http.Response, error) {
	if tenantID == "" {
		return nil, nil, ErrEmptyArgument
	}
	if createRequst == nil {
		return nil, nil, ErrEmptyPayloadNotAllowed
	}

	path := fmt.Sprintf("%v/%v", tenantID, persistentStorageBasePath)
	req, err := s.client.NewRequest(http.MethodPost, path, createRequst)
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

func (s *PersistentStorageService) Delete(ctx context.Context, tenantID, localID string) (*http.Response, error) {
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
