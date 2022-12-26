package xelon

import (
	"context"
	"fmt"
	"net/http"
)

const templatesBasePath = "templates"

// TemplatesService handles communication with the template related methods of the Xelon API.
type TemplatesService service

type Template struct {
	ID           int    `json:"id,omitempty"`
	HVSystemID   int    `json:"hv_system_id"`
	Name         string `json:"name,omitempty"`
	NICUnit      int    `json:"nicunit,omitempty"`
	TemplateType int    `json:"templatetype,omitempty"`
	Type         string `json:"type,omitempty"`
}

type Templates struct {
	LinuxTemplates []Template `json:"templates_linux,omitempty"`
	Templates      []Template `json:"templates,omitempty"`
}

func (s *TemplatesService) List(ctx context.Context) (*Templates, *Response, error) {
	path := fmt.Sprintf("device/%s", templatesBasePath)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	templates := new(Templates)
	resp, err := s.client.Do(ctx, req, templates)
	if err != nil {
		return nil, resp, err
	}

	return templates, resp, nil
}
