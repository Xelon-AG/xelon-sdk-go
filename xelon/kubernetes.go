package xelon

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"iter"
	"net/http"
	"net/netip"
	"time"
)

const kubernetesBasePath = "kubernetes"

// KubernetesService handles communication with the Kubernetes related methods of the Xelon API.
type KubernetesService service

// KubernetesCluster represents a Xelon Kubernetes cluster.
type KubernetesCluster struct {
	Cloud     *Cloud                   `json:"cloud,omitempty"`
	CreatedAt *time.Time               `json:"createdAt,omitempty"`
	Health    *KubernetesClusterHealth `json:"health,omitempty"`
	ID        string                   `json:"identifier,omitempty"`
	Name      string                   `json:"name,omitempty"`
	Status    string                   `json:"status,omitempty"`
}

type KubernetesClusterHealth struct {
	Status           string `json:"health,omitempty"`
	LastCheckingData string `json:"lastCheckingData,omitempty"`
}

type KubernetesClusterCreateRequest struct {
	CloudID              string                                     `json:"cloudIdentifier"`
	ControlPlaneCPUCores int                                        `json:"controlPlaneCpu"`
	ControlPlaneDiskSize int                                        `json:"controlPlaneDisk"`
	ControlPlaneRAM      int                                        `json:"controlPlaneRam"`
	ControlPlaneType     string                                     `json:"controlPlaneType"`
	KubernetesVersion    string                                     `json:"k8sVersion"`
	LoadBalancerCPUCores int                                        `json:"loadBalancerCpu"`
	LoadBalancerDiskSize int                                        `json:"loadBalancerDisk"`
	LoadBalancerRAM      int                                        `json:"loadBalancerRam"`
	LoadBalancerType     string                                     `json:"loadBalancerType"`
	Name                 string                                     `json:"clusterName"`
	PodCIDRBlock         string                                     `json:"podSubnet"`
	SendEmail            bool                                       `json:"notify"`
	ServiceCIDRBlock     string                                     `json:"serviceSubnet"`
	TalosVersion         string                                     `json:"talosVersion"`
	TenantID             string                                     `json:"tenantIdentifier"`
	WorkerPools          []KubernetesClusterCreateRequestWorkerPool `json:"workerPool"`
}

func (r KubernetesClusterCreateRequest) MarshalJSON() ([]byte, error) {
	type alias KubernetesClusterCreateRequest
	return json.Marshal(alias{
		CloudID:              r.CloudID,
		ControlPlaneCPUCores: r.ControlPlaneCPUCores,
		ControlPlaneDiskSize: r.ControlPlaneDiskSize,
		ControlPlaneRAM:      r.ControlPlaneRAM,
		ControlPlaneType:     r.ControlPlaneType,
		KubernetesVersion:    r.KubernetesVersion,
		LoadBalancerCPUCores: r.LoadBalancerCPUCores,
		LoadBalancerDiskSize: r.LoadBalancerDiskSize,
		LoadBalancerRAM:      r.LoadBalancerRAM,
		LoadBalancerType:     r.LoadBalancerType,
		Name:                 r.Name,
		PodCIDRBlock:         r.PodCIDRBlock,
		ServiceCIDRBlock:     r.ServiceCIDRBlock,
		TalosVersion:         r.TalosVersion,
		TenantID:             r.TenantID,
		WorkerPools:          nilToEmpty(r.WorkerPools),
	})
}

type KubernetesClusterCreateRequestWorkerPool struct {
	ExtraStorageEnabled  bool   `json:"workerNodeIsStorage,omitempty"`
	ExtraStorageDiskSize int    `json:"workerNodeExtraDisk,omitempty"`
	Index                string `json:"workerPoolIndex"`
	Name                 string `json:"workerPoolName"`
	NodeCount            int    `json:"workerNodeAmount"`
	NodeCPUCores         int    `json:"workerNodeCpu"`
	NodeDiskSize         int    `json:"workerNodeDisk"`
	NodeRAM              int    `json:"workerNodeRam"`
}

