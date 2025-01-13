package xelon

import (
	"context"
	"net/http"
)

const cloudsBasePath = "clouds"

// CloudsService handles communication with the organization's cloud related methods of the Xelon REST API.
type CloudsService service

type Cloud struct {
	ID     string `json:"identifier,omitempty"`
	Name   string `json:"name,omitempty"`
	Type   int    `json:"type,omitempty"`
	HVType string `json:"hvType,omitempty"`
}

// CloudListOptions specifies the optional parameters to the CloudsService.List.
type CloudListOptions struct {
	TenantID string `url:"tenantIdentifier,omitempty"`
}

func (v Cloud) String() string { return Stringify(v) }

func (s *CloudsService) List(ctx context.Context, opts *CloudListOptions) ([]Cloud, *Response, error) {
	path, err := addOptions(cloudsBasePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var clouds []Cloud
	resp, err := s.client.Do(ctx, req, &clouds)
	if err != nil {
		return nil, resp, err
	}

	return clouds, resp, nil
}
