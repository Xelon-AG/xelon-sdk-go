package xelon

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

const isoBasePath = "isos"

// ISOsService handles communication with the ISO related methods of the Xelon REST API.
type ISOsService service

// ISO represents a Xelon custom ISO.
type ISO struct {
	Active      bool   `json:"active,omitempty"`
	Category    string `json:"category,omitempty"`
	Cloud       *Cloud `json:"cloud,omitempty"`
	Description string `json:"description,omitempty"`
	ID          string `json:"identifier,omitempty"`
	Name        string `json:"name,omitempty"`
	Owner       string `json:"owner,omitempty"`
	Status      bool   `json:"status,omitempty"`
}

type ISOCreateRequest struct {
	CategoryID  int    `json:"categoryId"`
	CloudID     string `json:"cloudIdentifier"`
	Description string `json:"description,omitempty"`
	Name        string `json:"name"`
	TenantID    string `json:"tenantIdentifier,omitempty"`
	URL         string `json:"url"`
}

type ISOUpdateRequest struct {
	CategoryID  int    `json:"categoryId"`
	Description string `json:"description"`
	Name        string `json:"name"`
}

// ISOListOptions specifies the optional parameters to the ISOsService.List.
type ISOListOptions struct {
	Sort   string `url:"sort,omitempty"`
	Search string `url:"search,omitempty"`

	ListOptions
}

type isoRoot struct {
	ISO     *ISO   `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

type isosRoot struct {
	ISOs []ISO `json:"data"`
	Meta *Meta `json:"meta,omitempty"`
}

func (v ISO) String() string { return Stringify(v) }

// List provides a list of all custom ISOs.
func (s *ISOsService) List(ctx context.Context, opts *ISOListOptions) ([]ISO, *Response, error) {
	path, err := addOptions(isoBasePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(isosRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root.ISOs, resp, nil
}

// Get provides detailed information for custom ISO identified by id.
func (s *ISOsService) Get(ctx context.Context, isoID string) (*ISO, *Response, error) {
	if isoID == "" {
		return nil, nil, errors.New("failed to get iso: id must be supplied")
	}

	path := fmt.Sprintf("%v/%v", isoBasePath, isoID)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	iso := new(ISO)
	resp, err := s.client.Do(ctx, req, iso)
	if err != nil {
		return nil, resp, err
	}

	return iso, resp, err
}

// Create makes a new custom ISO with given payload.
func (s *ISOsService) Create(ctx context.Context, createRequest *ISOCreateRequest) (*ISO, *Response, error) {
	if createRequest == nil {
		return nil, nil, errors.New("failed to create iso: payload must be supplied")
	}

	req, err := s.client.NewRequest(http.MethodPost, isoBasePath, createRequest)
	if err != nil {
		return nil, nil, err
	}

	isoRoot := new(isoRoot)
	resp, err := s.client.Do(ctx, req, isoRoot)
	if err != nil {
		return nil, resp, err
	}

	return isoRoot.ISO, resp, nil
}

// Update changes custom ISO identified by id.
func (s *ISOsService) Update(ctx context.Context, isoID string, updateRequest *ISOUpdateRequest) (*ISO, *Response, error) {
	if isoID == "" {
		return nil, nil, errors.New("failed to update iso: id must be supplied")
	}
	if updateRequest == nil {
		return nil, nil, errors.New("failed to update iso: payload must be supplied")
	}

	path := fmt.Sprintf("%v/%v", isoBasePath, isoID)
	req, err := s.client.NewRequest(http.MethodPatch, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	isoRoot := new(isoRoot)
	resp, err := s.client.Do(ctx, req, isoRoot)
	if err != nil {
		return nil, resp, err
	}

	return isoRoot.ISO, resp, nil
}

// Delete removes custom ISO identified by id.
func (s *ISOsService) Delete(ctx context.Context, isoID string) (*Response, error) {
	if isoID == "" {
		return nil, errors.New("failed to delete iso: id must be supplied")
	}

	path := fmt.Sprintf("%v/%v", isoBasePath, isoID)
	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
