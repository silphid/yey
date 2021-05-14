package core

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/silphid/yey/cli/src/internal/statefile"
)

// UseContext validates that given context is valid and saves it as default context
// to state file. If an empty string is passed, prompts user to select from list of
// available contexts.
func (c Core) UseContext(name string) error {
	if name == "" {
		var err error
		name, err = c.promptContext()
		if err != nil {
			return err
		}
	} else {
		err := c.validateContextName(name)
		if err != nil {
			return err
		}
	}
	state, err := statefile.Load(c.homeDir)
	if err != nil {
		return err
	}
	state.CurrentContext = name
	return state.Save()
}

func (c Core) promptContext() (string, error) {
	// Get context names
	names, err := c.GetContextNames()
	if err != nil {
		return "", err
	}

	// Determine default context index
	defaultIndex := 0
	state, err := statefile.Load(c.homeDir)
	if err != nil {
		return "", err
	}
	for i, n := range names {
		if n == state.CurrentContext {
			defaultIndex = i
			break
		}
	}

	// Show selection prompt
	prompt := &survey.Select{
		Message: "Select context:",
		Options: names,
		Default: defaultIndex,
	}
	selectedIndex := defaultIndex
	if err := survey.AskOne(prompt, &selectedIndex); err != nil {
		return "", err
	}

	return names[selectedIndex], nil
}

// validateContextName returns an error if given context names is invalid
func (c Core) validateContextName(name string) error {
	names, err := c.GetContextNames()
	if err != nil {
		return err
	}
	found := false
	for _, n := range names {
		if n == name {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("context %q invalid (expecting one of: %s)", name, strings.Join(names, ", "))
	}
	return nil
}
