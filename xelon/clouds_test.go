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

	mux.HandleFunc("/clouds", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = fmt.Fprint(w, `[
{"identifier":"12345","name":"public cloud #1","type":2,"hvType":"12345"},
{"identifier":"67890","name":"public cloud #2","type":1,"hvType":"67890"}
]`)
	})
	expected := []Cloud{{
		ID:     "12345",
		Name:   "public cloud #1",
		Type:   2,
		HVType: "12345",
	}, {
		ID:     "67890",
		Name:   "public cloud #2",
		Type:   1,
		HVType: "67890",
	}}

	clouds, _, err := client.Clouds.List(ctx, nil)

	assert.NoError(t, err)
	assert.Equal(t, expected, clouds)
}
