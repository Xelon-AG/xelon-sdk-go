package xelon

import (
	"context"
	"fmt"
	"net/http"
)

const loadBalancerBasePath = "loadBalancer"

// LoadBalancerService handles communication with the load balancer related methods of the Xelon API.
type LoadBalancerService service

// LoadBalancer represents a Xelon load balancer.
type LoadBalancer struct {
	ID      int    `json:"id,omitempty"`
	LocalID string `json:"local_id,omitempty"`
	Name    string `json:"name,omitempty"`
	Type    int    `json:"type,omitempty"`
}

// List provides information about load balancers.
func (s *LoadBalancerService) List(ctx context.Context, tenantID string) ([]LoadBalancer, *http.Response, error) {
	if tenantID == "" {
		return nil, nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v", tenantID, loadBalancerBasePath)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var loadBalancers []LoadBalancer
	resp, err := s.client.Do(ctx, req, &loadBalancers)
	if err != nil {
		return nil, resp, err
	}

	return loadBalancers, resp, nil
}

// Get provides information about a load balancer identified by local id.
func (s *LoadBalancerService) Get(ctx context.Context, tenantID, localID string) (*LoadBalancer, *http.Response, error) {
	if tenantID == "" || localID == "" {
		return nil, nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v/%v", tenantID, loadBalancerBasePath, localID)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	loadBalancer := new(LoadBalancer)
	resp, err := s.client.Do(ctx, req, loadBalancer)
	if err != nil {
		return nil, resp, err
	}

	return loadBalancer, resp, nil
}
