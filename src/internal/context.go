package yey

import (
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

// Context represents execution configuration for some docker container
type Context struct {
	Name       string `yaml:",omitempty"`
	Remove     *bool
	Image      string
	Env        map[string]string
	Mounts     map[string]string
	Cmd        []string
	EntryPoint []string
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

func (c Context) ResolveMounts() (map[string]string, error) {
	mounts := make(map[string]string, len(c.Mounts))
	for key, value := range c.Mounts {
		dir, err := resolveLocalDir(key)
		if err != nil {
			return nil, err
		}
		mounts[dir] = value
	}
	return mounts, nil
}

func resolveLocalDir(dir string) (string, error) {
	var err error
	if dir == "~" {
		dir, err = homedir.Dir()
	} else {
		dir, err = homedir.Expand(dir)
	}
	if err != nil {
		return "", err
	}

	dir, err = filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	return dir, nil
}
