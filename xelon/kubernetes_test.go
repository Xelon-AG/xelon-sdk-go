package xelon

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKubernetes_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/kubernetes-talos/clusters", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = fmt.Fprint(w, `[{"clusterIdentifier":"abc","name":"test cluster","status":"Ready"}]`)
	})
	expected := []KubernetesCluster{{
		ID:     "abc",
		Name:   "test cluster",
		Status: "Ready",
	}}

	clusters, _, err := client.Kubernetes.List(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expected, clusters)
}

func TestKubernetes_ListClusterPools(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/kubernetes-talos/abc/cluster-pools", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = fmt.Fprint(w, `
[{
  "cpu": 2,
  "disk": 50,
  "identifier": "abc",
  "memory": 4,
  "name": "test cluster pool",
  "nodes": [
    {"identifier":"def","name":"node-1","status":"Created"},
    {"identifier":"ghi","name":"node-2","status":"Deployed"}
  ]
}]`)
	})
	expected := []ClusterPool{{
		CPUCoreCount: 2,
		DiskSize:     50,
		ID:           "abc",
		Memory:       4,
		Name:         "test cluster pool",
		Nodes: []ClusterPoolNode{
			{ID: "def", Name: "node-1", Status: "Created"},
			{ID: "ghi", Name: "node-2", Status: "Deployed"},
		},
	}}

	clusterPools, _, err := client.Kubernetes.ListClusterPools(ctx, "abc")

	assert.NoError(t, err)
	assert.Equal(t, expected, clusterPools)
}

func TestKubernetes_ListClusterPools_emptyKubernetesClusterID(t *testing.T) {
	_, _, err := client.Kubernetes.ListClusterPools(ctx, "")

	assert.Error(t, err)
	assert.Equal(t, ErrEmptyArgument, err)
}

func TestKubernetes_AddClusterNode(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/kubernetes-talos/abc/add-node/def", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		_, _ = fmt.Fprint(w, `{"success":"Cluster node will be added shortly"}`)
	})
	expected := &SuccessResponse{
		Success: "Cluster node will be added shortly",
	}

	successResponse, _, err := client.Kubernetes.AddClusterNode(ctx, "abc", "def")

	assert.NoError(t, err)
	assert.Equal(t, expected, successResponse)
}

func TestKubernetes_AddClusterNode_emptyKubernetesClusterID(t *testing.T) {
	_, _, err := client.Kubernetes.AddClusterNode(ctx, "", "def")

	assert.Error(t, err)
	assert.Equal(t, ErrEmptyArgument, err)
}

func TestKubernetes_AddClusterNode_emptyClusterPoolID(t *testing.T) {
	_, _, err := client.Kubernetes.AddClusterNode(ctx, "abc", "")

	assert.Error(t, err)
	assert.Equal(t, ErrEmptyArgument, err)
}

func TestKubernetes_DeleteClusterNode(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/kubernetes-talos/abc/delete-node/def", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		_, _ = fmt.Fprint(w, `{"success":"Cluster node has been deleted."}`)
	})
	expected := &SuccessResponse{
		Success: "Cluster node has been deleted.",
	}

	successResponse, _, err := client.Kubernetes.DeleteClusterNode(ctx, "abc", "def")

	assert.NoError(t, err)
	assert.Equal(t, expected, successResponse)
}

func TestKubernetes_DeleteClusterNode_emptyKubernetesClusterID(t *testing.T) {
	_, _, err := client.Kubernetes.DeleteClusterNode(ctx, "", "def")

	assert.Error(t, err)
	assert.Equal(t, ErrEmptyArgument, err)
}

func TestKubernetes_DeleteClusterNode_emptyClusterNodeID(t *testing.T) {
	_, _, err := client.Kubernetes.DeleteClusterNode(ctx, "abc", "")

	assert.Error(t, err)
	assert.Equal(t, ErrEmptyArgument, err)
}
