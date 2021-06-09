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
	var options options

	cmd := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"rm"},
		Short:   "Removes context container",
		Args:    cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			nameAndVariant := ""
			if len(args) == 1 {
				nameAndVariant = args[0]
			}
			return run(cmd.Context(), nameAndVariant, options)
		},
	}

	cmd.Flags().BoolVarP(&options.force, "force", "f", false, "force removes container")

	return cmd
}

type options struct {
	force bool
}

func run(ctx context.Context, nameAndVariant string, options options) error {
	contexts, err := yey.LoadContexts()
	if err != nil {
		return err
	}

	name, variant, err := cmd.GetOrPromptContextNameAndVariant(contexts, nameAndVariant)
	if err != nil {
		return fmt.Errorf("failed to prompt for context: %w", err)
	}

	context, err := contexts.GetContext(name, variant)
	if err != nil {
		return fmt.Errorf("failed to get context %q: %w", nameAndVariant, err)
	}

	container := yey.ContainerName(contexts.Path, context)

	return docker.Remove(ctx, container, docker.WithForceRemove(options.force))
}
