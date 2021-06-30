package yey

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizePathName(t *testing.T) {
	testCases := []struct {
		Name     string
		Value    string
		Expected string
	}{
		{
			Name:     "removes space from  name",
			Value:    "my folder with   spaces",
			Expected: "my_folder_with_spaces",
		},
		{
			Name:     "strips invalid characters",
			Value:    "my&folder#cool*stuff!",
			Expected: "myfoldercoolstuff",
		},
		{
			Name:     "padds with 0 if does not start with alphanumeric character",
			Value:    "_valid_if_not_for_underscore_start",
			Expected: "0_valid_if_not_for_underscore_start",
		},
		{
			Name:     "reduces multidashes to single dash",
			Value:    "my---folder----path",
			Expected: "my-folder-path",
		},
		{
			Name:     "strips trailing dashes from name",
			Value:    "my_folder----",
			Expected: "my_folder",
		},
		{
			Name:     "empty string",
			Value:    "",
			Expected: "",
		},
		{
			Name:     "full strip returns empty string",
			Value:    "$%^&",
			Expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			actual := sanitizePathName(tc.Value)
			if actual != tc.Expected {
				t.Fatalf("expected %s but got: %s", tc.Expected, actual)
			}
		})
	}
}

func TestContainerPathPrefix(t *testing.T) {
	testCases := []struct {
		Name     string
		Value    string
		Expected string
	}{
		{
			Name:     "simple base name with hash prefix",
			Value:    "/root/projectName/.yeyrc.yaml",
			Expected: "yey-projectName-45c6afaff136ad78",
		},
		{
			Name:     "root path",
			Value:    "/.yeyrc.yaml",
			Expected: "yey-8767bf26c174b05d",
		},
		{
			Name:     "path base that sanitizes to empty string",
			Value:    "/specialBase/#%!/.yeyrc.yaml",
			Expected: "yey-19bd96cb89cf189c",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			actual := ContainerPathPrefix(tc.Value)
			if actual != tc.Expected {
				t.Fatalf("expected %s but got: %s", tc.Expected, actual)
			}
		})
	}
}

func TestContainerNamePattern(t *testing.T) {
	names := [][]string{
		{"gcp", "aws"},
		{"devops"},
		{"stg", "prod"},
	}
	pattern := ContainerNamePattern(names)

	assert.True(t, pattern.MatchString("yey-path-123456-gcp-devops-stg-whatever-123456"))
	assert.True(t, pattern.MatchString("yey-path-123456-gcp-devops-prod-whatever-123456"))
	assert.True(t, pattern.MatchString("yey-path-123456-aws-devops-stg-whatever-123456"))
	assert.True(t, pattern.MatchString("yey-path-123456-aws-devops-prod-whatever-123456"))

	assert.False(t, pattern.MatchString("yey-path-123456-wrong-devops-stg-whatever-123456"))
	assert.False(t, pattern.MatchString("yey-path-123456-gcp-wrong-stg-whatever-123456"))
	assert.False(t, pattern.MatchString("yey-path-123456-gcp-devops-wrong-whatever-123456"))
	assert.False(t, pattern.MatchString("yey-path-123456-wrong-wrong-wrong-whatever-123456"))
}
