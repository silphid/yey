package cmd

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	yey "github.com/silphid/yey/src/internal"
)

// Parses given value into context name and variant and, as needed, prompt user for those values
func GetOrPromptContextNames(contexts yey.Contexts, names []string) ([]string, error) {
	availableNames := contexts.GetNamesInAllLayers()

	// Prompt unspecified names
	for i := len(names); i < len(contexts.Layers); i++ {
		// Don't prompt when single name in layer
		if len(availableNames[i]) == 1 {
			names = append(names, availableNames[i][0])
			continue
		}

		// Render all outputs to stderr
		renderer := &survey.Renderer{}
		renderer.WithStdio(terminal.Stdio{
			In:  os.Stdin,
			Out: os.Stderr,
			Err: os.Stderr,
		})

		prompt := &survey.Select{
			Message:  fmt.Sprintf("Select %s:", contexts.Layers[i].Name),
			Renderer: *renderer,
			Options:  availableNames[i],
		}
		selectedIndex := 0
		if err := survey.AskOne(prompt, &selectedIndex); err != nil {
			return nil, err
		}
		names = append(names, availableNames[i][selectedIndex])
	}

	return names, nil
}
