package xelon

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKubernetes_UpgradeHighAvailability(t *testing.T) {
	setup(t)
	defer teardown()

	mux.HandleFunc("POST /kubernetes/kubernetes-cluster-1/upgrade", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.Empty(t, body)

		w.WriteHeader(http.StatusOK)
	})

	resp, err := client.Kubernetes.UpgradeHighAvailability(ctx, "kubernetes-cluster-1")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestKubernetes_UpgradeHighAvailability_EmptyKubernetesClusterID(t *testing.T) {
	c := newTestClient(t)
	_, err := c.Kubernetes.UpgradeHighAvailability(ctx, "")

	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrEmptyArgument))
}

func TestKubernetes_ScaleNodePool(t *testing.T) {
	tests := map[string]int{
		"scale down to zero": 0,
		"scale up":           5,
	}

	for name, desiredNodeCount := range tests {
		t.Run(name, func(t *testing.T) {
			setup(t)
			defer teardown()

			mux.HandleFunc("PUT /kubernetes/kubernetes-cluster-1/pools/node-pool-1/scale", func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPut, r.Method)

				var payload map[string]int
				err := json.NewDecoder(r.Body).Decode(&payload)
				assert.NoError(t, err)
				assert.Equal(t, map[string]int{"desiredNodeCount": desiredNodeCount}, payload)

				w.WriteHeader(http.StatusAccepted)
			})

			resp, err := client.Kubernetes.ScaleNodePool(ctx, "kubernetes-cluster-1", "node-pool-1", desiredNodeCount)

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, http.StatusAccepted, resp.StatusCode)
		})
	}
}

func TestKubernetes_ScaleNodePool_ValidationErrors(t *testing.T) {
	c := newTestClient(t)
	tests := map[string]error{
		"empty kubernetes cluster id": func() error {
			_, err := c.Kubernetes.ScaleNodePool(ctx, "", "node-pool-1", 3)
			return err
		}(),
		"empty node pool id": func() error {
			_, err := c.Kubernetes.ScaleNodePool(ctx, "kubernetes-cluster-1", "", 3)
			return err
		}(),
	}

	for name, err := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Error(t, err)
			assert.True(t, errors.Is(err, ErrEmptyArgument))
		})
	}
}
