package tidy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForEachPossibleNameCombination(t *testing.T) {
	allNames := [][]string{
		{"dev", "stg", "prod"},
		{"go", "node"},
	}
	expected := [][]string{
		{"dev", "go"},
		{"dev", "node"},
		{"stg", "go"},
		{"stg", "node"},
		{"prod", "go"},
		{"prod", "node"},
	}
	var actual [][]string
	forEachPossibleNameCombination(allNames, nil, func(combo []string) error {
		actual = append(actual, combo)
		return nil
	})

	assert.Equal(t, expected, actual)
}
