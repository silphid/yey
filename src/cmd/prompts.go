package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	yey "github.com/silphid/yey/src/internal"
)

// Parses given value into context name and variant and, as needed, prompt user for those values
func GetOrPromptContextNames(contexts yey.Contexts, names []string) ([]string, error) {
	availableNames := contexts.GetNamesInAllLayers()

	// Prompt unspecified names
	for layer := len(names); layer < len(contexts.Layers); layer++ {
		// Don't prompt when single name in layer
		if len(availableNames[layer]) == 1 {
			names = append(names, availableNames[layer][0])
			continue
		}
		prompt := &survey.Select{
			Message: fmt.Sprintf("Select %s", contexts.Layers[layer].Name),
			Options: availableNames[layer],
		}
		var selectedName string
		if err := survey.AskOne(prompt, &selectedName); err != nil {
			return nil, err
		}
		names = append(names, selectedName)
	}

	return names, nil
}

// Parses given value into context name and variant and, as needed, prompt user for those values
func GetOrPromptMultipleContextNames(contexts yey.Contexts, names []string, predicate func(name string, layer int) bool) ([][]string, error) {
	availableNames := contexts.GetNamesInAllLayers()

	outputNames := make([][]string, 0, len(contexts.Layers))
	for layer := 0; layer < len(contexts.Layers); layer++ {
		// Context name for this layer already specified by user?
		if layer < len(names) {
			// Just take name specified by user
			outputNames = append(outputNames, []string{names[layer]})
		} else {
			// Filter context names through predicate
			filteredNames := make([]string, 0, len(availableNames[layer]))
			for _, name := range availableNames[layer] {
				if predicate(name, layer) {
					filteredNames = append(filteredNames, name)
				}
			}

			// Don't prompt when single option
			if len(filteredNames) == 1 {
				outputNames = append(outputNames, filteredNames)
				continue
			}

			// Prompt to multiselect context names for unspecified layer
			prompt := &survey.MultiSelect{
				Message: fmt.Sprintf("Select %s(s)", contexts.Layers[layer].Name),
				Options: filteredNames,
			}
			var selectedNames []string
			if err := survey.AskOne(prompt, &selectedNames); err != nil {
				return nil, err
			}
			outputNames = append(outputNames, selectedNames)
		}
	}

	return outputNames, nil
}

// Prompts user to multi-select among given images
func PromptImageNames(allImages []string) ([]string, error) {

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
