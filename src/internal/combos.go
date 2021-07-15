package yey

import (
	"sort"
)

// GetCombos returns the list of all possible context name combinations user can choose from
func (c Context) GetCombos() [][]string {
	// Any child layers?
	if len(c.Layers) > 0 {
		layerCombos := c.Layers.getCombos()
		if c.Name == "" {
			return layerCombos
		}

		// Prepend context name to all layer combos
		baseCombo := []string{c.Name}
		combos := make([][]string, 0)
		for _, layerCombo := range layerCombos {
			combo := append(baseCombo, layerCombo...)
			combos = append(combos, combo)
		}

		return combos
	}

	// Single context name
	if c.Name != "" {
		return [][]string{{c.Name}}
	}

	// Nothing
	return [][]string{}
}

func (layers Layers) getCombos() [][]string {
	var combos [][]string
	for _, layer := range layers {
		layerCombos := layer.getCombos()
		if combos == nil {
			combos = layerCombos
		} else {
			// Compute all combinations of current and layer combos
			var newCombos [][]string
			for _, combo := range combos {
				for _, layerCombo := range layerCombos {
					newCombos = append(newCombos, append(combo, layerCombo...))
				}
			}
			combos = newCombos
		}
	}
	return combos
}

func (layer Layer) getCombos() [][]string {
	// Extract layer's sorted context names
	names := make([]string, 0)
	for name := range layer.Contexts {
		names = append(names, name)
	}
	sort.Strings(names)

	// Recursively get all combos from child contexts
	var combos [][]string
	for _, name := range names {
		context := layer.Contexts[name]
		combos = append(combos, context.GetCombos()...)
	}
	return combos
}
