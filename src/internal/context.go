package yey

// Context represents an execution context for yey (env vars and volumes)
type Context struct {
	Name      string `yaml:"-"`
	Image     string
	Container string
	Env       map[string]string
	Mounts    map[string]string
}

var None = Context{
	Env:    make(map[string]string),
	Mounts: make(map[string]string),
}

// Clone returns a deep-copy of this context
func (c Context) Clone() Context {
	clone := c
	clone.Env = make(map[string]string)
	for key, value := range c.Env {
		clone.Env[key] = value
	}
	clone.Mounts = make(map[string]string)
	for key, value := range c.Mounts {
		clone.Mounts[key] = value
	}
	return clone
}

// Merge creates a deep-copy of this context and copies values from given source context on top of it
func (c Context) Merge(source Context) Context {
	context := c.Clone()
	if source.Container != "" {
		context.Container = source.Container
	}
	if source.Image != "" {
		context.Image = source.Image
	}
	for key, value := range source.Env {
		context.Env[key] = value
	}
	for key, value := range source.Mounts {
		context.Mounts[key] = value
	}
	return context
}
