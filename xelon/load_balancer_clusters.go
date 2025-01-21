package xelon

import (
	"context"
	"fmt"
	"net/http"
)

const loadBalancerClusterBasePath = "load-balancer-clusters"

// LoadBalancerClustersService handles communication with the load balancer cluster
// related methods of the Xelon API.
type LoadBalancerClustersService service

// LoadBalancerCluster represents a Xelon load balancer cluster.
type LoadBalancerCluster struct {
	Cloud               *Cloud   `json:"hv_system,omitempty"`
	ID                  string   `json:"identifier,omitempty"`
	KubernetesClusterID string   `json:"kubernetesClusterIdentifier,omitempty"`
	Name                string   `json:"name,omitempty"`
	Nodes               []string `json:"nodes,omitempty"`
	Status              string   `json:"status,omitempty"`
	TenantID            string   `json:"tenantIdentifier,omitempty"`
}

type LoadBalancerClusterCreateRequest struct {
	CloudID                     int                           `json:"cloudId"`
	KubernetesClusterIdentifier string                        `json:"kubernetesClusterIdentifier"`
	Name                        string                        `json:"name"`
	NodesSpec                   *LoadBalancerClusterNodesSpec `json:"nodesSpec"`
}

type LoadBalancerClusterCreateResponse struct {
	LoadBalancerClusterID string `json:"identifier"`
	Status                string `json:"status"`
}

type LoadBalancerClusterNodesSpec struct {
	CPUCoreCount int `json:"cpuCoreCount"`
	Disk         int `json:"disk"`
	Memory       int `json:"memory"`
}

type LoadBalancerClusterVirtualIP struct {
	ID             string `json:"identifier,omitempty"`
	IPAddress      string `json:"ipAddress,omitempty"`
	PoolIdentifier string `json:"vipPoolIdentifier,omitempty"`
	State          string `json:"state,omitempty"`
}

type LoadBalancerClusterForwardingRule struct {
	Backend  *LoadBalancerClusterForwardingRuleBackendConfiguration  `json:"backend,omitempty"`
	Frontend *LoadBalancerClusterForwardingRuleFrontendConfiguration `json:"frontend,omitempty"`
}

type LoadBalancerClusterForwardingRuleBackendConfiguration struct {
	ID            string `json:"identifier,omitempty"`
	Port          int    `json:"port,omitempty"`
	ProxyProtocol int    `json:"proxy_protocol"`
}

type LoadBalancerClusterForwardingRuleFrontendConfiguration struct {
	ID   string `json:"identifier,omitempty"`
	Port int    `json:"port,omitempty"`
}

type LoadBalancerClusterForwardingRuleUpdateResponse struct {
	Port          int `json:"port,omitempty"`
	ProxyProtocol int `json:"proxy_protocol,omitempty"`
}

// List provides information about load balancer clusters.
func (s *LoadBalancerClustersService) List(ctx context.Context) ([]LoadBalancerCluster, *Response, error) {
	path := loadBalancerClusterBasePath
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var loadBalancerClusters []LoadBalancerCluster
	resp, err := s.client.Do(ctx, req, &loadBalancerClusters)
	if err != nil {
		return nil, resp, err
	}

	return loadBalancerClusters, resp, nil
}

