package yey

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
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

const rcVerserion = 0

type RCFile struct {
	Context
	Parent   string `yaml:"parent"`
	Version  int    `yaml:"Version"`
	Contexts map[string]Context
}

func ParseRCFile(wd string) (*Contexts, error) {
	var rcBytes []byte
	var err error
	var yeyrcPath string

	for {
		yeyrcPath = filepath.Join(wd, ".yeyrc.yaml")
		rcBytes, err = os.ReadFile(yeyrcPath)
		if errors.Is(err, os.ErrNotExist) {
			if wd == "/" {
				return nil, fmt.Errorf("failed to find .yeyrc.yaml")
			}
			wd = filepath.Join(wd, "..")
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read .yeyrc.yaml: %w", err)
		}
		break
	}

	var rcFile RCFile
	if err := yaml.Unmarshal(rcBytes, &rcFile); err != nil {
		return nil, fmt.Errorf("failed to parse rcfile: %w", err)
	}

	if rcFile.Version != rcVerserion {
		return nil, fmt.Errorf("unsupported version %d (expected %d) in config file %q", rcVerserion, rcVerserion, yeyrcPath)
	}

	if rcFile.Parent != "" {
		// TODO RESOLVE RCFILE
	}

	return &Contexts{
		base:  rcFile.Context,
		named: rcFile.Contexts,
	}, nil
}
