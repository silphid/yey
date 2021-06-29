package remove

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/silphid/yey/src/cmd"
	yey "github.com/silphid/yey/src/internal"
	"github.com/silphid/yey/src/internal/docker"
)

// New creates a cobra command
func New() *cobra.Command {
	var options docker.RemoveOptions

	cmd := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"rm"},
		Short:   "Removes context container",
		Args:    cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.Context(), args, options)
		},
	}

	cmd.Flags().BoolVarP(&options.Force, "force", "f", false, "force removes container")

	return cmd
}

func run(ctx context.Context, names []string, options docker.RemoveOptions) error {
	contexts, err := yey.LoadContexts()
	if err != nil {
		return err
	}

	lastNames, err := cmd.LoadLastNames()
	if err != nil {
		return err
	}

	names, err = cmd.GetOrPromptContextNames(contexts, names, lastNames)
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

	container := yey.ContainerName(contexts.Path, context)

	return docker.Remove(ctx, container, options)
}
