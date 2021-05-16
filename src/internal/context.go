package yey

import (
	"fmt"
)

type Context struct {
	Name      string            `yaml:"-"`
	Image     string            `yaml:"image"`
	Container string            `yaml:"container"`
	Env       map[string]string `yaml:"env"`
	Mounts    map[string]string `yaml:"mounts"`
}

// Clone returns a deep-copy of this context
func (c Context) Clone() Context {
	clone := c

	clone.Env = make(map[string]string, len(c.Env))
	for key, value := range c.Env {
		clone.Env[key] = value
	}

	clone.Mounts = make(map[string]string, len(c.Mounts))
	for key, value := range c.Mounts {
		clone.Mounts[key] = value
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
	result.Name = name

	return &result, nil
}

func (ctx Contexts) Merge(other Contexts) Contexts {
	ctx.base = ctx.base.Merge(other.base)

	nameMap := map[string]struct{}{}
	for key := range ctx.named {
		nameMap[key] = struct{}{}
	}
	for key := range other.named {
		nameMap[key] = struct{}{}
	}

	named := ctx.named
	ctx.named = make(map[string]Context, len(nameMap))
	for key := range nameMap {
		ctx.named[key] = named[key].Merge(other.named[key])
	}

	return ctx
}
