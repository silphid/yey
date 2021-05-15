package ctxfile

import (
	"fmt"
	"io/ioutil"
	"sort"

	"github.com/silphid/yey/cli/src/internal/ctx"
	"github.com/silphid/yey/cli/src/internal/helpers"
	"gopkg.in/yaml.v2"
)

const (
	currentVersion = 0
	ContextBase    = "base"
)

// ContextFile represents yey's current config persisted to disk
type ContextFile struct {
	ctx.Context
	Version       int
	Parent        string
	NamedContexts map[string]ctx.Context
}

// Load loads ContextFile from given file
func Load(file string) (ContextFile, error) {
	var config ContextFile
	if !helpers.PathExists(file) {
		return config, nil
	}

	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return config, fmt.Errorf("loading config file: %w", err)
	}
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		return config, fmt.Errorf("unmarshalling yaml of config file %q: %w", file, err)
	}

	if config.Version != currentVersion {
		return config, fmt.Errorf("unsupported version %d (expected %d) in config file %q", config.Version, currentVersion, file)
	}

	return config, nil
}

// Clone returns a deep-copy of this context
func (cf ContextFile) Clone() ContextFile {
	contextFile := ContextFile{
		Version: cf.Version,
		Parent:  cf.Parent,
	}
	for key, value := range cf.NamedContexts {
		contextFile.NamedContexts[key] = value.Clone()
	}
	return contextFile
}

// Merge creates a deep-copy of this context and copies values from given source context on top of it
func (cf ContextFile) Merge(source ContextFile) ContextFile {
	clone := cf.Clone()
	if source.Version != 0 {
		clone.Version = source.Version
	}
	if source.Parent != "" {
		clone.Parent = source.Parent
	}
	clone.Context.Merge(source.Context)
	for key, value := range cf.NamedContexts {
		existing, ok := cf.NamedContexts[key]
		if ok {
			cf.NamedContexts[key] = existing.Merge(value)
		} else {
			cf.NamedContexts[key] = value
		}
	}
	return clone
}

// GetContextNames returns the list of all context names user can choose from,
// including the special "base" contexts.
func (cf ContextFile) GetContextNames() ([]string, error) {
	// Extract unique names
	namesMap := make(map[string]bool)
	for name := range cf.NamedContexts {
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
	names = append(names, "base")
	names = append(names, sortedNames...)

	return names, nil
}

// GetContext returns context with given name, or base context
// if name is "base".
func (cf ContextFile) GetContext(name string) (ctx.Context, error) {
	if name == "base" {
		return cf.Context, nil
	}
	context, ok := cf.NamedContexts[name]
	if !ok {
		return ctx.Context{}, fmt.Errorf("named context not found: %s", name)
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
