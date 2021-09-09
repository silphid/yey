package yey

// Contexts represents a combinaison of base and named contexts
type Contexts struct {
	Path string
	Context
}

// Merge creates a deep-copy of this object and copies values from given source object on top of it
func (c Contexts) Merge(source Contexts) Contexts {
	return Contexts{
		Context: c.Context.Merge(source.Context, true),
	}
}
