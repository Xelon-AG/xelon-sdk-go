package xelon

import (
	"context"
	"fmt"
	"net/http"
)

const cloudsBasePath = "hv"

// CloudsServiceV1 handles communication with the organization's cloud related methods of the Xelon API.
// Deprecated.
type CloudsServiceV1 service

type Cloud struct {
	ID        int    `json:"id,omitempty"`
	Name      string `json:"display_name,omitempty"`
	ShortName string `json:"display_short_name,omitempty"`
	Type      int    `json:"type,omitempty"`
}

func (s *CloudsServiceV1) List(ctx context.Context, tenantID string) ([]Cloud, *Response, error) {
	if tenantID == "" {
		return nil, nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/list/%v", cloudsBasePath, tenantID)

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
