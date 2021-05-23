package docker

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetShortImageName(t *testing.T) {

	cases := []struct {
		name     string
		image    string
		expected string
		error    string
	}{
		{
			name:  "empty",
			image: "",
			error: `malformed image name ""`,
		},
		{
			name:     "plain",
			image:    "image",
			expected: "image",
		},
		{
			name:     "with string tag",
			image:    "image:tag",
			expected: "image",
		},
		{
			name:     "with version tag",
			image:    "image:123.456.789-beta",
			expected: "image",
		},
		{
			name:     "with prefix",
			image:    "gcr.io/123456/image",
			expected: "image",
		},
		{
			name:     "with prefix and version tag",
			image:    "gcr.io/123456/image:123.456.789-beta",
			expected: "image",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := require.New(t)

			actual, err := GetShortImageName(c.image)

			if c.error != "" {
				r.NotNil(err)
				r.Equal(c.error, err.Error())
				return
			}

			r.NoError(err)
			r.Equal(c.expected, actual)
		})
	}
}
