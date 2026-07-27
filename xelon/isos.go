package xelon

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

const isoBasePath = "isos"

// ISOsService handles communication with the ISO related methods of the Xelon REST API.
type ISOsService service

// ISO represents a Xelon custom ISO.
type ISO struct {
	Active      bool       `json:"active,omitempty"`
	Category    string     `json:"category,omitempty"`
	Cloud       *Cloud     `json:"cloud,omitempty"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	Description string     `json:"description,omitempty"`
	ID          string     `json:"identifier,omitempty"`
	Name        string     `json:"name,omitempty"`
	Owner       string     `json:"owner,omitempty"`
	Status      bool       `json:"status,omitempty"`
}

type ISOCreateRequest struct {
	CategoryID  int    `json:"categoryId"`
	CloudID     string `json:"cloudIdentifier"`
	Description string `json:"description,omitempty"`
	Name        string `json:"name"`
	TenantID    string `json:"tenantIdentifier,omitempty"`
	URL         string `json:"url"`
}

// ISOUploadRequest specifies a local ISO and its upload metadata.
type ISOUploadRequest struct {
	CategoryID  int
	CloudID     string
	Description string
	File        io.Reader
	Filename    string
	Name        string
	TenantID    string
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

// Upload uploads a local ISO file to a cloud datastore.
func (s *ISOsService) Upload(ctx context.Context, uploadRequest *ISOUploadRequest) (*ISO, *Response, error) {
	if uploadRequest == nil {
		return nil, nil, fmt.Errorf("payload: %w", ErrEmptyPayloadNotAllowed)
	}
	if uploadRequest.File == nil {
		return nil, nil, fmt.Errorf("file: %w", ErrEmptyArgument)
	}
	if uploadRequest.Filename == "" {
		return nil, nil, fmt.Errorf("filename: %w", ErrEmptyArgument)
	}
	if uploadRequest.Name == "" {
		return nil, nil, fmt.Errorf("name: %w", ErrEmptyArgument)
	}
	if uploadRequest.CloudID == "" {
		return nil, nil, fmt.Errorf("cloud id: %w", ErrEmptyArgument)
	}
	if uploadRequest.CategoryID <= 0 {
		return nil, nil, fmt.Errorf("category id: %w", ErrEmptyArgument)
	}
	uploadContext, cancel := s.client.withFallbackTimeout(ctx, defaultISOUploadTimeout)
	if cancel != nil {
		defer cancel()
	}
	if err := uploadContext.Err(); err != nil {
		return nil, nil, err
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writeISOUploadMultipart(uploadContext, writer, uploadRequest); err != nil {
		_ = writer.Close()
		return nil, nil, err
	}
	contentType := writer.FormDataContentType()
	if err := writer.Close(); err != nil {
		return nil, nil, fmt.Errorf("close multipart body: %w", err)
	}

	req, err := s.client.newRequest(http.MethodPost, isoBasePath+"/upload", &body, contentType)
	if err != nil {
		return nil, nil, fmt.Errorf("create upload request: %w", err)
	}
	if req.ContentLength <= 0 {
		return nil, nil, errors.New("ISO upload request has no known positive content length")
	}

	root := new(isoRoot)
	resp, err := s.client.do(uploadContext, req, root, defaultISOUploadTimeout)
	if err != nil {
		return nil, resp, fmt.Errorf("upload ISO: %w", err)
	}
	if root.ISO == nil {
		return nil, resp, errors.New("iso data is empty")
	}

	return root.ISO, resp, nil
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

func writeISOUploadMultipart(ctx context.Context, writer *multipart.Writer, uploadRequest *ISOUploadRequest) error {
	fields := []struct {
		name  string
		value string
	}{
		{name: "name", value: uploadRequest.Name},
		{name: "categoryId", value: strconv.Itoa(uploadRequest.CategoryID)},
		{name: "cloudIdentifier", value: uploadRequest.CloudID},
	}
	if uploadRequest.Description != "" {
		fields = append(fields, struct {
			name  string
			value string
		}{name: "description", value: uploadRequest.Description})
	}
	if uploadRequest.TenantID != "" {
		fields = append(fields, struct {
			name  string
			value string
		}{name: "tenantIdentifier", value: uploadRequest.TenantID})
	}

	for _, field := range fields {
		if err := writer.WriteField(field.name, field.value); err != nil {
			return fmt.Errorf("write multipart field %q: %w", field.name, err)
		}
	}

	part, err := writer.CreateFormFile("file", filepath.Base(uploadRequest.Filename))
	if err != nil {
		return fmt.Errorf("create multipart file part: %w", err)
	}
	if _, err := io.Copy(part, contextReader{ctx: ctx, reader: uploadRequest.File}); err != nil {
		return fmt.Errorf("copy ISO file: %w", err)
	}

	return nil
}

type contextReader struct {
	ctx    context.Context
	reader io.Reader
}

func (r contextReader) Read(p []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}

	return r.reader.Read(p)
}
