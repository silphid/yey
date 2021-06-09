package cmd

import (
	"testing"

	_assert "github.com/stretchr/testify/assert"
)

func TestParseNameAndVariant(t *testing.T) {
	assert := _assert.New(t)

	cases := []struct {
		name            string
		value           string
		expectedName    string
		expectedVariant string
		error           string
	}{
		{
			name:            "empty",
			value:           "",
			expectedName:    "",
			expectedVariant: "",
		},
		{
			name:            "just the name",
			value:           "some_name",
			expectedName:    "some_name",
			expectedVariant: "",
		},
		{
			name:            "name and variant",
			value:           "some_name/some_variant",
			expectedName:    "some_name",
			expectedVariant: "some_variant",
		},
		{
			name:  "just the slash",
			value: "/",
			error: `malformed context name: "/"`,
		},
		{
			name:  "just the variant",
			value: "/some_variant",
			error: `malformed context name: "/some_variant"`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			actualName, actualVariant, err := parseNameAndVariant(tc.value)

			if tc.error != "" {
				assert.NotNil(err)
				assert.Equal(tc.error, err.Error())
				return
			}

			assert.Equal(tc.expectedName, actualName, "name")
			assert.Equal(tc.expectedVariant, actualVariant, "variant")
		})
	}
}
