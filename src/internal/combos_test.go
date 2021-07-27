package yey

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCombos(t *testing.T) {
	ctx := Context{
		Variations: Variations{
			Variation{
				Name: "variation1",
				Contexts: map[string]Context{
					"dev": {
						Name: "dev",
						Variations: Variations{
							Variation{
								Name: "childVariation",
								Contexts: map[string]Context{
									"dev1": {Name: "dev1"},
									"dev2": {Name: "dev2"},
								},
							},
						},
					},
					"stg":  {Name: "stg"},
					"prod": {Name: "prod"},
				},
			},
			Variation{
				Name: "variation2",
				Contexts: map[string]Context{
					"go":   {Name: "go"},
					"node": {Name: "node"},
				},
			},
		},
	}

	expected := [][]string{
		{"dev", "dev1", "go"},
		{"dev", "dev1", "node"},
		{"dev", "dev2", "go"},
		{"dev", "dev2", "node"},
		{"prod", "go"},
		{"prod", "node"},
		{"stg", "go"},
		{"stg", "node"},
	}

	actual := ctx.GetCombos()

	assert.Equal(t, expected, actual)
}
