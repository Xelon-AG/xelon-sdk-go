package xelon

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"iter"
	"net/http"
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
	CloudID              string                                     `json:"cloud_identifier"`
	ControlPlaneCPUCores int                                        `json:"control_plane_cpu"`
	ControlPlaneDiskSize int                                        `json:"control_plane_disk"`
	ControlPlaneRAM      int                                        `json:"control_plane_ram"`
	ControlPlaneType     string                                     `json:"control_plane_type"`
	KubernetesVersion    string                                     `json:"k8s_version"`
	LoadBalancerCPUCores int                                        `json:"load_balancer_cpu"`
	LoadBalancerDiskSize int                                        `json:"load_balancer_disk"`
	LoadBalancerRAM      int                                        `json:"load_balancer_ram"`
	LoadBalancerType     string                                     `json:"load_balancer_type"`
	Name                 string                                     `json:"cluster_name"`
	PodCIDRBlock         string                                     `json:"pod_subnet"`
	ServiceCIDRBlock     string                                     `json:"service_subnet"`
	TalosVersion         string                                     `json:"talos_version"`
	TenantID             string                                     `json:"tenant_identifier"`
	WorkerPools          []KubernetesClusterCreateRequestWorkerPool `json:"worker_pool"`
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
	ExtraStorageEnabled  bool   `json:"worker_node_is_storage,omitempty"`
	ExtraStorageDiskSize int    `json:"worker_node_extra_disk,omitempty"`
	Index                string `json:"worker_pool_index"`
	Name                 string `json:"worker_pool_name"`
	NodeCount            int    `json:"worker_node_amount"`
	NodeCPUCores         int    `json:"worker_node_cpu"`
	NodeDiskSize         int    `json:"worker_node_disk"`
	NodeRAM              int    `json:"worker_node_ram"`
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
