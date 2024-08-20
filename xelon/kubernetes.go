package xelon

import (
	"context"
	"fmt"
	"net/http"
)

const kubernetesBasePath = "kubernetes-talos"

// KubernetesService handles communication with the Kubernetes
// related methods of the Xelon API.
type KubernetesService service

// KubernetesCluster represents a Xelon Kubernetes cluster.
type KubernetesCluster struct {
	Cloud     *Cloud                   `json:"hv_system,omitempty"`
	CreatedAt string                   `json:"createdAt,omitempty"`
	Health    *KubernetesClusterHealth `json:"health,omitempty"`
	ID        string                   `json:"clusterIdentifier,omitempty"`
	Name      string                   `json:"name,omitempty"`
	Status    string                   `json:"status,omitempty"`
}

type KubernetesClusterHealth struct {
	Health           string `json:"health,omitempty"`
	LastCheckingData string `json:"lastCheckingData,omitempty"`
}

type ClusterControlPlane struct {
	CPUCoreCount int                       `json:"control_plane_cpu,omitempty"`
	DiskSize     int                       `json:"control_plane_disk,omitempty"`
	Memory       int                       `json:"control_plane_ram,omitempty"`
	Nodes        []ClusterControlPlaneNode `json:"nodes,omitempty"`
}

type ClusterControlPlaneNode struct {
	ID        string `json:"identifier,omitempty"`
	LocalVMID string `json:"localvmid,omitempty"`
	Name      string `json:"name,omitempty"`
}

type ClusterPool struct {
	CPUCoreCount int               `json:"cpu,omitempty"`
	DiskSize     int               `json:"disk,omitempty"`
	ID           string            `json:"identifier,omitempty"`
	Memory       int               `json:"memory,omitempty"`
	Name         string            `json:"name,omitempty"`
	Nodes        []ClusterPoolNode `json:"nodes,omitempty"`
}

type ClusterPoolNode struct {
	ID        string `json:"identifier,omitempty"`
	LocalVMID string `json:"localvmid,omitempty"`
	Name      string `json:"name,omitempty"`
	Status    string `json:"status,omitempty"`
}

func (v KubernetesCluster) String() string {
	return Stringify(v)
}

func (v ClusterControlPlane) String() string {
	return Stringify(v)
}

func (v ClusterPool) String() string {
	return Stringify(v)
}

// List provides information about Kubernetes clusters.
func (s *KubernetesService) List(ctx context.Context) ([]KubernetesCluster, *Response, error) {
	path := fmt.Sprintf("%v/clusters", kubernetesBasePath)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var kubernetesClusters []KubernetesCluster
	resp, err := s.client.Do(ctx, req, &kubernetesClusters)
	if err != nil {
		return nil, resp, err
	}

	return kubernetesClusters, resp, nil
}

// ListControlPlanes provides information about control plane on Kubernetes cluster.
func (s *KubernetesService) ListControlPlanes(ctx context.Context, kubernetesClusterID string) (*ClusterControlPlane, *Response, error) {
	if kubernetesClusterID == "" {
		return nil, nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v/cluster-control-planes", kubernetesBasePath, kubernetesClusterID)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	clusterControlPlane := new(ClusterControlPlane)
	resp, err := s.client.Do(ctx, req, clusterControlPlane)
	if err != nil {
		return nil, resp, err
	}

	return clusterControlPlane, resp, nil
}

// ListClusterPools provides information about cluster pools on Kubernetes cluster.
func (s *KubernetesService) ListClusterPools(ctx context.Context, kubernetesClusterID string) ([]ClusterPool, *Response, error) {
	if kubernetesClusterID == "" {
		return nil, nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v/cluster-pools", kubernetesBasePath, kubernetesClusterID)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var clusterPools []ClusterPool
	resp, err := s.client.Do(ctx, req, &clusterPools)
	if err != nil {
		return nil, resp, err
	}

	return clusterPools, resp, nil
}

// AddClusterNode creates and adds a new cluster node in the specified cluster pool.
func (s *KubernetesService) AddClusterNode(ctx context.Context, kubernetesClusterID, clusterPoolID string) (*SuccessResponse, *Response, error) {
	if kubernetesClusterID == "" || clusterPoolID == "" {
		return nil, nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v/add-node/%v", kubernetesBasePath, kubernetesClusterID, clusterPoolID)
	req, err := s.client.NewRequest(http.MethodPost, path, nil)
	if err != nil {
		return nil, nil, err
	}

	successResponse := new(SuccessResponse)
	resp, err := s.client.Do(ctx, req, successResponse)
	if err != nil {
		return nil, resp, err
	}

	return successResponse, resp, nil
}

// DeleteClusterNode removes a cluster node.
func (s *KubernetesService) DeleteClusterNode(ctx context.Context, kubernetesClusterID, clusterNodeID string) (*SuccessResponse, *Response, error) {
	if kubernetesClusterID == "" || clusterNodeID == "" {
		return nil, nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v/delete-node/%v", kubernetesBasePath, kubernetesClusterID, clusterNodeID)
	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, nil, err
	}

	successResponse := new(SuccessResponse)
	resp, err := s.client.Do(ctx, req, successResponse)
	if err != nil {
		return nil, resp, err
	}

	return successResponse, resp, nil
}
