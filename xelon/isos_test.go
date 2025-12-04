package xelon

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestISOs_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/isos", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = fmt.Fprint(w, `
{
  "data": [
    {"identifier":"abc123","name":"iso-1","status":true},
    {"identifier":"def456","name":"iso-2","status":false}
  ],
  "meta": {"total":20,"lastPage":2,"perPage":10,"currentPage":1,"from":1,"to":10}
}
`)
	})
	expectedISOs := []ISO{
		{ID: "abc123", Name: "iso-1", Status: true},
		{ID: "def456", Name: "iso-2", Status: false},
	}
	expectedMeta := &Meta{Total: 20, LastPage: 2, PerPage: 10, Page: 1, From: 1, To: 10}

	actualISOs, resp, err := client.ISOs.List(ctx, nil)

	assert.NoError(t, err)
	assert.Equal(t, expectedISOs, actualISOs)
	assert.Equal(t, expectedMeta, resp.Meta)
}
