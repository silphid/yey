package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
)

func getLastNamesFilePath() (string, error) {
	return homedir.Expand("~/.yeylast")
}

// SaveLastNames saves to home dir the context names that were last selected by user
func SaveLastNames(names []string) error {
	file, err := getLastNamesFilePath()
	if err != nil {
		return fmt.Errorf("failed to determine last context names file path: %w", err)
	}

	err = ioutil.WriteFile(file, ([]byte)(strings.Join(names, " ")), 0644)
	if err != nil {
		return fmt.Errorf("failed to write last context names file: %w", err)
	}

	return nil
}

// LoadLastNames loads from home dir the context names that were last selected by user
func LoadLastNames() ([]string, error) {
	file, err := getLastNamesFilePath()
	if err != nil {
		return nil, fmt.Errorf("failed to determine last context names file path: %w", err)
	}

	text, err := ioutil.ReadFile(file)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read last context names file: %w", err)
	}

	return strings.Split((string)(text), " "), nil
}
