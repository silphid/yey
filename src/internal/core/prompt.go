package core

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

func (c Core) GetOrPromptContextName(name string) (string, error) {
	if name == "" {
		return c.promptContext()
	}
	if err := c.validateContextName(name); err != nil {
		return "", err
	}
	return name, nil
}

func (c Core) promptContext() (string, error) {
	// Get context names
	names, err := c.GetContextNames()
	if err != nil {
		return "", err
	}
	if len(names) == 0 {
		return "", fmt.Errorf("no context defined")
	}

	// Only one context defined, no need to prompt
	if len(names) == 1 {
		return names[0], nil
	}

	// Show selection prompt
	prompt := &survey.Select{
		Message: "Select context:",
		Options: names,
	}
	selectedIndex := 0
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
