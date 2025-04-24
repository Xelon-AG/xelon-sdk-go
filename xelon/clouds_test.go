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
{"identifier":"12345","name":"public cloud #1","type":"public","hvType":"12345"},
{"identifier":"67890","name":"private cloud #2","type":"private","hvType":"67890"}
]`)
	})
	expected := []Cloud{{
		ID:     "12345",
		Name:   "public cloud #1",
		Type:   "public",
		HVType: "12345",
	}, {
		ID:     "67890",
		Name:   "private cloud #2",
		Type:   "private",
		HVType: "67890",
	}}

	clouds, _, err := client.Clouds.List(ctx, nil)

	assert.NoError(t, err)
	assert.Equal(t, expected, clouds)
}
