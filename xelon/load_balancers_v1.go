package xelon

import (
	"context"
	"fmt"
	"net/http"
)

// const loadBalancerBasePathV1 = "loadBalancer"

// LoadBalancersServiceV1 handles communication with the load balancer related methods of the Xelon API.
// Deprecated.
type LoadBalancersServiceV1 service

// LoadBalancerV1 represents a Xelon load balancer.
type LoadBalancerV1 struct {
	ForwardingRules []LoadBalancerForwardingRuleV1 `json:"forwarding_rules,omitempty"`
	Health          string                         `json:"health,omitempty"`
	HealthCheck     LoadBalancerHealthCheckV1      `json:"health_check,omitempty"`
	ID              int                            `json:"id,omitempty"`
	InternalIP      string                         `json:"internalIp,omitempty"`
	IP              string                         `json:"ip,omitempty"`
	LocalID         string                         `json:"local_id,omitempty"`
	Name            string                         `json:"name,omitempty"`
	Type            int                            `json:"type,omitempty"`
}

// LoadBalancerForwardingRuleV1 represents a Xelon load balancer forwarding rule.
type LoadBalancerForwardingRuleV1 struct {
	ID    int      `json:"id,omitempty"`
	IP    []string `json:"ip,omitempty"`
	Ports []int    `json:"ports,omitempty"`
	Type  string   `json:"type,omitempty"`
}

// LoadBalancerHealthCheckV1 represents a Xelon load balancer health check.
type LoadBalancerHealthCheckV1 struct {
	BadThreshold  int    `json:"bad_threshold,omitempty"`
	GoodThreshold int    `json:"good_threshold,omitempty"`
	Interval      int    `json:"interval,omitempty"`
	Path          string `json:"path,omitempty"`
	Port          int    `json:"port,omitempty"`
	Timeout       int    `json:"timeout,omitempty"`
}

type LoadBalancerCreateRequestV1 struct {
	CloudID         string                         `json:"cloudId,omitempty"`
	ForwardingRules []LoadBalancerForwardingRuleV1 `json:"forwarding_rules,omitempty"`
	Name            string                         `json:"name,omitempty"`
	ServerID        []string                       `json:"server_id,omitempty"`
	Type            int                            `json:"type,omitempty"`
}

type LoadBalancerUpdateForwardingRulesRequestV1 struct {
	ForwardingRules []LoadBalancerForwardingRuleV1 `json:"forwarding_rules,omitempty"`
}

// List provides information about load balancers.
func (s *LoadBalancersServiceV1) List(ctx context.Context, tenantID string) ([]LoadBalancerV1, *Response, error) {
	if tenantID == "" {
		return nil, nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v", tenantID, loadBalancerBasePath)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var loadBalancers []LoadBalancerV1
	resp, err := s.client.Do(ctx, req, &loadBalancers)
	if err != nil {
		return nil, resp, err
	}

	return loadBalancers, resp, nil
}

// Get provides information about a load balancer identified by local id.
func (s *LoadBalancersServiceV1) Get(ctx context.Context, tenantID, localID string) (*LoadBalancerV1, *Response, error) {
	if tenantID == "" || localID == "" {
		return nil, nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v/%v", tenantID, loadBalancerBasePath, localID)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	loadBalancer := new(LoadBalancerV1)
	resp, err := s.client.Do(ctx, req, loadBalancer)
	if err != nil {
		return nil, resp, err
	}

	return loadBalancer, resp, nil
}

// Create makes a new load balancer with given payload.
func (s *LoadBalancersServiceV1) Create(ctx context.Context, tenantID string, createRequest *LoadBalancerCreateRequestV1) (*APIResponse, *Response, error) {
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
func (s *LoadBalancersServiceV1) Delete(ctx context.Context, tenantID, localID string) (*Response, error) {
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

func (s *LoadBalancersServiceV1) UpdateForwardingRules(ctx context.Context, tenantID, localID string, updateRequest *LoadBalancerUpdateForwardingRulesRequestV1) (*APIResponse, *Response, error) {
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
