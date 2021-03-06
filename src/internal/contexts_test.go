package yey

import (
	"testing"

	"github.com/go-test/deep"
	_assert "github.com/stretchr/testify/assert"
)

func loadContexts(baseFile, ctx1Key, ctx1File, ctx2Key, ctx2File string) Contexts {
	ctx := Contexts{
		Context: loadContext(baseFile),
	}
	ctx.Variations = Variations{
		Variation{
			Name: "variation",
			Contexts: map[string]Context{
				ctx1Key: loadContext(ctx1File),
				ctx2Key: loadContext(ctx2File),
			},
		},
	}
	return ctx
}

func TestGetContext(t *testing.T) {
	assert := _assert.New(t)

	parent := loadContexts("base1", "ctx1", "ctx1", "ctx2", "ctx2")
	child := loadContexts("base1b", "ctx1", "ctx1b", "ctx3", "ctx3")
	merged := parent.Merge(child)

	cases := []struct {
		name     string
		expected string
		error    string
	}{
		{
			name:     "ctx1",
			expected: "base1_base1b_ctx1_ctx1b",
		},
		{
			name:     "ctx2",
			expected: "base1_base1b_ctx2",
		},
		{
			name:     "ctx3",
			expected: "base1_base1b_ctx3",
		},
		{
			name:  "unknown",
			error: `context "unknown" not found in variation "variation"`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			actual, err := merged.GetContext([]string{c.name})
			actual.Name = c.name

			if c.error != "" {
				assert.NotNil(err)
				assert.Equal(c.error, err.Error())
				return
			}

			expected := loadContext(c.expected)
			expected.Name = c.name

			assert.NoError(err)
			if diff := deep.Equal(expected, actual); diff != nil {
				t.Error(diff)
			}
		})
	}
}
