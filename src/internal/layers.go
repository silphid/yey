package yey

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Layers []Layer

// Clone returns a deep-copy of this layer
func (l Layers) Clone() Layers {
	clone := Layers{}
	for _, layer := range l {
		clone = append(clone, layer.Clone())
	}
	return clone
}

// GetByName returns layer with given name and whether it was found
func (l Layers) GetByName(name string) (Layer, bool) {
	for _, layer := range l {
		if layer.Name == name {
			return layer, true
		}
	}
	return Layer{}, false
}

// Merge creates a deep-copy of this layer and copies values from given source layer on top of it
func (l Layers) Merge(source Layers) Layers {
	merged := l.Clone()
	for i, layer := range source {
		existing, ok := merged.GetByName(layer.Name)
		if ok {
			merged[i] = existing.Merge(layer)
		} else {
			merged = append(merged, layer.Clone())
		}
	}
	return merged
}

func (layers *Layers) UnmarshalYAML(n *yaml.Node) error {
	if n.Kind != yaml.MappingNode {
		return fmt.Errorf(`expecting map for "layers" property at line %d, column %d`, n.Line, n.Column)
	}
	for i := 0; i < len(n.Content); i += 2 {
		layerName := n.Content[i].Value
		layerMap := n.Content[i+1]
		if layerMap.Kind != yaml.MappingNode {
			return fmt.Errorf(`expecting map for layer item at line %d, column %d`, n.Line, n.Column)
		}
		var contexts map[string]Context
		if err := layerMap.Decode(&contexts); err != nil {
			return fmt.Errorf("failed to parse contexts for layer %q: %w", layerName, err)
		}
		layer := Layer{
			Name:     layerName,
			Contexts: contexts,
		}
		*layers = append(*layers, layer)
	}
	return nil
}
