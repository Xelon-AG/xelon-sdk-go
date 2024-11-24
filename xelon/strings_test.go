package xelon

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringify(t *testing.T) {
	var nilPointer *string

	cases := []struct {
		description string
		input       interface{}
		expected    string
	}{
		{"BasicTypeString", "foo", `"foo"`},
		{"BasicTypeNumber", 123, `123`},
		{"BasicTypeFloatNumber", 1.8, `1.8`},
		{"BasicTypeBoolean", false, `false`},
		{"BasicTypeArray", []string{"a", "b"}, `["a" "b"]`},
		{"BasicTypeStruct", struct{ A []string }{nil}, `{}`},
		{"BasicTypeStructNoNameType", struct{ A string }{"foo"}, `{A:"foo"}`},

		{"PointerTypeNil", nilPointer, `<nil>`},
		{"PointerTypeString", stringPointer("foo"), `"foo"`},
		{"PointerTypeInt", intPointer(123), `123`},
		{"PointerTypeBool", boolPointer(false), `false`},
		{"PointerTypeArray", []*string{stringPointer("a"), stringPointer("b")}, `["a" "b"]`},

		{
			"XelonTypeSSHKey",
			SSHKeyV1{Fingerprint: "fingerprint-text", ID: 1, Name: "name-text", PublicKey: "public-key-text"},
			`xelon.SSHKeyV1{CreatedAt:"", Fingerprint:"fingerprint-text", ID:1, Name:"name-text", PublicKey:"public-key-text"}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			s := Stringify(tc.input)
			assert.Equal(t, tc.expected, s)
		})
	}
}

func stringPointer(s string) *string {
	return &s
}

func intPointer(i int) *int {
	return &i
}

func boolPointer(b bool) *bool {
	return &b
}
