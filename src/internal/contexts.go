package yey

import (
	"sort"
)

// Contexts represents a combinaison of base and named contexts
type Contexts struct {
	Path string
	Context
}

// Merge creates a deep-copy of this object and copies values from given source object on top of it
func (c Contexts) Merge(source Contexts) Contexts {
	return Contexts{
		Context: c.Context.Merge(source.Context),
	}
}

// GetAllImages returns the list of image names referenced in all contexts
func (c Contexts) GetAllImages() []string {
	namesMap := make(map[string]struct{})

	if c.Context.Image != "" {
		namesMap[c.Context.Image] = struct{}{}
	}

	for _, layer := range c.Layers {
		for _, ctx := range layer.Contexts {
			if ctx.Image != "" {
				namesMap[ctx.Image] = struct{}{}
			}
		}
	}

	// TODO: consider nested layers

	// Sort
	sortedNames := make([]string, 0, len(namesMap))
	for name := range namesMap {
		sortedNames = append(sortedNames, name)
	}
	sort.Strings(sortedNames)
	return sortedNames
}
