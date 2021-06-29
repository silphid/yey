package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	yey "github.com/silphid/yey/src/internal"
)

// Parses given value into context name and variant and, as needed, prompt user for those values
func GetOrPromptContextNames(contexts yey.Contexts, names []string, lastNames []string) ([]string, error) {
	if len(names) == 1 && names[0] == "-" {
		return lastNames, nil
	}

	availableNames := contexts.GetNamesInAllLayers()

	// Prompt unspecified names
	for i := len(names); i < len(contexts.Layers); i++ {
		// Don't prompt when single name in layer
		if len(availableNames[i]) == 1 {
			names = append(names, availableNames[i][0])
			continue
		}
		prompt := &survey.Select{
			Message: fmt.Sprintf("Select %s", contexts.Layers[i].Name),
			Options: availableNames[i],
		}
		if i < len(lastNames) {
			prompt.Default = lastNames[i]
		}
		selectedIndex := 0
		if err := survey.AskOne(prompt, &selectedIndex); err != nil {
			return nil, err
		}
		names = append(names, availableNames[i][selectedIndex])
	}

	return names, nil
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
