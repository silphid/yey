package run

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTagFromImageName(t *testing.T) {
	assert.Equal(t, "", getTagFromImageName(""))
	assert.Equal(t, "", getTagFromImageName("gcr.io/project-1a2b3c4d5e/abcdef/image"))
	assert.Equal(t, "latest", getTagFromImageName("gcr.io/project-1a2b3c4d5e/abcdef/image:latest"))
	assert.Equal(t, "0.123.4", getTagFromImageName("gcr.io/project-1a2b3c4d5e/abcdef/image:0.123.4"))
}

func TestShouldPull(t *testing.T) {
	assert.Equal(t, true, shouldPull(""))
	assert.Equal(t, true, shouldPull("gcr.io/project-1a2b3c4d5e/abcdef/image"))
	assert.Equal(t, true, shouldPull("gcr.io/project-1a2b3c4d5e/abcdef/image:latest"))
	assert.Equal(t, false, shouldPull("gcr.io/project-1a2b3c4d5e/abcdef/image:0.123.4"))
}
