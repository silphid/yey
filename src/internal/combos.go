package yey

import (
	"sort"
)

// GetCombos returns the list of all possible context name combinations user can choose from
func (c Context) GetCombos() [][]string {
	// Any child variations?
	if len(c.Variations) > 0 {
		variationCombos := c.Variations.getCombos()
		if c.Name == "" {
			return variationCombos
		}

		// Prepend context name to all variation combos
		baseCombo := []string{c.Name}
		combos := make([][]string, 0)
		for _, variationCombo := range variationCombos {
			combo := append(baseCombo, variationCombo...)
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

func (variations Variations) getCombos() [][]string {
	var combos [][]string
	for _, variation := range variations {
		variationCombos := variation.getCombos()
		if combos == nil {
			combos = variationCombos
		} else {
			// Compute all combinations of current and variation combos
			var newCombos [][]string
			for _, combo := range combos {
				for _, variationCombo := range variationCombos {
					newCombos = append(newCombos, append(combo, variationCombo...))
				}
			}
			combos = newCombos
		}
	}
	return combos
}

func (variation Variation) getCombos() [][]string {
	// Extract variation's sorted context names
	names := make([]string, 0)
	for name := range variation.Contexts {
		names = append(names, name)
	}
	sort.Strings(names)

	// Recursively get all combos from child contexts
	var combos [][]string
	for _, name := range names {
		context := variation.Contexts[name]
		combos = append(combos, context.GetCombos()...)
	}
	return combos
}
