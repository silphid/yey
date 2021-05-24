package yey

import (
	"fmt"
	"sort"
)

const (
	BaseContextName = "base"
)

// Contexts represents a combinaison of base and named contexts
type Contexts struct {
	Context `yaml:",inline"`
	Named   map[string]Context
}

// Merge creates a deep-copy of this object and copies values from given source object on top of it
func (c Contexts) Merge(source Contexts) Contexts {
	merged := Contexts{
		Context: c.Context.Merge(source.Context),
		Named:   make(map[string]Context),
	}
	for key, value := range c.Named {
		merged.Named[key] = value.Clone()
	}
	for key, value := range source.Named {
		existing, ok := merged.Named[key]
		if ok {
			merged.Named[key] = existing.Merge(value)
		} else {
			merged.Named[key] = value
		}
	}
	return merged
}

// GetNames returns the list of all context names user can choose from,
// including the special "base" contexts.
func (c Contexts) GetNames() []string {
	// Extract unique names
	namesMap := make(map[string]bool)
	for name := range c.Named {
		namesMap[name] = true
	}

	// Sort
	sortedNames := make([]string, 0, len(namesMap))
	for name := range namesMap {
		sortedNames = append(sortedNames, name)
	}
	sort.Strings(sortedNames)

	// Prepend special contexts
	names := make([]string, 0, len(sortedNames)+1)
	names = append(names, BaseContextName)
	names = append(names, sortedNames...)

	return names
}

// GetContext returns context with given name, or base context
// if name is "base".
func (c Contexts) GetContext(name string) (Context, error) {
	base := c.Context
	if name == BaseContextName {
		base.Name = "base"
		return base, nil
	}
	named, ok := c.Named[name]
	if !ok {
		return Context{}, fmt.Errorf("context %q not found", name)
	}
	merged := base.Merge(named)
	merged.Name = name
	return merged, nil
}
