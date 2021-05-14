package statefile

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/silphid/yey/cli/src/internal/helpers"
	"gopkg.in/yaml.v2"
)

const (
	currentVersion = "2021.05"
	fileName       = "state.yaml"
)

// State represents yey's current state persisted to disk
type State struct {
	dir            string
	Version        string `yaml:"version"`
	CloneDir       string `yaml:"clone"`
	ImageTag       string `yaml:"tag"`
	CurrentContext string `yaml:"context"`
}

// Save saves state file to given directory
func (s State) Save() error {
	s.Version = currentVersion

	doc, err := yaml.Marshal(s)
	if err != nil {
		return err
	}

	path := filepath.Join(s.dir, fileName)
	return ioutil.WriteFile(path, doc, 0644)
}

// Load loads the state file from given directory
func Load(dir string) (State, error) {
	state := State{
		dir: dir,
	}
	path := filepath.Join(dir, fileName)
	if !helpers.PathExists(path) {
		return state, nil
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return state, fmt.Errorf("loading state file: %w", err)
	}
	err = yaml.Unmarshal(buf, &state)
	if err != nil {
		return state, fmt.Errorf("unmarshalling state file yaml: %w", err)
	}

	if state.Version != currentVersion {
		return state, fmt.Errorf("unsupported state file %q version %s (expected %s)", path, state.Version, currentVersion)
	}

	return state, nil
}
