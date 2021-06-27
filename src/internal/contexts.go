package yey

import (
	"fmt"
	"sort"
)

// Contexts represents a combinaison of base and named contexts
type Contexts struct {
	Path string
	Context
	Layers Layers
}

// Merge creates a deep-copy of this object and copies values from given source object on top of it
func (c Contexts) Merge(source Contexts) Contexts {
	return Contexts{
		Context: c.Context.Merge(source.Context),
		Layers:  c.Layers.Merge(source.Layers),
	}
}

// GetNamesInAllLayers returns the list of all context names user can choose from,
func (c Contexts) GetNamesInAllLayers() [][]string {
	names := make([][]string, 0, len(c.Layers))

	for _, layer := range c.Layers {
		// Extract unique names
		namesMap := make(map[string]bool)
		for name := range layer.Contexts {
			namesMap[name] = true
		}

		// Sort
		sortedNames := make([]string, 0, len(namesMap))
		for name := range namesMap {
			sortedNames = append(sortedNames, name)
		}
		sort.Strings(sortedNames)

		names = append(names, sortedNames)
	}

	return names
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

	// Sort
	sortedNames := make([]string, 0, len(namesMap))
	for name := range namesMap {
		sortedNames = append(sortedNames, name)
	}
	sort.Strings(sortedNames)
	return sortedNames
}

// GetContext returns context with given name (or base context, if name is "base") and
// variant (or no variant, if variant name is "")
func (c Contexts) GetContext(names []string) (Context, error) {
	if len(names) != len(c.Layers) {
		return Context{}, fmt.Errorf("number of context names (%d) does not match number of layers (%d)", len(names), len(c.Layers))
	}

	// Start with base context
	ctx := c.Context
	compositeName := ""

	for i, layer := range c.Layers {
		name := names[i]

		// Merge layer context
		layerContext, ok := layer.Contexts[name]
		if !ok {
			return Context{}, fmt.Errorf("context %q not found in layer %q", name, layer.Name)
		}
		ctx = ctx.Merge(layerContext)

		// Accumulate composite name
		if compositeName == "" {
			compositeName = name
		} else {
			compositeName = fmt.Sprintf("%s %s", compositeName, name)
		}
	}

	ctx.Name = compositeName
	return ctx, nil
}
