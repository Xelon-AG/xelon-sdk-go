package xelon

import (
	"context"
	"net/http"
)

const tenantBasePath = "tenant"

// TenantsServiceV1 handles communication with the user related methods of the Xelon API.
// Deprecated.
type TenantsServiceV1 service

type Tenant struct {
	Active   bool   `json:"active,omitempty"`
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Parent   string `json:"parent,omitempty"`
	TenantID string `json:"tenant_identifier,omitempty"`
}

// GetCurrent provides information about organization.
//
// Note, after calling this method only TenantID field is filled.
func (s *TenantsServiceV1) GetCurrent(ctx context.Context) (*Tenant, *Response, error) {
	path := tenantBasePath

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	tenant := new(Tenant)
	resp, err := s.client.Do(ctx, req, tenant)
	if err != nil {
		return nil, resp, err
	}

	return tenant, resp, nil
}

// List provides information about tenant (aka organizations).
func (s *TenantsServiceV1) List(ctx context.Context) ([]Tenant, *Response, error) {
	path := "tenants"

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var tenants []Tenant
	resp, err := s.client.Do(ctx, req, &tenants)
	if err != nil {
		return nil, resp, err
	}

	return tenants, resp, nil
}
