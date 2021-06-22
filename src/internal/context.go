package yey

import (
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
	Remove     *bool
	Image      string
	Build      DockerBuild
	Env        map[string]string
	Mounts     map[string]string
	Cmd        []string
	EntryPoint []string
	Network    string
}

// Clone returns a deep-copy of this context
func (c Context) Clone() Context {
	clone := c
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

// String returns a user-friendly yaml representation of this context
func (c Context) String() string {
	buf, err := yaml.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(buf)
}
