package xelon

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFirewalls_All(t *testing.T) {
	setup(t)
	defer teardown()

	requestCount := 0
	mux.HandleFunc("GET /firewalls", func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, []string{"tenant-1", "tenant-2"}, r.URL.Query()["tenantIds[]"])
		assert.Equal(t, "edge", r.URL.Query().Get("search"))
		assert.Equal(t, "name", r.URL.Query().Get("sort"))
		assert.Equal(t, "1", r.URL.Query().Get("perPage"))

		switch r.URL.Query().Get("page") {
		case "1":
			_, _ = w.Write([]byte(`{"data":[{"identifier":"firewall-1"}],"meta":{"lastPage":2,"currentPage":1}}`))
		case "2":
			_, _ = w.Write([]byte(`{"data":[{"identifier":"firewall-2"}],"meta":{"lastPage":2,"currentPage":2}}`))
		default:
			t.Fatalf("unexpected page %q", r.URL.Query().Get("page"))
		}
	})

	opts := &FirewallListOptions{
		Search:    "edge",
		Sort:      "name",
		TenantIDs: []string{"tenant-1", "tenant-2"},
		ListOptions: ListOptions{
			Page:    1,
			PerPage: 1,
		},
	}
	seq, errFn := client.Firewalls.All(ctx, opts)

	var actualFirewalls []Firewall
	for firewall := range seq {
		actualFirewalls = append(actualFirewalls, firewall)
	}

	assert.NoError(t, errFn())
	assert.Equal(t, 2, requestCount)
	assert.Equal(t, []Firewall{{ID: "firewall-1"}, {ID: "firewall-2"}}, actualFirewalls)
	assert.Equal(t, 1, opts.Page)
	assert.Equal(t, 1, opts.PerPage)
	assert.Equal(t, []string{"tenant-1", "tenant-2"}, opts.TenantIDs)
}

func TestFirewallForwardingRule_UnmarshalJSON(t *testing.T) {
	type testCase struct {
		input  string
		expect FirewallForwardingRule
	}
	tests := map[string]testCase{
		"single sourceIp_multiple destinationIp": {
			input: `
{
  "identifier": "123abc456def",
  "port": 10000,
  "externalPort": 10000,
  "protocol": "tcp",
  "sourceIp": "10.0.0.80",
  "destinationIp": ["0.0.0.0\/0"],
  "type": "outbound"
}
`,
			expect: FirewallForwardingRule{
				ID:                          "123abc456def",
				InternalPort:                10000,
				ExternalPort:                10000,
				Protocol:                    "tcp",
				Type:                        "outbound",
				DestinationIPAddressWrapper: []any{"0.0.0.0/0"},
				DestinationIPAddress:        "",
				DestinationIPAddresses:      []string{"0.0.0.0/0"},
				SourceIPAddressWrapper:      "10.0.0.80",
				SourceIPAddress:             "10.0.0.80",
				SourceIPAddresses:           nil,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var actual FirewallForwardingRule
			err := json.Unmarshal([]byte(test.input), &actual)

			assert.Nil(t, err)
			assert.Equal(t, test.expect, actual)
		})
	}
}

func TestFirewallForwardingRule_MarshalJSON(t *testing.T) {
	type testCase struct {
		input  *FirewallForwardingRule
		expect string
	}
	tests := map[string]testCase{
		"single sourceIp_multiple destinationIps": {
			input: &FirewallForwardingRule{
				ID:                          "123abc456def",
				InternalPort:                10000,
				ExternalPort:                10000,
				Protocol:                    "tcp",
				Type:                        "outbound",
				DestinationIPAddressWrapper: []any{"0.0.0.0/0"},
				DestinationIPAddress:        "",
				DestinationIPAddresses:      []string{"0.0.0.0/0"},
				SourceIPAddressWrapper:      "10.0.0.80",
				SourceIPAddress:             "10.0.0.80",
				SourceIPAddresses:           nil,
			},
			expect: "{\"externalPort\":10000,\"identifier\":\"123abc456def\",\"port\":10000,\"protocol\":\"tcp\",\"type\":\"outbound\",\"destinationIp\":[\"0.0.0.0/0\"],\"sourceIp\":\"10.0.0.80\"}",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			bytes, err := json.Marshal(test.input)

			actual := strings.TrimSpace(string(bytes))

			assert.Nil(t, err)
			assert.Equal(t, test.expect, actual)
		})
	}
}
