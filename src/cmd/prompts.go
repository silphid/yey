package cmd

import (
	"fmt"
	"regexp"

	"github.com/AlecAivazis/survey/v2"
	yey "github.com/silphid/yey/src/internal"
)

// Parses given value into context name and variant and, as needed, prompt user for those values
func GetOrPromptContextNameAndVariant(contexts yey.Contexts, nameAndVariant string) (string, string, error) {
	name, variant, err := parseNameAndVariant(nameAndVariant)
	if err != nil {
		return "", "", err
	}

	if name == "" {
		var err error
		name, err = PromptContext(contexts)
		if err != nil {
			return "", "", fmt.Errorf("failed to prompt for context name: %w", err)
		}
	}

	if variant == "" {
		var err error
		variant, err = PromptVariant(contexts)
		if err != nil {
			return "", "", fmt.Errorf("failed to prompt for context variant: %w", err)
		}
	}

	return name, variant, nil
}

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

func PromptVariant(contexts yey.Contexts) (string, error) {
	// Get variants names
	names, displayNames := contexts.GetVariants()

	// Only one variant defined, no need to prompt, it forcibly is the "none" variant
	if len(names) == 1 {
		return yey.NoneVariantName, nil
	}

	// Show selection prompt
	prompt := &survey.Select{
		Message: "Select variant:",
		Options: displayNames,
	}

	selectedIndex := 0
	if err := survey.AskOne(prompt, &selectedIndex); err != nil {
		return "", err
	}

	return names[selectedIndex], nil
}

var nameAndVariantRegex = regexp.MustCompile(`^([^/]+)(/(.+))?$`)

func parseNameAndVariant(name string) (string, string, error) {
	if name == "" {
		return "", yey.NoneVariantName, nil
	}

	matches := nameAndVariantRegex.FindStringSubmatch(name)
	if matches == nil {
		return "", "", fmt.Errorf("malformed context name: %q", name)
	}

	// Format is "name/variant"
	if len(matches) == 4 {
		return matches[1], matches[3], nil
	}

	// Format is "name" (no variant)
	return matches[1], yey.NoneVariantName, nil
}
