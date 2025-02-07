package xelon

import (
	"context"
	"fmt"
	"net/http"
)

const tenantBasePath = "tenants"

// TenantsService handles communication with the user related methods of the Xelon REST API.
type TenantsService service

// Tenant represents a top-level entity in the Xelon cloud.
type Tenant struct {
	ID     string `json:"identifier,omitempty"`
	Name   string `json:"name,omitempty"`
	Parent string `json:"parent,omitempty"`
	Status string `json:"status,omitempty"`
	Type   string `json:"type,omitempty"`
}

// TenantListOptions specifies the optional parameters to the TenantsService.List.
type TenantListOptions struct {
	Sort   string `url:"sort,omitempty"`
	Search string `url:"search,omitempty"`

	ListOptions
}

type tenantsRoot struct {
	Tenants []Tenant `json:"data"`
	Meta    *Meta    `json:"meta,omitempty"`
}

func (v Tenant) String() string { return Stringify(v) }

// GetCurrent provides detailed information for current tenant.
func (s *TenantsService) GetCurrent(ctx context.Context) (*Tenant, *Response, error) {
	path := fmt.Sprintf("%v/current", tenantBasePath)
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

// List providers a list of all tenants.
func (s *TenantsService) List(ctx context.Context, opts *TenantListOptions) ([]Tenant, *Response, error) {
	path, err := addOptions(tenantBasePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(tenantsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root.Tenants, resp, nil
}
