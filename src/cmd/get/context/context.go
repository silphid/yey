package context

import (
	"fmt"

	"github.com/silphid/yey/src/cmd"
	yey "github.com/silphid/yey/src/internal"
	"github.com/spf13/cobra"
)

// New creates a cobra command
func New() *cobra.Command {
	return &cobra.Command{
		Use:   "context",
		Short: "Displays resolved values of given context",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			name := ""
			if len(args) == 1 {
				name = args[0]
			}
			return run(name)
		},
	}
}

func run(name string) error {
	contexts, err := yey.ReadAndParseContextFile()
	if err != nil {
		return err
	}

	if name == "" {
		var err error
		name, err = cmd.PromptContext(contexts)
		if err != nil {
			return fmt.Errorf("failed to prompt for desired context: %w", err)
		}
	}

	context, err := contexts.GetContext(name)
	if err != nil {
		return fmt.Errorf("failed to get context with name %q: %w", name, err)
	}

	fmt.Println(context.String())
	return nil
}
