package xelon

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTenant_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/tenant/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = fmt.Fprint(w, `{"tenant_identifier":"long-id"}`)
	})
	expected := &Tenant{
		TenantID: "long-id",
	}

	tenant, _, err := client.Tenants.GetCurrent(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expected, tenant)
}
