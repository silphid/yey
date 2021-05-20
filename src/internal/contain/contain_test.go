package contain

import (
	"testing"

	_assert "github.com/stretchr/testify/assert"
)

func TestGetShortImageName(t *testing.T) {
	assert := _assert.New(t)

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

			actual, err := getShortImageName(c.image)

			if c.error != "" {
				assert.NotNil(err)
				assert.Equal(c.error, err.Error())
				return
			}

			assert.NoError(err)
			assert.Equal(c.expected, actual)
		})
	}
}
