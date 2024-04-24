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
	Cloud    *Cloud   `json:"hv_system,omitempty"`
	ID       string   `json:"identifier,omitempty"`
	Name     string   `json:"name,omitempty"`
	Nodes    []string `json:"nodes,omitempty"`
	Status   string   `json:"status,omitempty"`
	TenantID string   `json:"tenantIdentifier,omitempty"`
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
	Backend  *LoadBalancerClusterForwardingRuleConfiguration `json:"backend,omitempty"`
	Frontend *LoadBalancerClusterForwardingRuleConfiguration `json:"frontend,omitempty"`
}

type LoadBalancerClusterForwardingRuleConfiguration struct {
	ID   string `json:"identifier,omitempty"`
	Port int    `json:"port,omitempty"`
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
func (s *LoadBalancerClustersService) Get(ctx context.Context, id string) (*LoadBalancerCluster, *Response, error) {
	if id == "" {
		return nil, nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v", loadBalancerClusterBasePath, id)
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
func (s *LoadBalancerClustersService) Delete(ctx context.Context, id string) (*Response, error) {
	if id == "" {
		return nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v", loadBalancerClusterBasePath, id)
	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *LoadBalancerClustersService) ListVirtualIPs(ctx context.Context, id string) ([]LoadBalancerClusterVirtualIP, *Response, error) {
	if id == "" {
		return nil, nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v/virtual-ips", loadBalancerClusterBasePath, id)
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