// Get provides information about a load balancer cluster identified by id.
func (s *LoadBalancerClustersService) Get(ctx context.Context, loadBalancerClusterID string) (*LoadBalancerCluster, *Response, error) {
	if loadBalancerClusterID == "" {
		return nil, nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v", loadBalancerClusterBasePath, loadBalancerClusterID)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	loadBalancerCluster := new(LoadBalancerCluster)
	resp, err := s.client.Do(ctx, req, loadBalancerCluster)
	if err != nil {
		return nil, resp, err
	}

	return loadBalancerCluster, resp, nil
}

// Create makes a new load balancer cluster.
func (s *LoadBalancerClustersService) Create(ctx context.Context, createRequest *LoadBalancerClusterCreateRequest) (*LoadBalancerClusterCreateResponse, *Response, error) {
	if createRequest == nil {
		return nil, nil, ErrEmptyPayloadNotAllowed
	}

	path := loadBalancerClusterBasePath
	req, err := s.client.NewRequest(http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	apiResponse := new(LoadBalancerClusterCreateResponse)
	resp, err := s.client.Do(ctx, req, apiResponse)
	if err != nil {
		return nil, resp, err
	}

	return apiResponse, resp, nil
}

// Delete removes a load balancer.
func (s *LoadBalancerClustersService) Delete(ctx context.Context, loadBalancerClusterID string) (*Response, error) {
	if loadBalancerClusterID == "" {
		return nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v", loadBalancerClusterBasePath, loadBalancerClusterID)
	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// ListVirtualIPs provides information about virtual IP addresses.
func (s *LoadBalancerClustersService) ListVirtualIPs(ctx context.Context, loadBalancerClusterID string) ([]LoadBalancerClusterVirtualIP, *Response, error) {
	if loadBalancerClusterID == "" {
		return nil, nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v/virtual-ips", loadBalancerClusterBasePath, loadBalancerClusterID)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var virtualIPs []LoadBalancerClusterVirtualIP
	resp, err := s.client.Do(ctx, req, &virtualIPs)
	if err != nil {
		return nil, resp, err
	}

	return virtualIPs, resp, nil
}

// GetVirtualIP provides information about a virtual IP address identified by id.
func (s *LoadBalancerClustersService) GetVirtualIP(ctx context.Context, loadBalancerClusterID, virtualIPID string) (*LoadBalancerClusterVirtualIP, *Response, error) {
	if loadBalancerClusterID == "" || virtualIPID == "" {
		return nil, nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v/virtual-ips/%v", loadBalancerClusterBasePath, loadBalancerClusterID, virtualIPID)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	virtualIP := new(LoadBalancerClusterVirtualIP)
	resp, err := s.client.Do(ctx, req, virtualIP)
	if err != nil {
		return nil, resp, err
	}

	return virtualIP, resp, nil
}

// ListForwardingRules provides information about forwarding rules.
func (s *LoadBalancerClustersService) ListForwardingRules(ctx context.Context, loadBalancerClusterID, virtualIPID string) ([]LoadBalancerClusterForwardingRule, *Response, error) {
	if loadBalancerClusterID == "" || virtualIPID == "" {
		return nil, nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v/virtual-ips/%v/forwarding-rules", loadBalancerClusterBasePath, loadBalancerClusterID, virtualIPID)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var forwardingRules []LoadBalancerClusterForwardingRule
	resp, err := s.client.Do(ctx, req, &forwardingRules)
	if err != nil {
		return nil, resp, err
	}

	return forwardingRules, resp, nil
}

// CreateForwardingRules makes a new forwarding rule.
func (s *LoadBalancerClustersService) CreateForwardingRules(ctx context.Context, loadBalancerClusterID, virtualIPID string, createRequest []LoadBalancerClusterForwardingRule) ([]LoadBalancerClusterForwardingRule, *Response, error) {
	if loadBalancerClusterID == "" || virtualIPID == "" {
		return nil, nil, ErrEmptyArgument
	}
	if createRequest == nil {
		return nil, nil, ErrEmptyPayloadNotAllowed
	}

	path := fmt.Sprintf("%v/%v/virtual-ips/%v/forwarding-rules", loadBalancerClusterBasePath, loadBalancerClusterID, virtualIPID)
	req, err := s.client.NewRequest(http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	var forwardingRules []LoadBalancerClusterForwardingRule
	resp, err := s.client.Do(ctx, req, &forwardingRules)
	if err != nil {
		return nil, resp, err
	}

	return forwardingRules, resp, nil
}

// UpdateForwardingRule changes the configuration of a forwarding rule.
func (s *LoadBalancerClustersService) UpdateForwardingRule(ctx context.Context, loadBalancerClusterID, virtualIPID, forwardingRuleID string, updateRequest *LoadBalancerClusterForwardingRuleUpdateResponse) (*APIResponse, *Response, error) {
	if loadBalancerClusterID == "" || virtualIPID == "" || forwardingRuleID == "" {
		return nil, nil, ErrEmptyArgument
	}
	if updateRequest == nil {
		return nil, nil, ErrEmptyPayloadNotAllowed
	}

	path := fmt.Sprintf("%v/%v/virtual-ips/%v/forwarding-rules/%v", loadBalancerClusterBasePath, loadBalancerClusterID, virtualIPID, forwardingRuleID)
	req, err := s.client.NewRequest(http.MethodPatch, path, updateRequest)
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

// DeleteForwardingRule removes a forwarding rule.
func (s *LoadBalancerClustersService) DeleteForwardingRule(ctx context.Context, loadBalancerClusterID, virtualIPID, forwardingRuleID string) (*Response, error) {
	if loadBalancerClusterID == "" || virtualIPID == "" || forwardingRuleID == "" {
		return nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v/virtual-ips/%v/forwarding-rules/%v", loadBalancerClusterBasePath, loadBalancerClusterID, virtualIPID, forwardingRuleID)
	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
