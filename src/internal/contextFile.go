package yey

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v3"
)

const (
	currentVersion = 0
	yeyRCFileName  = ".yeyrc.yaml"
)

// ContextFile represents yey's current config persisted to disk
type ContextFile struct {
	Version int
	Parent  string
	Path    string `yaml:"-"`
	Context `yaml:",inline"`
	Layers  Layers `yaml:"layers"`
}

// readContextFileFromWorkingDirectory scans the current directory and searches for a .yeyrc.yaml file and returns
// the bytes in the file, the absolute path to contextFile and an error if encountered.
// If none is found it climbs the directory hierarchy.
func readContextFileFromWorkingDirectory() ([]byte, string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, "", err
	}

	for {
		candidate := filepath.Join(wd, yeyRCFileName)
		data, err := os.ReadFile(candidate)

		if errors.Is(err, os.ErrNotExist) {
			if wd == "/" {
				data, candidate, err = readContextFileFromHomeDirectory()
				if err != nil {
					return nil, "", fmt.Errorf("could not find %q in current directory hierarchy or home dir: %w", yeyRCFileName, err)
				}
				return data, candidate, nil
			}
			wd = filepath.Join(wd, "..")
			continue
		}

		if err != nil {
			return nil, "", fmt.Errorf("failed to read context file: %w", err)
		}

		return data, candidate, nil
	}
}

// readContextFileFromHomeDirectory looks in home directory for a .yeyrc.yaml file and returns
// the bytes in the file, the absolute path to contextFile and an error if encountered.
// If none is found it climbs the directory hierarchy.
func readContextFileFromHomeDirectory() ([]byte, string, error) {
	dir, err := homedir.Dir()
	if err != nil {
		return nil, "", fmt.Errorf("could not determine home dir: %w", err)
	}
	file := filepath.Join(dir, yeyRCFileName)
	data, err := os.ReadFile(file)
	if errors.Is(err, os.ErrNotExist) {
		return nil, "", fmt.Errorf("no %q in home dir", yeyRCFileName)
	}
	return data, file, nil
}

// readContextFileFromFilePath reads the contextfile from the fs
func readContextFileFromFilePath(path string) ([]byte, error) {
	path, err := homedir.Expand(path)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(path)
}

// readContextFileFromNetwork reads the contextfile from the network over http
func readContextFileFromNetwork(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed fetching context file from network: %w", err)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// parseContextFile unmarshals the contextFile data and resolves any parent contextfiles
func parseContextFile(dir string, data []byte) (Contexts, error) {
	var ctxFile ContextFile
	if err := yaml.Unmarshal(data, &ctxFile); err != nil {
		return Contexts{}, fmt.Errorf("failed to decode context file: %w", err)
	}

	if ctxFile.Version != currentVersion {
		return Contexts{}, fmt.Errorf("unsupported context file version")
	}

	contexts := Contexts{
		Context: ctxFile.Context,
		Layers:  ctxFile.Layers,
	}

	if dir != "" {
		var err error
		contexts, err = resolveContextsPaths(dir, contexts)
		if err != nil {
			return Contexts{}, err
		}
	}

	if ctxFile.Parent != "" {
		parent, err := readAndParseContextFileFromURI(ctxFile.Parent)
		if err != nil {
			return Contexts{}, fmt.Errorf("failed to resolve parent context %q: %w", ctxFile.Parent, err)
		}
		contexts = parent.Merge(contexts)
	}

	contexts = resolveEnvironmentVariables(contexts)

	return contexts, nil
}

// readAndParseContextFileFromURI reads and parses the context file from an URI, which can either
// be an URL or local path
func readAndParseContextFileFromURI(path string) (Contexts, error) {
	Log("loading context file: %s", path)

	var bytes []byte
	var err error
	var dir string

	if strings.HasPrefix(path, "https:") || strings.HasPrefix(path, "http:") {
		bytes, err = readContextFileFromNetwork(path)
	} else {
		dir = filepath.Dir(path)
		bytes, err = readContextFileFromFilePath(path)
	}

	if err != nil {
		return Contexts{}, fmt.Errorf("failed to read context file: %w", err)
	}

	return parseContextFile(dir, bytes)
}

// LoadContexts reads the context file and returns the contexts. It starts by reading from current
// working directory and resolves all parent context files.
func LoadContexts() (Contexts, error) {
	bytes, path, err := readContextFileFromWorkingDirectory()
	if err != nil {
		return Contexts{}, fmt.Errorf("failed to read context file: %w", err)
	}

	Log("loading context file: %s", path)
	contexts, err := parseContextFile(filepath.Dir(path), bytes)
	if err != nil {
		return Contexts{}, err
	}
	contexts.Path = path

	return contexts, nil
}

func resolveContextsPaths(dir string, contexts Contexts) (Contexts, error) {
	var err error
	contexts.Context, err = resolveContextPaths(dir, contexts.Context)
	if err != nil {
		return Contexts{}, err
	}
	for _, layer := range contexts.Layers {
		for name, context := range layer.Contexts {
			layer.Contexts[name], err = resolveContextPaths(dir, context)
			if err != nil {
				return Contexts{}, err
			}
		}
	}
	return contexts, nil
}

func resolveContextPaths(dir string, context Context) (Context, error) {
	clone := context.Clone()

	// Resolve dockerfile path
	var err error
	clone.Build.Dockerfile, err = resolvePath(dir, context.Build.Dockerfile)
	if err != nil {
		return Context{}, err
	}

	// Resolve build context dir
	clone.Build.Context, err = resolvePath(dir, clone.Build.Context)
	if err != nil {
		return Context{}, err
	}

	// Resolve mount dirs
	clone.Mounts = make(map[string]string, len(context.Mounts))
	for key, value := range context.Mounts {
		key, err = resolvePath(dir, key)
		if err != nil {
			return Context{}, err
		}
		clone.Mounts[key] = value
	}

	return clone, nil
}

func resolvePath(dir, path string) (string, error) {
	if path == "" {
		return "", nil
	}

	// Resolve home dir
	var err error
	if path == "~" {
		path, err = homedir.Dir()
	} else {
		path, err = homedir.Expand(path)
	}
	if err != nil {
		return "", err
	}

	if filepath.IsAbs(path) {
		return path, nil
	}

	return filepath.Join(dir, path), nil
}

func resolveEnvironmentVariables(contexts Contexts) Contexts {
	clone := contexts
	clone.Context = clone.Context.Clone()
	clone.Layers = clone.Layers.Clone()

	for key, value := range clone.Env {
		clone.Context.Env[key] = os.ExpandEnv(value)
	}

	for _, layer := range clone.Layers {
		for _, ctx := range layer.Contexts {
			for key, value := range ctx.Env {
				ctx.Env[key] = os.ExpandEnv(value)
			}
		}
	}

	return clone
}
