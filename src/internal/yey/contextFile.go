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

// readContextFileFromWorkingDirectory scans the current directory and searches for a .yeyrc.yaml file and returns
// the bytes in the file, and an error if encountered. If none is found it climbs the directory hierarchy.
func readContextFileFromWorkingDirectory() ([]byte, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	for {
		candidate := filepath.Join(wd, ".yeyrc.yaml")
		data, err := os.ReadFile(candidate)

		if errors.Is(err, os.ErrNotExist) {
			if wd == "/" {
				return nil, fmt.Errorf("failed to find a .yeyrc in directory hierarchy")
			}
			wd = filepath.Join(wd, "..")
			continue
		}

		if err != nil {
			return nil, fmt.Errorf("failed to read context file: %w", err)
		}

		return data, nil
	}
}

// readContextFileFromFilePath reads the contextfile from the fs
func readContextFileFromFilePath(path string) ([]byte, error) {
	path, err := homedir.Expand(path)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(path)
}

// readContextFileFromNetwork reads the contextfile from the network over http.
func readContextFileFromNetwork(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed fetching context file from network: %w", err)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// parseContextFile unmarshals the contextFile data and resolves any parent contextfiles
func parseContextFile(data []byte) (Contexts, error) {
	var ctxFile ContextFile
	if err := yaml.Unmarshal(data, &ctxFile); err != nil {
		return Contexts{}, fmt.Errorf("failed to decode context file: %w", err)
	}

	if ctxFile.Version != currentVersion {
		return Contexts{}, fmt.Errorf("unsupported context file version")
	}

	contexts := Contexts{
		Context:  ctxFile.Context,
		contexts: ctxFile.contexts,
	}

	if ctxFile.Parent != "" {
		parent, err := readAndParseContextFile(ctxFile.Parent)
		if err != nil {
			return Contexts{}, fmt.Errorf("failed to resolve parent context %q: %w", ctxFile.Parent, err)
		}
		contexts = parent.Merge(contexts)
	}

	return contexts, nil
}

// readAndParseContextFile reads and parses the context file from a path. If empty will work from current working directory
// looking for default .yeyrc.yaml file, if starts with https: will download from network. Otherwise searches path in filesystem
func readAndParseContextFile(path string) (Contexts, error) {
	var bytes []byte
	var err error

	if path == "" {
		bytes, err = readContextFileFromWorkingDirectory()
	} else if strings.HasPrefix(path, "https:") {
		bytes, err = readContextFileFromNetwork(path)
	} else {
		bytes, err = readContextFileFromFilePath(path)
	}

	if err != nil {
		return Contexts{}, fmt.Errorf("failed to read contextfile: %w", err)
	}

	return parseContextFile(bytes)
}

// ReadAndParseContextFile reads the context file and returns the contexts. It starts by reading from current working directory
// and resolves all parent context files
func ReadAndParseContextFile() (Contexts, error) {
	return readAndParseContextFile("")
}
