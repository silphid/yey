package yey

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

const rcVerserion = 0

type RCFile struct {
	Context
	Parent   string `yaml:"parent"`
	Version  int    `yaml:"Version"`
	Contexts map[string]Context
}

func ParseRCFile(ctx context.Context, resource string) (*Contexts, error) {
	var rcBytes []byte
	var err error

	if resource == "" {
		rcBytes, err = readRCFromWorkingDirectory()
	} else if strings.HasPrefix(resource, "http:") || strings.HasPrefix(resource, "https:") {
		rcBytes, err = readRCFromNetwork(ctx, resource)
	} else {
		rcBytes, err = readRCFromFilepath(resource)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to read rcfile: %w", err)
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

func readRCFromWorkingDirectory() ([]byte, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	for {
		rcBytes, err := os.ReadFile(filepath.Join(wd, ".yeyrc.yaml"))
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
		return rcBytes, nil
	}
}

func readRCFromFilepath(path string) ([]byte, error) {
	return os.ReadFile(filepath.Clean(path))
}

func readRCFromNetwork(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to form request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %w", err)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
