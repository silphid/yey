package start

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
	return &cobra.Command{
		Use:   "start",
		Short: "Starts container",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := ""
			if len(args) == 1 {
				name = args[0]
			}
			return run(cmd.Context(), name)
		},
	}
}

func run(ctx context.Context, name string) error {
	contexts, err := yey.ReadAndParseContextFile()
	if err != nil {
		return err
	}

	if name == "" {
		var err error
		name, err = cmd.PromptContext(contexts)
		if err != nil {
			return err
		}
	}

	yeyContext, err := contexts.GetContext(name)
	if err != nil {
		return fmt.Errorf("failed to get context with name %q: %w", name, err)
	}

	shortImageName, err := docker.GetShortImageName(yeyContext.Image)
	if err != nil {
		return err
	}

	containerName := fmt.Sprintf("yey-%s-%s", shortImageName, yeyContext.Name)

	api, err := docker.NewAPI()
	if err != nil {
		return fmt.Errorf("failed to connect to docker client: %w", err)
	}

	return api.Start(ctx, yeyContext, containerName)
}
