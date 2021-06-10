package yey

type Layer struct {
	Name     string
	Contexts map[string]Context
}

// Clone returns a deep-copy of this layer
func (l Layer) Clone() Layer {
	clone := l
	clone.Contexts = make(map[string]Context, len(l.Contexts))
	for key, value := range l.Contexts {
		clone.Contexts[key] = value.Clone()
	}
	return clone
}

// Merge creates a deep-copy of this layer and copies values from given source layer on top of it
func (l Layer) Merge(source Layer) Layer {
	merged := Layer{
		Name:     l.Name,
		Contexts: make(map[string]Context),
	}
	for key, value := range l.Contexts {
		merged.Contexts[key] = value.Clone()
	}
	for key, value := range source.Contexts {
		existing, ok := merged.Contexts[key]
		if ok {
			merged.Contexts[key] = existing.Merge(value)
		} else {
			merged.Contexts[key] = value
		}
	}
	return merged
}
