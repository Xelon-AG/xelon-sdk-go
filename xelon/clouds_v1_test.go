package xelon

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClouds_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/hv/list/12345", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = fmt.Fprint(w, `[{"id":1,"type":2,"display_name":"public cloud #1"},{"id":5,"type":1,"display_name":"private cloud #5"}]`)
	})
	expected := []Cloud{
		{
			ID:   1,
			Type: 2,
			Name: "public cloud #1",
		},
		{
			ID:   5,
			Type: 1,
			Name: "private cloud #5",
		},
	}

	clouds, _, err := client.Clouds.List(ctx, "12345")

	assert.NoError(t, err)
	assert.Equal(t, expected, clouds)
}

func TestClouds_List_emptyTenantID(t *testing.T) {
	_, _, err := client.Clouds.List(ctx, "")

	assert.Error(t, err)
	assert.Equal(t, ErrEmptyArgument, err)
}
