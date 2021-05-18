package yey

import (
	"fmt"
	"io/ioutil"
	"sort"

	"github.com/silphid/yey/cli/src/internal/helpers"
	"gopkg.in/yaml.v2"
)

const (
	currentVersion  = 0
	BaseContextName = "base"
)

// ContextFile represents yey's current config persisted to disk
type ContextFile struct {
	Context
	Version  int
	Parent   string
	Contexts map[string]Context
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

// Clone returns a deep-copy of this context
func (cf ContextFile) Clone() ContextFile {
	clone := cf
	for key, value := range cf.Contexts {
		clone.Contexts[key] = value.Clone()
	}
	return clone
}

// Merge creates a deep-copy of this context and copies values from given source context on top of it
func (cf ContextFile) Merge(source ContextFile) ContextFile {
	clone := cf.Clone()
	if source.Version != clone.Version {
		panic(fmt.Errorf("trying to merge contexts with incompatible versions %d and %d", clone.Version, source.Version))
	}
	if source.Parent != "" {
		clone.Parent = source.Parent
	}
	clone.Context.Merge(source.Context)
	for key, value := range cf.Contexts {
		existing, ok := cf.Contexts[key]
		if ok {
			cf.Contexts[key] = existing.Merge(value)
		} else {
			cf.Contexts[key] = value
		}
	}
	return clone
}

// GetContextNames returns the list of all context names user can choose from,
// including the special "base" contexts.
func (cf ContextFile) GetContextNames() ([]string, error) {
	// Extract unique names
	namesMap := make(map[string]bool)
	for name := range cf.Contexts {
		namesMap[name] = true
	}

	// Sort
	sortedNames := make([]string, 0, len(namesMap))
	for name := range namesMap {
		sortedNames = append(sortedNames, name)
	}
	sort.Strings(sortedNames)

	// Prepend special contexts
	names := make([]string, 0, len(sortedNames)+2)
	names = append(names, BaseContextName)
	names = append(names, sortedNames...)

	return names, nil
}

// GetContext returns context with given name, or base context
// if name is "base".
func (cf ContextFile) GetContext(name string) (Context, error) {
	if name == BaseContextName {
		return cf.Context, nil
	}
	context, ok := cf.Contexts[name]
	if !ok {
		return Context{}, fmt.Errorf("context %q not found", name)
	}
	return context, nil
}

// MergeAll returns the result of merging all given ContextFiles
func MergeAll(contextFiles []ContextFile) ContextFile {
	merged := ContextFile{}
	for _, file := range contextFiles {
		merged = merged.Merge(file)
	}
	return merged
}
