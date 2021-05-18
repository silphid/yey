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
	Context  `yaml:",inline"`
	contexts map[string]Context
}

// Clone returns a deep-copy of this object
func (c Contexts) Clone() Contexts {
	clone := c
	for key, value := range c.contexts {
		clone.contexts[key] = value.Clone()
	}
	return clone
}

// Merge creates a deep-copy of this object and copies values from given source object on top of it
func (c Contexts) Merge(source Contexts) Contexts {
	clone := c.Clone()
	clone.Context.Merge(source.Context)
	for key, value := range source.contexts {
		existing, ok := clone.contexts[key]
		if ok {
			c.contexts[key] = existing.Merge(value)
		} else {
			c.contexts[key] = value
		}
	}
	return clone
}

// GetContextNames returns the list of all context names user can choose from,
// including the special "base" contexts.
func (c Contexts) GetContextNames() ([]string, error) {
	// Extract unique names
	namesMap := make(map[string]bool)
	for name := range c.contexts {
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
	base := c.Context
	if name == BaseContextName {
		return base, nil
	}
	context, ok := c.contexts[name]
	if !ok {
		return Context{}, fmt.Errorf("context %q not found", name)
	}
	return base.Merge(context), nil
}
