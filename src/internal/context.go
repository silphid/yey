package yey

import (
	"fmt"
	"sort"
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
	Name       string     `yaml:",omitempty"`
	Variations Variations `yaml:"variations"`
	Remove     *bool
	Image      string
	Build      DockerBuild
	Env        map[string]string
	Mounts     map[string]string
	EntryPoint string `yaml:"entrypoint,omitempty"`
	Cmd        []string
	Network    string
	DockerArgs []string `yaml:"dockerArgs,omitempty"`
}

// Clone returns a deep-copy of this context
func (c Context) Clone() Context {
	clone := c
	clone.Variations = c.Variations.Clone()
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
func (c Context) Merge(source Context, withVariations bool) Context {
	merged := c.Clone()
	if source.Name != "" {
		merged.Name = source.Name
	}
	if withVariations {
		if source.Variations != nil {
			merged.Variations = merged.Variations.Merge(source.Variations)
		}
	} else {
		merged.Variations = nil
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
	merged.Cmd = append(merged.Cmd, source.Cmd...)
	merged.DockerArgs = append(merged.DockerArgs, source.DockerArgs...)
	return merged
}

// GetContext returns context resulting from merging contexts with given names from all variations
func (c Context) GetContext(names []string) (Context, error) {
	ctx, remainingNames, err := c.getContextRecursively(names)
	if err != nil {
		return Context{}, err
	}
	if len(remainingNames) > 0 {
		return Context{}, fmt.Errorf("extraneous context names: %s", strings.Join(remainingNames, " "))
	}
	ctx.Name = strings.Join(names, " ")
	return ctx, nil
}

// getContext returns context resulting from merging contexts with given names from all variations,
// as well as list of remaining names that were not used to resolve variations (used for recursivity)
func (c Context) getContextRecursively(names []string) (Context, []string, error) {
	// Start with context itself
	ctx := c
	ctx.Variations = nil

	for _, variation := range c.Variations {
		if len(names) == 0 {
			return Context{}, nil, fmt.Errorf("too few context names")
		}
		name := names[0]
		names = names[1:]

		// Merge variation context
		variationContext, ok := variation.Contexts[name]
		if !ok {
			return Context{}, nil, fmt.Errorf("context %q not found in variation %q", name, variation.Name)
		}
		ctx = ctx.Merge(variationContext, false)

		// Get child contexts recursively
		if len(variationContext.Variations) > 0 {
			childContext, remainingNames, err := variationContext.getContextRecursively(names)
			if err != nil {
				return Context{}, nil, err
			}
			ctx = ctx.Merge(childContext, false)
			names = remainingNames
		}
	}

	return ctx, names, nil
}

// GetAllImages returns the list of image names referenced in context recursively
func (c Context) GetAllImages() []string {
	namesMap := make(map[string]struct{})

	if c.Image != "" {
		namesMap[c.Image] = struct{}{}
	}

	// Recurse into all child variations/contexts
	for _, variation := range c.Variations {
		for _, ctx := range variation.Contexts {
			childImages := ctx.GetAllImages()
			for _, childImage := range childImages {
				namesMap[childImage] = struct{}{}
			}
		}
	}

	// Sort
	sortedNames := make([]string, 0, len(namesMap))
	for name := range namesMap {
		sortedNames = append(sortedNames, name)
	}
	sort.Strings(sortedNames)
	return sortedNames
}

// String returns a user-friendly yaml representation of this context
func (c Context) String() string {
	buf, err := yaml.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(buf)
}
