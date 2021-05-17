package yey

import (
	"reflect"
	"testing"

	"github.com/go-test/deep"
	_assert "github.com/stretchr/testify/assert"
)

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
	assert := _assert.New(t)

	context1 := Context{
		Env: map[string]string{
			"ENV1": "value1",
			"ENV2": "value2",
		},
		Mounts: map[string]string{
			"/local/mount1": "/container/mount1",
			"/local/mount2": "/container/mount2",
		},
	}

	context2 := Context{
		Env: map[string]string{
			"ENV1": "value1b",
			"ENV3": "value3",
		},
		Mounts: map[string]string{
			"/local/mount1": "/container/mount1b",
			"/local/mount3": "/container/mount3",
		},
	}

	expected := Context{
		Env: map[string]string{
			"ENV1": "value1b",
			"ENV2": "value2",
			"ENV3": "value3",
		},
		Mounts: map[string]string{
			"/local/mount1": "/container/mount1b",
			"/local/mount2": "/container/mount2",
			"/local/mount3": "/container/mount3",
		},
	}

	merged := context1.Merge(context2)

	diff := deep.Equal(merged, expected)
	if diff != nil {
		assert.Fail("Merged context different from expected context", diff)
	}

	assertNotSameMapStringString(t, merged.Env, context1.Env)
	assertNotSameMapStringString(t, merged.Env, context2.Env)
	assertNotSameMapStringString(t, merged.Mounts, context1.Mounts)
	assertNotSameMapStringString(t, merged.Mounts, context2.Mounts)
}

func assertNotSameMapStringString(t *testing.T, map1, map2 map[string]string, msgAndArgs ...interface{}) {
	_assert.NotEqual(t, reflect.ValueOf(map1).Pointer(), reflect.ValueOf(map2).Pointer(), msgAndArgs)
}
