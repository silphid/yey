package run

import (
	"context"
	"fmt"

	yey "github.com/silphid/yey/src/internal"
	"github.com/silphid/yey/src/internal/docker"

	"github.com/silphid/yey/src/cmd"

	"github.com/spf13/cobra"
)

// New creates a cobra command
func New() *cobra.Command {
	options := Options{Remove: new(bool)}

	cmd := &cobra.Command{
		Use:   "run",
		Short: "runs container",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := ""
			if len(args) == 1 {
				name = args[0]
			}
			if !cmd.Flag("rm").Changed {
				options.Remove = nil
			}
			return run(cmd.Context(), name, options)
		},
	}

	cmd.Flags().BoolVar(options.Remove, "rm", false, "remove container upon exit")
	cmd.Flags().BoolVar(&options.Reset, "reset", false, "remove previous container before starting a fresh one")

	return cmd
}

type Options struct {
	Remove *bool
	Reset  bool
}

func run(ctx context.Context, name string, options Options) error {
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

	yeyContext, err := contexts.GetContext(name)
	if err != nil {
		return fmt.Errorf("failed to get context with name %q: %w", name, err)
	}
	if options.Remove != nil {
		yeyContext.Remove = options.Remove
	}

	containerName := yey.ContainerName(contexts.Path, yeyContext)

	if options.Reset {
		if err := docker.Remove(ctx, containerName); err != nil {
			return fmt.Errorf("failed to remove container %q: %w", containerName, err)
		}
	}

	return docker.Start(ctx, yeyContext, containerName)
}
