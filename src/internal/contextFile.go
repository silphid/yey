package yey

import (
	"fmt"
	"io/ioutil"

	"github.com/silphid/yey/cli/src/internal/helpers"
	"gopkg.in/yaml.v2"
)

const (
	currentVersion = 0
)

// ContextFile represents yey's current config persisted to disk
type ContextFile struct {
	Version  int
	Parent   string
	Contexts `yaml:",inline"`
}

// Load loads ContextFile from given file
func Load(file string) (ContextFile, error) {
	var cf ContextFile
	if !helpers.PathExists(file) {
		return cf, nil
	}

	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return cf, fmt.Errorf("loading context file: %w", err)
	}

	err = yaml.Unmarshal(buf, &cf)
	if err != nil {
		return cf, fmt.Errorf("unmarshalling yaml of context file %q: %w", file, err)
	}

	if cf.Version != currentVersion {
		return cf, fmt.Errorf("unsupported version %d (expected %d) in context file %q", cf.Version, currentVersion, file)
	}

	return cf, nil
}
