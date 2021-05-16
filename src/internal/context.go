package yey

import (
	"fmt"
)

type Context struct {
	// Do we need a name? Is it not set in the context map as the key?
	// Name      string `yaml:"-"`
	Image     string            `yaml:"image"`
	Container string            `yaml:"container"`
	Env       map[string]string `yaml:"env"`
	Mounts    map[string]string `yaml:"mounts"`
}

// Clone returns a deep-copy of this context
func (c Context) Clone() Context {
	env := c.Env
	c.Env = make(map[string]string, len(env))
	for key, value := range env {
		c.Env[key] = value
	}

	mounts := c.Mounts
	c.Mounts = make(map[string]string, len(mounts))
	for key, value := range mounts {
		c.Mounts[key] = value
	}

	return c
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

type Contexts struct {
	base  Context
	named map[string]Context
}

func (ctx Contexts) Get(name string) (*Context, error) {
	if name == "" {
		return &ctx.base, nil
	}

	selected, ok := ctx.named[name]
	if !ok {
		return nil, fmt.Errorf("no context named %q detected", name)
	}

	result := ctx.base.Merge(selected)

	return &result, nil
}
