package xelon

import (
	"context"
	"fmt"
	"net/http"
)

const templatesBasePath = "templates"

// TemplatesServiceV1 handles communication with the template related methods of the Xelon API.
// Deprecated.
type TemplatesServiceV1 service

type Template struct {
	Description  string `json:"description,omitempty"`
	Category     string `json:"category,omitempty"`
	Cloud        *Cloud `json:"hv_system,omitempty"`
	ID           int    `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	NICUnit      int    `json:"nicunit,omitempty"`
	TemplateType int    `json:"templatetype,omitempty"`
	Type         string `json:"type,omitempty"`
}

type Templates struct {
	Firewalls []Template `json:"templates_firewalls,omitempty"`
	Linux     []Template `json:"templates_linux,omitempty"`
	Templates []Template `json:"templates,omitempty"`
	Windows   []Template `json:"templates_windows,omitempty"`
}

// List provides a list of available templates.
//
// Note, passing 0 (zero) for cloudID will retrieve templates across all available clouds.
func (s *TemplatesServiceV1) List(ctx context.Context, cloudID int) (*Templates, *Response, error) {
	path := fmt.Sprintf("device/%s", templatesBasePath)
	if cloudID != 0 {
		path = fmt.Sprintf("%v?cloudId=%v", path, cloudID)
	}

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
