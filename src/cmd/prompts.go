package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	yey "github.com/silphid/yey/src/internal"
)

// Parses given value into context name and variant and, as needed, prompt user for those values
func GetOrPromptContextNames(context yey.Context, argNames []string, lastNames []string) ([]string, error) {
	names, remainingArgNames, _, err := getOrPromptContextNamesRecursively(context, argNames, lastNames)
	if err != nil {
		return nil, err
	}
	if len(remainingArgNames) > 0 {
		return nil, fmt.Errorf("extraneous context names: %s", strings.Join(remainingArgNames, " "))
	}
	return names, nil
}

func getOrPromptContextNamesRecursively(context yey.Context, argNames []string, lastNames []string) ([]string, []string, []string, error) {
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
			childNames, argNames, lastNames, err = getOrPromptContextNamesRecursively(selectedContext, argNames, lastNames)
			if err != nil {
				return nil, nil, nil, err
			}
			selectedNames = append(selectedNames, childNames...)
		}
	}

	return selectedNames, argNames, lastNames, nil
}

// // Parses given value into context name and variant and, as needed, prompt user for those values
// func GetOrPromptMultipleContextNames(contexts yey.Contexts, providedNames []string, getMatchingContainers func(names [][]string) []string) ([][]string, error) {
// 	availableNames := contexts.GetNamesInAllLayers()

// 	outputNames := make([][]string, 0, len(contexts.Layers))
// 	for layer := 0; layer < len(contexts.Layers); layer++ {
// 		// Context name for this layer already specified by user?
// 		if layer < len(providedNames) {
// 			// Just take name specified by user
// 			outputNames = append(outputNames, []string{providedNames[layer]})
// 		} else {
// 			// Filter context names through predicate
// 			selectableNames := make([]string, 0, len(availableNames[layer]))
// 			selectableTitles := make([]string, 0, len(availableNames[layer]))
// 			for _, name := range availableNames[layer] {
// 				matchingContainers := len(getMatchingContainers(append(outputNames, []string{name})))
// 				if matchingContainers > 0 {
// 					selectableNames = append(selectableNames, name)
// 					selectableTitles = append(selectableTitles, fmt.Sprintf("%s %s",
// 						color.Ize(color.Gray, name),
// 						color.Ize(color.Blue, fmt.Sprintf("(%d)", matchingContainers))))
// 				}
// 			}

// 			layerName := contexts.Layers[layer].Name

// 			// Don't prompt when single option
// 			if len(selectableNames) == 1 {
// 				fmt.Fprintf(os.Stderr, "Auto-selecting only %s: %s\n",
// 					color.Ize(color.Green, layerName),
// 					color.Ize(color.Blue, selectableNames[0]))
// 				outputNames = append(outputNames, selectableNames)
// 				continue
// 			}

// 			// Prompt to multiselect context names for unspecified layer
// 			prompt := &survey.MultiSelect{
// 				Message: fmt.Sprintf("Select %s(s)", layerName),
// 				Options: selectableTitles,
// 			}
// 			var selectedIndices []int
// 			if err := survey.AskOne(prompt, &selectedIndices); err != nil {
// 				return nil, err
// 			}

// 			// Look-up names for selected indices
// 			selectedNames := make([]string, 0, len(selectedIndices))
// 			for _, selectedIndex := range selectedIndices {
// 				selectedNames = append(selectedNames, selectableNames[selectedIndex])
// 			}
// 			outputNames = append(outputNames, selectedNames)
// 		}
// 	}

// 	return outputNames, nil
// }

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
