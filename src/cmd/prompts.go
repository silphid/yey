package cmd

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/TwinProduction/go-color"
	yey "github.com/silphid/yey/src/internal"
)

// Parses given value into context name and variant and, as needed, prompt user for those values
func GetOrPromptContextNames(contexts yey.Contexts, names []string, lastNames []string) ([]string, error) {
	if len(names) == 1 && names[0] == "-" {
		return lastNames, nil
	}

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
		if layer < len(lastNames) {
			prompt.Default = lastNames[layer]
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
func GetOrPromptMultipleContextNames(contexts yey.Contexts, providedNames []string, getMatchingContainers func(names [][]string) []string) ([][]string, error) {
	availableNames := contexts.GetNamesInAllLayers()

	outputNames := make([][]string, 0, len(contexts.Layers))
	for layer := 0; layer < len(contexts.Layers); layer++ {
		// Context name for this layer already specified by user?
		if layer < len(providedNames) {
			// Just take name specified by user
			outputNames = append(outputNames, []string{providedNames[layer]})
		} else {
			// Filter context names through predicate
			selectableNames := make([]string, 0, len(availableNames[layer]))
			selectableTitles := make([]string, 0, len(availableNames[layer]))
			for _, name := range availableNames[layer] {
				matchingContainers := len(getMatchingContainers(append(outputNames, []string{name})))
				if matchingContainers > 0 {
					selectableNames = append(selectableNames, name)
					selectableTitles = append(selectableTitles, fmt.Sprintf("%s %s",
						color.Ize(color.Gray, name),
						color.Ize(color.Blue, fmt.Sprintf("(%d)", matchingContainers))))
				}
			}

			layerName := contexts.Layers[layer].Name

			// Don't prompt when single option
			if len(selectableNames) == 1 {
				fmt.Fprintf(os.Stderr, "Auto-selecting only %s: %s\n",
					color.Ize(color.Green, layerName),
					color.Ize(color.Blue, selectableNames[0]))
				outputNames = append(outputNames, selectableNames)
				continue
			}

			// Prompt to multiselect context names for unspecified layer
			prompt := &survey.MultiSelect{
				Message: fmt.Sprintf("Select %s(s)", layerName),
				Options: selectableTitles,
			}
			var selectedIndices []int
			if err := survey.AskOne(prompt, &selectedIndices); err != nil {
				return nil, err
			}

			// Look-up names for selected indices
			selectedNames := make([]string, 0, len(selectedIndices))
			for _, selectedIndex := range selectedIndices {
				selectedNames = append(selectedNames, selectableNames[selectedIndex])
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
