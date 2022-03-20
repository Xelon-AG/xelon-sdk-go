package xelon

import (
	"context"
	"net/http"
)

const tenantBasePath = "tenant"

// TenantService handles communication with the user related methods of the Xelon API.
type TenantService service

type Tenant struct {
	TenantID string `json:"tenant_identifier"`
}

// Get provides information about user especially tenant id.
func (s *TenantService) Get(ctx context.Context) (*Tenant, *Response, error) {
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
