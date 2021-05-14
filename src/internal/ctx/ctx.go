package ctx

import (
	"fmt"
	"io/ioutil"

	"github.com/silphid/yey/cli/src/internal/helpers"
	"gopkg.in/yaml.v2"
)

// RegistryType represents the type of docker registry yey image should be retrieved from
type RegistryType string

const (
	RegistryGCR       RegistryType = "gcr"
	RegistryECR       RegistryType = "ecr"
	RegistryDockerHub RegistryType = "dockerhub"
)

// Context represents an execution context for yey (env vars and volumes)
type Context struct {
	Name      string `yaml:"-"`
	Registry  RegistryType
	Image     string
	Container string
	Env       map[string]string
	Mounts    map[string]string
}

var None = Context{
	Env:    make(map[string]string),
	Mounts: make(map[string]string),
}

// Load loads the context from given file
func Load(file string) (Context, error) {
	var context Context
	if !helpers.PathExists(file) {
		return context, nil
	}

	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return context, fmt.Errorf("loading context file: %w", err)
	}
	err = yaml.Unmarshal(buf, &context)
	if err != nil {
		return context, fmt.Errorf("unmarshalling yaml of context file %q: %w", file, err)
	}

	return context, nil
}

// Clone returns a deep-copy of this context
func (c Context) Clone() Context {
	context := Context{
		Registry:  c.Registry,
		Container: c.Container,
		Name:      c.Name,
		Image:     c.Image,
		Env:       make(map[string]string),
		Mounts:    make(map[string]string),
	}

	for key, value := range c.Env {
		context.Env[key] = value
	}

	// Mounts
	for key, value := range c.Mounts {
		context.Mounts[key] = value
	}

	return context
}

// Merge creates a deep-copy of this context and copies values from given source context on top of it
func (c Context) Merge(source Context) Context {
	context := c.Clone()

	// Registry
	if source.Registry != "" {
		context.Registry = source.Registry
	}

	// Container
	if source.Container != "" {
		context.Container = source.Container
	}

	// Image
	if source.Image != "" {
		context.Image = source.Image
	}

	// Env
	for key, value := range source.Env {
		context.Env[key] = value
	}

	// Volumes
	for key, value := range source.Mounts {
		context.Mounts[key] = value
	}

	return context
}
