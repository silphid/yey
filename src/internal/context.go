package yey

import (
	"encoding/hex"
	"hash/crc64"
	"io"

	"gopkg.in/yaml.v2"
)

// Context represents execution configuration for some docker container
type Context struct {
	Name   string `yaml:",omitempty"`
	Image  string
	Env    map[string]string
	Mounts map[string]string
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
	merged := c.Clone()
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

func (c Context) Hash() string {
	hasher := crc64.New(crc64.MakeTable(crc64.ECMA))
	io.WriteString(hasher, c.String())
	return hex.EncodeToString(hasher.Sum(nil))
}
