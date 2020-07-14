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
	ForwardingRules []LoadBalancerForwardingRule `json:"forwarding_rules,omitempty"`
	Health          string                       `json:"health,omitempty"`
	HealthCheck     LoadBalancerHealthCheck      `json:"health_check,omitempty"`
	ID              int                          `json:"id,omitempty"`
	InternalIP      string                       `json:"internalIp,omitempty"`
	IP              string                       `json:"ip,omitempty"`
	LocalID         string                       `json:"local_id,omitempty"`
	Name            string                       `json:"name,omitempty"`
	Type            int                          `json:"type,omitempty"`
}

// LoadBalancerForwardingRule represents a Xelon load balancer forwarding rule.
type LoadBalancerForwardingRule struct {
	IP    string `json:"ip,omitempty"`
	Ports []int  `json:"ports,omitempty"`
}

// LoadBalancerHealthCheck represents a Xelon load balancer health check.
type LoadBalancerHealthCheck struct {
	BadThreshold  int    `json:"bad_threshold,omitempty"`
	GoodThreshold int    `json:"good_threshold,omitempty"`
	Interval      int    `json:"interval,omitempty"`
	Path          string `json:"path,omitempty"`
	Port          int    `json:"port,omitempty"`
	Timeout       int    `json:"timeout,omitempty"`
}

type LoadBalancerCreateRequest struct {
	Name     string   `json:"name,omitempty"`
	ServerID []string `json:"server_id,omitempty"`
	Type     int      `json:"type,omitempty"`
}

type LoadBalancerUpdateForwardingRulesRequest struct {
	ForwardingRules []LoadBalancerForwardingRule `json:"forwarding_rules,omitempty"`
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

// Create makes a new load balancer with given payload.
func (s *LoadBalancerService) Create(ctx context.Context, tenantID string, createRequest *LoadBalancerCreateRequest) (*APIResponse, *http.Response, error) {
	if tenantID == "" {
		return nil, nil, ErrEmptyArgument
	}
	if createRequest == nil {
		return nil, nil, ErrEmptyPayloadNotAllowed
	}

	path := fmt.Sprintf("%v/%v", tenantID, loadBalancerBasePath)
	req, err := s.client.NewRequest(http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	apiResponse := new(APIResponse)
	resp, err := s.client.Do(ctx, req, apiResponse)
	if err != nil {
		return nil, resp, err
	}

	return apiResponse, resp, nil
}

// Delete removes a load balancer.
func (s *LoadBalancerService) Delete(ctx context.Context, tenantID, localID string) (*http.Response, error) {
	if tenantID == "" || localID == "" {
		return nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v/%v", tenantID, loadBalancerBasePath, localID)
	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *LoadBalancerService) UpdateForwardingRules(ctx context.Context, tenantID, localID string, updateRequest *LoadBalancerUpdateForwardingRulesRequest) (*APIResponse, *http.Response, error) {
	if tenantID == "" || localID == "" {
		return nil, nil, ErrEmptyArgument
	}
	if updateRequest == nil {
		return nil, nil, ErrEmptyPayloadNotAllowed
	}

	path := fmt.Sprintf("%v/%v/%v/forwardingRules", tenantID, loadBalancerBasePath, localID)
	req, err := s.client.NewRequest(http.MethodPut, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	apiResponse := new(APIResponse)
	resp, err := s.client.Do(ctx, req, apiResponse)
	if err != nil {
		return nil, resp, err
	}

	return apiResponse, resp, nil
}
