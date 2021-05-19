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
	base          Context
	namedContexts map[string]Context
}

// Merge creates a deep-copy of this object and copies values from given source object on top of it
func (c Contexts) Merge(source Contexts) Contexts {
	merged := Contexts{
		base:          c.base.Merge(source.base),
		namedContexts: make(map[string]Context),
	}

	for key, value := range c.namedContexts {
		merged.namedContexts[key] = value.Clone()
	}

	for key, value := range source.namedContexts {
		existing, ok := merged.namedContexts[key]
		if ok {
			merged.namedContexts[key] = existing.Merge(value)
		} else {
			merged.namedContexts[key] = value
		}
	}

	return merged
}

// GetContextNames returns the list of all context names user can choose from,
// including the special "base" contexts.
func (c Contexts) GetContextNames() ([]string, error) {
	// Extract unique names
	namesMap := make(map[string]bool)
	for name := range c.namedContexts {
		namesMap[name] = true
	}

	// Sort
	sortedNames := make([]string, 0, len(namesMap))
	for name := range namesMap {
		sortedNames = append(sortedNames, name)
	}
	sort.Strings(sortedNames)

	// Prepend special contexts
	names := make([]string, 0, len(sortedNames)+2)
	names = append(names, BaseContextName)
	names = append(names, sortedNames...)

	return names, nil
}

// GetContext returns context with given name, or base context
// if name is "base".
func (c Contexts) GetContext(name string) (Context, error) {
	base := c.base
	if name == BaseContextName {
		return base, nil
	}
	context, ok := c.namedContexts[name]
	if !ok {
		return Context{}, fmt.Errorf("context %q not found", name)
	}
	return base.Merge(context), nil
}
