package yey

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

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
