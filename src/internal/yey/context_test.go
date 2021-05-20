package yey

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/go-test/deep"
	"github.com/silphid/yey/cli/src/internal/helpers"
	_assert "github.com/stretchr/testify/assert"
)

func loadContext(file string) Context {
	path := filepath.Join("testdata", file+".yaml")
	if !helpers.PathExists(path) {
		panic(fmt.Errorf("context file not found: %q", path))
	}
	contexts, err := readAndParseContextFile(path)
	if err != nil {
		panic(fmt.Errorf("loading context from %q: %w", path, err))
	}
	return contexts.Context
}

func assertNotSameMapStringString(t *testing.T, map1, map2 map[string]string, msgAndArgs ...interface{}) {
	_assert.NotEqual(t, reflect.ValueOf(map1).Pointer(), reflect.ValueOf(map2).Pointer(), msgAndArgs)
}

func TestClone(t *testing.T) {
	assert := _assert.New(t)

	original := Context{
		Image: "image",
		Env: map[string]string{
			"ENV1": "value1",
			"ENV2": "value2",
		},
		Mounts: map[string]string{
			"/local/mount1": "/container/mount1",
			"/local/mount2": "/container/mount2",
		},
	}

	clone := original.Clone()

	diff := deep.Equal(original, clone)
	if diff != nil {
		assert.Fail("Cloned context different from original", diff)
	}

	assertNotSameMapStringString(t, original.Env, clone.Env)
}

func TestMerge(t *testing.T) {
	cases := []struct {
		parent string
		child  string
	}{
		{
			parent: "base1",
			child:  "base1b",
		},
		{
			parent: "ctx1",
			child:  "ctx1b",
		},
	}

	for _, c := range cases {
		expectedName := fmt.Sprintf("%s_%s", c.parent, c.child)
		t.Run(expectedName, func(t *testing.T) {

			parent := loadContext(c.parent)
			child := loadContext(c.child)
			actual := parent.Merge(child)
			expected := loadContext(expectedName)

			if diff := deep.Equal(expected, actual); diff != nil {
				t.Error(diff)
			}

			assertNotSameMapStringString(t, actual.Env, parent.Env)
			assertNotSameMapStringString(t, actual.Env, child.Env)
			assertNotSameMapStringString(t, actual.Mounts, parent.Mounts)
			assertNotSameMapStringString(t, actual.Mounts, child.Mounts)
		})
	}
}
