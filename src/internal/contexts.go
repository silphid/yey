package yey

import (
	"fmt"
	"sort"
)

const (
	BaseContextName        = "base"
	NoneVariantDisplayName = "none"
	NoneVariantName        = ""
)

// Contexts represents a combinaison of base and named contexts
type Contexts struct {
	Path     string `yaml:"-"`
	Context  `yaml:",inline"`
	Named    map[string]Context
	Variants map[string]Context
}

// Merge creates a deep-copy of this object and copies values from given source object on top of it
func (c Contexts) Merge(source Contexts) Contexts {
	merged := Contexts{
		Context:  c.Context.Merge(source.Context),
		Named:    make(map[string]Context),
		Variants: make(map[string]Context),
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
	for key, value := range source.Variants {
		existing, ok := merged.Variants[key]
		if ok {
			merged.Variants[key] = existing.Merge(value)
		} else {
			merged.Variants[key] = value
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

	// Prepend special "base" context
	names := make([]string, 0, len(sortedNames)+1)
	names = append(names, BaseContextName)
	names = append(names, sortedNames...)

	return names
}

// GetVariants returns the list of all context variant names and display names user can choose from
func (c Contexts) GetVariants() ([]string, []string) {
	// Extract unique names
	namesMap := make(map[string]bool)
	for name := range c.Variants {
		namesMap[name] = true
	}

	// Sort
	sortedNames := make([]string, 0, len(namesMap))
	for name := range namesMap {
		sortedNames = append(sortedNames, name)
	}
	sort.Strings(sortedNames)

	// Prepend special "none" variant
	names := make([]string, 0, len(sortedNames)+1)
	names = append(names, NoneVariantName)
	names = append(names, sortedNames...)
	displayNames := make([]string, 0, len(sortedNames)+1)
	displayNames = append(displayNames, NoneVariantDisplayName)
	displayNames = append(displayNames, sortedNames...)

	return names, displayNames
}

// GetContext returns context with given name (or base context, if name is "base") and
// variant (or no variant, if variant name is "")
func (c Contexts) GetContext(name, variant string) (Context, error) {
	// Start with base context
	ctx := c.Context

	// Merge named context, if any
	if name != BaseContextName {
		named, ok := c.Named[name]
		if !ok {
			return Context{}, fmt.Errorf("named context %q not found", name)
		}
		ctx = ctx.Merge(named)
	}
	ctx.Name = name

	// Merge variant context, if any
	if variant != "" {
		named, ok := c.Variants[variant]
		if !ok {
			return Context{}, fmt.Errorf("variant context %q not found", variant)
		}
		ctx = ctx.Merge(named)
		ctx.Name = fmt.Sprintf("%s %s", ctx.Name, variant)
	}

	return ctx, nil
}
