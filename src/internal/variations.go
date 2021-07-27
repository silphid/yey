package yey

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Variations []Variation

// Clone returns a deep-copy of this variation
func (l Variations) Clone() Variations {
	var clone Variations
	for _, variation := range l {
		clone = append(clone, variation.Clone())
	}
	return clone
}

// GetByName returns variation with given name and whether it was found
func (l Variations) GetByName(name string) (Variation, bool) {
	for _, variation := range l {
		if variation.Name == name {
			return variation, true
		}
	}
	return Variation{}, false
}

// Merge creates a deep-copy of this variation and copies values from given source variation on top of it
func (l Variations) Merge(source Variations) Variations {
	merged := l.Clone()
	for i, variation := range source {
		existing, ok := merged.GetByName(variation.Name)
		if ok {
			merged[i] = existing.Merge(variation)
		} else {
			merged = append(merged, variation.Clone())
		}
	}
	return merged
}

func (variations *Variations) UnmarshalYAML(n *yaml.Node) error {
	if n.Kind != yaml.MappingNode {
		return fmt.Errorf(`expecting map for "variations" property at line %d, column %d`, n.Line, n.Column)
	}
	for i := 0; i < len(n.Content); i += 2 {
		variationName := n.Content[i].Value
		variationMap := n.Content[i+1]
		if variationMap.Kind != yaml.MappingNode {
			return fmt.Errorf(`expecting map for variation item at line %d, column %d`, n.Line, n.Column)
		}
		var contexts map[string]Context
		if err := variationMap.Decode(&contexts); err != nil {
			return fmt.Errorf("failed to parse contexts for variation %q: %w", variationName, err)
		}
		variation := Variation{
			Name:     variationName,
			Contexts: contexts,
		}
		for name, context := range variation.Contexts {
			context.Name = name
			variation.Contexts[name] = context
		}
		*variations = append(*variations, variation)
	}
	return nil
}
