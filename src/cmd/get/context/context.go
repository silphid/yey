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
			nameAndVariant := ""
			if len(args) == 1 {
				nameAndVariant = args[0]
			}
			return run(nameAndVariant)
		},
	}
}

func run(nameAndVariant string) error {
	contexts, err := yey.LoadContexts()
	if err != nil {
		return err
	}

	name, variant, err := cmd.GetOrPromptContextNameAndVariant(contexts, nameAndVariant)
	if err != nil {
		return err
	}

	context, err := contexts.GetContext(name, variant)
	if err != nil {
		return fmt.Errorf("failed to get context with name %q: %w", name, err)
	}

	fmt.Println(context.String())
	return nil
}
