package yey

type Variation struct {
	Name     string
	Contexts map[string]Context
}

// Clone returns a deep-copy of this variation
func (l Variation) Clone() Variation {
	clone := l
	clone.Contexts = make(map[string]Context, len(l.Contexts))
	for key, value := range l.Contexts {
		clone.Contexts[key] = value.Clone()
	}
	return clone
}

// Merge creates a deep-copy of this variation and copies values from given source variation on top of it
func (l Variation) Merge(source Variation) Variation {
	merged := Variation{
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
