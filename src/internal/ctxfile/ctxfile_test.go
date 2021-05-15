package ctxfile

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/go-test/deep"
	"github.com/silphid/yey/cli/src/internal/ctx"
	"github.com/silphid/yey/cli/src/internal/helpers"
	_assert "github.com/stretchr/testify/assert"
)

func loadContext(file string) ctx.Context {
	path := filepath.Join("testdata", file+".yaml")
	if !helpers.PathExists(path) {
		panic(fmt.Errorf("context file not found: %q", path))
	}
	context, err := ctx.Load(path)
	if err != nil {
		panic(fmt.Errorf("loading context from %q: %w", path, err))
	}
	return context
}

func loadConfig(version, baseFile, ctx1Key, ctx1File, ctx2Key, ctx2File string) ContextFile {
	return ContextFile{
		Version: version,
		Base:    loadContext(baseFile),
		NamedContexts: map[string]ctx.Context{
			ctx1Key: loadContext(ctx1File),
			ctx2Key: loadContext(ctx2File),
		},
	}
}

func TestGetContext(t *testing.T) {
	assert := _assert.New(t)

	sharedConfig := loadConfig("version1", "base1", "ctx1", "ctx1", "ctx2", "ctx2")
	userConfig := loadConfig("version2", "base1b", "ctx1", "ctx1b", "ctx3", "ctx3")

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
			name:     "base",
			expected: "base1_base1b",
		},
		{
			name:     "none",
			expected: "none",
		},
		{
			name:  "unknown",
			error: `context not found "unknown"`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual, err := getContext(sharedConfig, userConfig, c.name)

			if c.error != "" {
				assert.NotNil(err)
				assert.Equal(c.error, err.Error())
				return
			}

			expected := loadContext(c.expected)
			assert.NoError(err)
			if diff := deep.Equal(expected, actual); diff != nil {
				t.Error(diff)
			}
		})
	}
}
