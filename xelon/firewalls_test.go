package xelon

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
