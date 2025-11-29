package xelon

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

const templatesBasePath = "templates"

// TemplatesService handles communication with the template related methods of the Xelon API.
type TemplatesService service

// Template represents a Xelon base image.
type Template struct {
	Description string `json:"description,omitempty"`
	Category    string `json:"category,omitempty"`
	CloudID     string `json:"cloudIdentifier,omitempty"`
	ID          string `json:"identifier,omitempty"`
	Name        string `json:"name,omitempty"`
	Status      int    `json:"status,omitempty"`
	Type        string `json:"type,omitempty"`
}

type TemplateCreateRequest struct {
	Description     string `json:"description,omitempty"`
	DeviceID        string `json:"deviceId"`
	Name            string `json:"name"`
	SendEmail       bool   `json:"sendEmail"`
	TemplateOwnerID string `json:"ownerTenantId,omitempty"`
	TenantID        string `json:"tenantId"`
}

type TemplateUpdateRequest struct {
	Description string `json:"description,omitempty"`
	Name        string `json:"name"`
}

// TemplateListOptions specifies the optional parameters to the TemplatesService.List.
type TemplateListOptions struct {
	Sort   string `url:"sort,omitempty"`
	Search string `url:"search,omitempty"`
	Type   string `url:"type,omitempty"`

	ListOptions
}

type templateRoot struct {
	Template *Template `json:"data,omitempty"`
	Message  string    `json:"message,omitempty"`
}

type templatesRoot struct {
	Templates []Template `json:"data"`
	Meta      *Meta      `json:"meta,omitempty"`
}

func (v Template) String() string { return Stringify(v) }

// List provides a list of available templates.
func (s *TemplatesService) List(ctx context.Context, opts *TenantListOptions) ([]Template, *Response, error) {
	path, err := addOptions(templatesBasePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(templatesRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root.Templates, resp, nil
}

// Get provides detailed information for template identified by id.
func (s *TemplatesService) Get(ctx context.Context, templateID string) (*Template, *Response, error) {
	if templateID == "" {
		return nil, nil, errors.New("failed to get template: id must be supplied")
	}

	path := fmt.Sprintf("%v/%v", templatesBasePath, templateID)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	template := new(Template)
	resp, err := s.client.Do(ctx, req, template)
	if err != nil {
		return nil, resp, err
	}

	return template, resp, err
}

// Create makes a template with given payload.
func (s *TemplatesService) Create(ctx context.Context, createRequest *TemplateCreateRequest) (*Template, *Response, error) {
	if createRequest == nil {
		return nil, nil, errors.New("failed to create template: payload must be supplied")
	}

	path := fmt.Sprintf("%v/create-from-device", templatesBasePath)
	req, err := s.client.NewRequest(http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	templateRoot := new(templateRoot)
	resp, err := s.client.Do(ctx, req, templateRoot)
	if err != nil {
		return nil, resp, err
	}

	return templateRoot.Template, resp, nil
}

// Update changes template identified by id.
func (s *TemplatesService) Update(ctx context.Context, templateID string, updateRequest *TemplateUpdateRequest) (*Template, *Response, error) {
	if templateID == "" {
		return nil, nil, errors.New("failed to update template: id must be supplied")
	}
	if updateRequest == nil {
		return nil, nil, errors.New("failed to update template: payload must be supplied")
	}

	path := fmt.Sprintf("%v/%v", templatesBasePath, templateID)
	req, err := s.client.NewRequest(http.MethodPatch, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	templateRoot := new(templateRoot)
	resp, err := s.client.Do(ctx, req, templateRoot)
	if err != nil {
		return nil, resp, err
	}

	return templateRoot.Template, resp, nil
}

// Delete removes template identified by id.
func (s *TemplatesService) Delete(ctx context.Context, templateID string) (*Response, error) {
	if templateID == "" {
		return nil, errors.New("failed to delete template: id must be supplied")
	}

	path := fmt.Sprintf("%v/%v", templatesBasePath, templateID)
	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
