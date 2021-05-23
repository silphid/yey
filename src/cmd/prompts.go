package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	yey "github.com/silphid/yey/src/internal"
)

func PromptContext(contexts yey.Contexts) (string, error) {
	// Get context names
	names := contexts.GetNames()
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
