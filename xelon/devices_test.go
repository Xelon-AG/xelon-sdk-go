package xelon

import (
	"encoding/json"
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDevices_DeviceNetworkIPAddresses_UnmarshalJSON(t *testing.T) {
	type testCase struct {
		input     string
		expect    []netip.Addr
		expectErr bool
	}
	tests := map[string]testCase{
		"single string": {
			input:     `"10.0.0.1"`,
			expect:    []netip.Addr{netip.MustParseAddr("10.0.0.1")},
			expectErr: false,
		},
		"array of strings": {
			input: `["10.0.0.1", "2001:db8::1"]`,
			expect: []netip.Addr{
				netip.MustParseAddr("10.0.0.1"),
				netip.MustParseAddr("2001:db8::1"),
			},
			expectErr: false,
		},
		"empty array": {
			input:     `[]`,
			expect:    []netip.Addr{},
			expectErr: false,
		},
		"null becomes empty": {
			input:     `null`,
			expect:    []netip.Addr{},
			expectErr: false,
		},
		"invalid ip rejected": {
			input:     `"not-an-ip"`,
			expect:    nil,
			expectErr: true,
		},
		"one bad entry rejects whole list": {
			input:     `["10.0.0.1", "not-an-ip"]`,
			expect:    nil,
			expectErr: true,
		},
		"wrong type rejected": {
			input:     `42`,
			expect:    nil,
			expectErr: true,
		},
		"unspecified rejected": {
			input:     `"0.0.0.0"`,
			expect:    nil,
			expectErr: true,
		},
		"duplicates rejected": {
			input:     `["10.0.0.1", ""10.0.0.1"]`,
			expect:    nil,
			expectErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var actual DeviceNetworkIPAddresses
			err := json.Unmarshal([]byte(test.input), &actual)

			if test.expectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.ElementsMatch(t, test.expect, actual)
		})
	}
}
