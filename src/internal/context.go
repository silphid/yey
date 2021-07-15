package yey

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type DockerBuild struct {
	Dockerfile string
	Args       map[string]string
	Context    string
}

// Context represents execution configuration for some docker container
type Context struct {
	Name       string `yaml:",omitempty"`
	Layers     Layers `yaml:"layers"`
	Remove     *bool
	Image      string
	Build      DockerBuild
	Env        map[string]string
	Mounts     map[string]string
	EntryPoint string `yaml:"entrypoint,omitempty"`
	Cmd        []string
	Network    string
}

// Clone returns a deep-copy of this context
func (c Context) Clone() Context {
	clone := c
	clone.Layers = c.Layers.Clone()
	if clone.Remove != nil {
		value := *clone.Remove
		clone.Remove = &value
	}
	clone.Env = make(map[string]string)
	for key, value := range c.Env {
		clone.Env[key] = value
	}
	clone.Mounts = make(map[string]string)
	for key, value := range c.Mounts {
		clone.Mounts[key] = value
	}
	clone.Build.Args = make(map[string]string)
	for key, value := range c.Build.Args {
		clone.Build.Args[key] = value
	}
	return clone
}

// Merge creates a deep-copy of this context and copies values from given source context on top of it
func (c Context) Merge(source Context) Context {
	merged := c.Clone()
	if source.Name != "" {
		merged.Name = source.Name
	}
	if source.Layers != nil {
		merged.Layers = merged.Layers.Merge(source.Layers)
	}
	if source.Remove != nil {
		value := *source.Remove
		merged.Remove = &value
	}
	if source.Image != "" {
		merged.Image = source.Image
	}
	for key, value := range source.Env {
		merged.Env[key] = value
	}
	for key, value := range source.Mounts {
		merged.Mounts[key] = value
	}
	if source.Build.Dockerfile != "" {
		merged.Build.Dockerfile = source.Build.Dockerfile
	}
	if source.Build.Context != "" {
		merged.Build.Context = source.Build.Context
	}
	for key, value := range source.Build.Args {
		merged.Build.Args[key] = value
	}
	if source.Network != "" {
		merged.Network = source.Network
	}
	return merged
}

// GetContext returns context resulting from merging contexts with given names from all layers
func (c Context) GetContext(names []string) (Context, error) {
	ctx, remainingNames, err := c.getContextRecursively(names)
	if err != nil {
		return Context{}, err
	}
	if len(remainingNames) > 0 {
		return Context{}, fmt.Errorf("extraneous context names: %s", strings.Join(remainingNames, " "))
	}
	return ctx, nil
}

// getContext returns context resulting from merging contexts with given names from all layers,
// as well as list of remaining names that were not used to resolve layers (used for recursivity)
func (c Context) getContextRecursively(names []string) (Context, []string, error) {
	// Start with context itself
	ctx := c

	for _, layer := range c.Layers {
		if len(names) == 0 {
			return Context{}, nil, fmt.Errorf("too few context names")
		}
		name := names[0]
		names = names[1:]

		// Merge layer context
		layerContext, ok := layer.Contexts[name]
		if !ok {
			return Context{}, nil, fmt.Errorf("context %q not found in layer %q", name, layer.Name)
		}
		ctx = ctx.Merge(layerContext)
		if len(ctx.Name) > 0 {
			ctx.Name = fmt.Sprintf("%s %s", ctx.Name, name)
		} else {
			ctx.Name = name
		}

		// Get child contexts recursively
		if len(layerContext.Layers) > 0 {
			childContext, remainingNames, err := layerContext.getContextRecursively(names)
			if err != nil {
				return Context{}, nil, err
			}
			ctx = ctx.Merge(childContext)
			ctx.Name = fmt.Sprintf("%s %s", ctx.Name, childContext.Name)
			names = remainingNames
		}
	}

	return ctx, names, nil
}

// String returns a user-friendly yaml representation of this context
func (c Context) String() string {
	buf, err := yaml.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(buf)
}
