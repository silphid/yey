package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	yey "github.com/silphid/yey/src/internal"
)

// Parses given value into context name and variant and, as needed, prompt user for those values
func GetOrPromptContextNames(contexts yey.Contexts, names []string) ([]string, error) {
	availableNames := contexts.GetNames()

	// Prompt unspecified names
	for i := len(names); i < len(contexts.Layers); i++ {
		prompt := &survey.Select{
			Message: fmt.Sprintf("Select %s:", contexts.Layers[i].Name),
			Options: availableNames[i],
		}
		selectedIndex := 0
		if err := survey.AskOne(prompt, &selectedIndex); err != nil {
			return nil, err
		}
		names = append(names, availableNames[i][selectedIndex])
	}

	return names, nil
}
