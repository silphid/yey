package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/TwinProduction/go-color"
	yey "github.com/silphid/yey/src/internal"
)

// Parses given value into context name and variant and, as needed, prompt user for those values
func GetOrPromptContexts(context yey.Context, argNames []string, lastNames []string) ([]string, error) {
	names, remainingArgNames, _, err := getOrPromptContextsRecursively(context, argNames, lastNames)
	if err != nil {
		return nil, err
	}
	if len(remainingArgNames) > 0 {
		return nil, fmt.Errorf("extraneous context names: %s", strings.Join(remainingArgNames, " "))
	}
	return names, nil
}

func getOrPromptContextsRecursively(context yey.Context, argNames []string, lastNames []string) ([]string, []string, []string, error) {
	if len(argNames) == 1 && argNames[0] == "-" {
		return lastNames, nil, nil, nil
	}

	selectedNames := make([]string, 0, len(argNames))
	for _, layer := range context.Layers {
		// determine context name for layer
		var selectedName string
		if len(argNames) > 0 {
			// use name passed as argument
			selectedName = argNames[0]
			argNames = argNames[1:]
		} else {
			// prompt for name
			prompt := &survey.Select{
				Message: fmt.Sprintf("Select %s", layer.Name),
			}
			for k := range layer.Contexts {
				prompt.Options = append(prompt.Options, k)
			}
			sort.Strings(prompt.Options)
			if len(lastNames) > 0 {
				prompt.Default = lastNames[0]
			}
			if err := survey.AskOne(prompt, &selectedName); err != nil {
				return nil, nil, nil, err
			}
		}

		// Consume one last name, if any
		if len(lastNames) > 0 {
			lastNames = lastNames[1:]
		}

		selectedNames = append(selectedNames, selectedName)

		// Prompt recursively for layer's own child layers
		selectedContext, ok := layer.Contexts[selectedName]
		if !ok {
			return nil, nil, nil, fmt.Errorf("layer %q has no context %q", layer.Name, selectedName)
		}
		if len(selectedContext.Layers) > 0 {
			var childNames []string
			var err error
			childNames, argNames, lastNames, err = getOrPromptContextsRecursively(selectedContext, argNames, lastNames)
			if err != nil {
				return nil, nil, nil, err
			}
			selectedNames = append(selectedNames, childNames...)
		}
	}

	return selectedNames, argNames, lastNames, nil
}

// Prompts user to multi-select among given images
func PromptImages(allImages []string) ([]string, error) {

	prompt := &survey.MultiSelect{
		Message: "Select images to pull",
		Options: allImages,
	}
	selectedImages := []string{}
	if err := survey.AskOne(prompt, &selectedImages); err != nil {
		return nil, err
	}

	return selectedImages, nil
}

// Prompts user to multi-select among given containers and optionally also
// other containers (which are displayed in yellow)
func PromptContainers(containers []string, otherContainers []string, message string) ([]string, error) {
	// Combine containers and otherContainers (in yellow)
	options := containers
	for _, container := range otherContainers {
		options = append(options, color.Ize(color.Yellow, container))
	}

	prompt := &survey.MultiSelect{
		Message: message,
		Options: options,
	}
	selectedIndices := []int{}
	if err := survey.AskOne(prompt, &selectedIndices); err != nil {
		return nil, err
	}

	// Lookup container names based on indices
	allContainers := append(containers, otherContainers...)
	var selectedContainers []string
	for _, selectedIndex := range selectedIndices {
		selectedContainers = append(selectedContainers, allContainers[selectedIndex])
	}
	return selectedContainers, nil
}
