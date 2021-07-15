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
		Args:  cobra.ArbitraryArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			return run(args)
		},
	}
}

func run(names []string) error {
	contexts, err := yey.LoadContexts()
	if err != nil {
		return err
	}

	lastNames, err := cmd.LoadLastNames()
	if err != nil {
		return err
	}

	names, err = cmd.GetOrPromptContextNames(contexts.Context, names, lastNames)
	if err != nil {
		return err
	}

	err = cmd.SaveLastNames(names)
	if err != nil {
		return err
	}

	context, err := contexts.GetContext(names)
	if err != nil {
		return fmt.Errorf("failed to get context: %w", err)
	}

	fmt.Println(context.String())
	return nil
}
