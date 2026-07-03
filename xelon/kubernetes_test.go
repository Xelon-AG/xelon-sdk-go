package xelon

import (
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKubernetes_UpgradeHighAvailability(t *testing.T) {
	setup()
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
	_, err := client.Kubernetes.UpgradeHighAvailability(ctx, "")

	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrEmptyArgument))
}