type kubernetesClusterRoot struct {
	KubernetesCluster *KubernetesCluster `json:"data,omitempty"`
	Message           string             `json:"message,omitempty"`
}

type kubernetesClustersRoot struct {
	KubernetesClusters []KubernetesCluster `json:"data"`
	Meta               *Meta               `json:"meta,omitempty"`
}

func (v KubernetesCluster) String() string { return Stringify(v) }

// List provides a list of available Kubernetes clusters.
func (s *KubernetesService) List(ctx context.Context, opts *ListOptions) ([]KubernetesCluster, *Response, error) {
	path, err := addOptions(kubernetesBasePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(kubernetesClustersRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root.KubernetesClusters, resp, nil
}

// All returns an iterator to paginate over all available Kubernetes clusters.
//
// The return iterator can be used in a for...range loop to easily process all Kubernetes clusters.
func (s *KubernetesService) All(ctx context.Context, opts *ListOptions) (iter.Seq2[KubernetesCluster, *Response], func() error) {
	return newPaginator[KubernetesCluster](ctx, s.client, kubernetesBasePath, opts)
}

// Get provides detailed information for Kubernetes cluster identified by id.
func (s *KubernetesService) Get(ctx context.Context, kubernetesClusterID string) (*KubernetesCluster, *Response, error) {
	if kubernetesClusterID == "" {
		return nil, nil, errors.New("failed to get kubernetes cluster: id must be supplied")
	}

	path := fmt.Sprintf("%v/%v", kubernetesBasePath, kubernetesClusterID)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(kubernetesClusterRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.KubernetesCluster, resp, nil
}

// Create makes a new Kubernetes cluster with given payload.
func (s *KubernetesService) Create(ctx context.Context, createRequest *KubernetesClusterCreateRequest) (*KubernetesCluster, *Response, error) {
	if createRequest == nil {
		return nil, nil, errors.New("failed to create kubernetes cluster: payload must be supplied")
	}

	req, err := s.client.NewRequest(http.MethodPost, kubernetesBasePath, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(kubernetesClusterRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.KubernetesCluster, resp, nil
}

// Delete removes Kubernetes cluster identified by id.
func (s *KubernetesService) Delete(ctx context.Context, kubernetesClusterID string) (*Response, error) {
	if kubernetesClusterID == "" {
		return nil, errors.New("failed to delete kubernetes cluster: id must be supplied")
	}

	path := fmt.Sprintf("%v/%v", kubernetesBasePath, kubernetesClusterID)
	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// GetKubeConfig returns the raw Kubernetes config in YAML format.
func (s *KubernetesService) GetKubeConfig(ctx context.Context, kubernetesClusterID string) ([]byte, *Response, error) {
	if kubernetesClusterID == "" {
		return nil, nil, errors.New("failed to get kube config: id must be supplied")
	}
	path := fmt.Sprintf("%v/%v/config/kube", kubernetesBasePath, kubernetesClusterID)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/yaml")

	var buf bytes.Buffer
	resp, err := s.client.Do(ctx, req, &buf)
	if err != nil {
		return nil, resp, err
	}

	return buf.Bytes(), resp, nil
}

type KubernetesClusterControlPlane struct {
	CPUCores int                     `json:"controlPlaneCpu,omitempty"`
	DiskSize int                     `json:"controlPlaneDisk,omitempty"`
	RAM      int                     `json:"controlPlaneRam,omitempty"`
	Nodes    []KubernetesClusterNode `json:"nodes,omitempty"`
}

type KubernetesClusterNodePool struct {
	CPUCores             int                     `json:"cpu,omitempty"`
	DiskSize             int                     `json:"disk,omitempty"`
	ExtraStorageDiskSize int                     `json:"extraStorage,omitempty"`
	ID                   string                  `json:"identifier,omitempty"`
	Name                 string                  `json:"name,omitempty"`
	Nodes                []KubernetesClusterNode `json:"nodes,omitempty"`
	RAM                  int                     `json:"memory,omitempty"`
}

type KubernetesClusterNode struct {
	ID        string `json:"identifier,omitempty"`
	LocalVMID string `json:"localvmid,omitempty"`
	Name      string `json:"name,omitempty"`
	Status    string `json:"status,omitempty"`
}

type KubernetesClusterNodePoolCreateRequest struct {
	WorkerPoolName      string `json:"workerPoolName"`
	WorkerNodeAmount    int    `json:"workerNodeAmount"`
	WorkerNodeCpu       int    `json:"workerNodeCpu"`
	WorkerNodeRam       int    `json:"workerNodeRam"`
	WorkerNodeDisk      int    `json:"workerNodeDisk"`
	WorkerNodeIsStorage bool   `json:"workerNodeIsStorage"`
	WorkerNodeExtraDisk int    `json:"workerNodeExtraDisk,omitempty"`
}

type kubernetesClusterNodePoolRoot struct {
	KubernetesClusterNodePool *KubernetesClusterNodePool `json:"data,omitempty"`
	Message                   string                     `json:"message,omitempty"`
}

func (v KubernetesClusterControlPlane) String() string { return Stringify(v) }
func (v KubernetesClusterNodePool) String() string     { return Stringify(v) }
func (v KubernetesClusterNode) String() string         { return Stringify(v) }

// ListControlPlane provides information about control planes on Kubernetes cluster.
func (s *KubernetesService) ListControlPlane(ctx context.Context, kubernetesClusterID string) (*KubernetesClusterControlPlane, *Response, error) {
	if kubernetesClusterID == "" {
		return nil, nil, errors.New("failed to list control plane: kubernetes cluster id must be supplied")
	}

	path := fmt.Sprintf("%v/%v/control-planes", kubernetesBasePath, kubernetesClusterID)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	controlPlane := new(KubernetesClusterControlPlane)
	resp, err := s.client.Do(ctx, req, controlPlane)
	if err != nil {
		return nil, resp, err
	}

	return controlPlane, resp, nil
}

// ListNodePools provides information about nodes pools on Kubernetes cluster.
func (s *KubernetesService) ListNodePools(ctx context.Context, kubernetesClusterID string) ([]KubernetesClusterNodePool, *Response, error) {
	if kubernetesClusterID == "" {
		return nil, nil, errors.New("failed to list nodes pools: kubernetes cluster id must be supplied")
	}

	path := fmt.Sprintf("%v/%v/pools", kubernetesBasePath, kubernetesClusterID)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var nodePools []KubernetesClusterNodePool
	resp, err := s.client.Do(ctx, req, &nodePools)
	if err != nil {
		return nil, resp, err
	}

	return nodePools, resp, nil
}

// CreateNodePool makes a nodes pool on Kubernetes cluster with given payload.
func (s *KubernetesService) CreateNodePool(ctx context.Context, kubernetesClusterID string, createRequest *KubernetesClusterNodePoolCreateRequest) (*KubernetesClusterNodePool, *Response, error) {
	if kubernetesClusterID == "" {
		return nil, nil, errors.New("failed to create nodes pool: kubernetes cluster id must be supplied")
	}
	if createRequest == nil {
		return nil, nil, errors.New("failed to create nodes pool: payload must be supplied")
	}

	path := fmt.Sprintf("%v/%v/pools", kubernetesBasePath, kubernetesClusterID)
	req, err := s.client.NewRequest(http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(kubernetesClusterNodePoolRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.KubernetesClusterNodePool, resp, nil
}

// DeleteNodePool removes the node pool.
func (s *KubernetesService) DeleteNodePool(ctx context.Context, kubernetesClusterID, nodePoolID string) (*Response, error) {
	if kubernetesClusterID == "" {
		return nil, errors.New("failed to delete nodes pool: kubernetes cluster id must be supplied")
	}
	if nodePoolID == "" {
		return nil, errors.New("failed to delete nodes pool: id must be supplied")
	}

	path := fmt.Sprintf("%v/%v/pools/%v", kubernetesBasePath, kubernetesClusterID, nodePoolID)
	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// CreateNode makes a new node in Kubernetes cluster.
func (s *KubernetesService) CreateNode(ctx context.Context, kubernetesClusterID, nodePoolID string) (*Response, error) {
	if kubernetesClusterID == "" {
		return nil, errors.New("failed to create node: kubernetes cluster id must be supplied")
	}
	if nodePoolID == "" {
		return nil, errors.New("failed to create node: node pool id must be supplied")
	}

	path := fmt.Sprintf("%v/%v/pools/%v/nodes", kubernetesBasePath, kubernetesClusterID, nodePoolID)
	req, err := s.client.NewRequest(http.MethodPost, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// DeleteNode removes node in Kubernetes cluster.
func (s *KubernetesService) DeleteNode(ctx context.Context, kubernetesClusterID, nodeID string) (*Response, error) {
	if kubernetesClusterID == "" {
		return nil, errors.New("failed to delete node: kubernetes cluster id must be supplied")
	}
	if nodeID == "" {
		return nil, errors.New("failed to delete node: id must be supplied")
	}

	path := fmt.Sprintf("%v/%v/nodes/%v", kubernetesBasePath, kubernetesClusterID, nodeID)
	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

type KubernetesClusterLoadBalancer struct {
	CPUCores  int                                     `json:"loadBalancerCpu,omitempty"`
	DiskSize  int                                     `json:"loadBalancerDisk,omitempty"`
	Name      string                                  `json:"clusterName,omitempty"`
	RAM       int                                     `json:"loadBalancerRam,omitempty"`
	Instances []KubernetesClusterLoadBalancerInstance `json:"loadBalancers,omitempty"`
}

type KubernetesClusterLoadBalancerInstance struct {
	Name      string     `json:"loadBalancerName,omitempty"`
	IPAddress netip.Addr `json:"loadBalancerIp,omitempty"`
}

func (v KubernetesClusterLoadBalancer) String() string         { return Stringify(v) }
func (v KubernetesClusterLoadBalancerInstance) String() string { return Stringify(v) }

// ListLoadBalancer provides information about load balancer cluster.
func (s *KubernetesService) ListLoadBalancer(ctx context.Context, kubernetesClusterID string) (*KubernetesClusterLoadBalancer, *Response, error) {
	if kubernetesClusterID == "" {
		return nil, nil, errors.New("failed to list load balancers: kubernetes cluster id must be supplied")
	}

	path := fmt.Sprintf("%v/%v/load-balancers", kubernetesBasePath, kubernetesClusterID)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var loadBalancer []KubernetesClusterLoadBalancer
	resp, err := s.client.Do(ctx, req, &loadBalancer)
	if err != nil {
		return nil, resp, err
	}
	if len(loadBalancer) > 0 {
		// backend returns load balancer as array with single element
		return &loadBalancer[0], resp, nil
	} else {
		return nil, resp, nil
	}
}

// KubernetesClusterVersionMapping maps a Talos version to a list of compatible
// Kubernetes versions for use in Kubernetes cluster provisioning and management.
//
// Example:
//
//	{
//	  "1.10.9": ["1.28.13", "1.29.10", ...],
//	  "1.11.6": ["1.29.10", "1.30.10", ...]
//	}
type KubernetesClusterVersionMapping map[string][]string

// ListVersionMapping retrieves the mapping of Talos versions to their compatible Kubernetes versions.
func (s *KubernetesService) ListVersionMapping(ctx context.Context, cloudID string) (KubernetesClusterVersionMapping, *Response, error) {
	if cloudID == "" {
		return nil, nil, errors.New("failed to list version mapping: cloud id must be supplied")
	}

	path := fmt.Sprintf("%v/versions/%v", kubernetesBasePath, cloudID)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var mapping KubernetesClusterVersionMapping
	resp, err := s.client.Do(ctx, req, &mapping)
	if err != nil {
		return nil, resp, err
	}

	return mapping, resp, nil
}
