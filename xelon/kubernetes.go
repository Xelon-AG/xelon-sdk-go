package xelon

import (
	"context"
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
	CloudID              string                                     `json:"cloud_id"`
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
	KubernetesClusters *KubernetesCluster `json:"data,omitempty"`
	Message            string             `json:"message,omitempty"`
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

// Create makes a new Kubernetes cluster with given payload.
func (s *KubernetesService) Create(ctx context.Context, createRequest *KubernetesClusterCreateRequest) (*Response, error) {
	if createRequest == nil {
		return nil, errors.New("failed to create kubernetes cluster: payload must be supplied")
	}

	req, err := s.client.NewRequest(http.MethodPost, kubernetesBasePath, createRequest)
	if err != nil {
		return nil, err
	}

	root := new(kubernetesClusterRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return resp, err
	}

	return resp, nil
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
