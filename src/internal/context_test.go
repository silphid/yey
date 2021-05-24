package yey

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
)

func loadContext(file string) Context {
	path := filepath.Join("testdata", file+".yaml")
	contexts, err := readAndParseContextFileFromURI(path)
	if err != nil {
		panic(fmt.Errorf("loading context from %q: %w", path, err))
	}
	return contexts.Context
}

func assertNotSameMapStringString(t *testing.T, map1, map2 map[string]string, msgAndArgs ...interface{}) {
	assert.NotEqual(t, reflect.ValueOf(map1).Pointer(), reflect.ValueOf(map2).Pointer(), msgAndArgs)
}

func TestClone(t *testing.T) {
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
		assert.Fail(t, "Cloned context different from original", diff)
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

func TestSameHashesForSameContexts(t *testing.T) {
	ctx1 := getCtx1()
	hash1 := ctx1.Hash()

	ctx2 := getCtx1()
	hash2 := ctx2.Hash()

	assert.Equal(t, hash1, hash2)
}

func TestDifferentHashesForDifferentEnvs(t *testing.T) {
	ctx1 := getCtx1()
	hash1 := ctx1.Hash()

	ctx2 := getCtx1()
	ctx2.Env["ENV1"] = "value1b"
	hash2 := ctx2.Hash()

	assert.NotEqual(t, hash1, hash2)
}

func TestDifferentHashesForDifferentMounts(t *testing.T) {
	ctx1 := getCtx1()
	hash1 := ctx1.Hash()

	ctx2 := getCtx1()
	ctx2.Mounts["/local/mount1"] = "/container/mount1b"
	hash2 := ctx2.Hash()

	assert.NotEqual(t, hash1, hash2)
}

func getCtx1() Context {
	return Context{
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
}
