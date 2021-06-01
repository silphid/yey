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
			name := ""
			if len(args) == 1 {
				name = args[0]
			}
			return run(cmd.Context(), name, options)
		},
	}

	cmd.Flags().BoolVarP(&options.force, "force", "f", false, "force removes container")

	return cmd
}

type options struct {
	force bool
}

func run(ctx context.Context, name string, options options) error {
	contexts, err := yey.LoadContexts()
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

	container := yey.ContainerName(contexts.Path, context)

	return docker.Remove(ctx, container, docker.WithForceRemove(options.force))
}
