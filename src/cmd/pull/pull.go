package pull

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/TwinProduction/go-color"
	"github.com/silphid/yey/src/cmd"
	yey "github.com/silphid/yey/src/internal"
	"github.com/silphid/yey/src/internal/docker"
)

type pullOptions struct {
	all bool
}

// New creates a cobra command
func New() *cobra.Command {
	var options pullOptions

	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull image(s) from registry",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.Context(), args, options)
		},
	}

	cmd.Flags().BoolVarP(&options.all, "all", "a", false, "pull all images")

	return cmd
}

func run(ctx context.Context, names []string, options pullOptions) error {
	contexts, err := yey.LoadContexts()
	if err != nil {
		return err
	}

	// Determine which images to pull
	images := contexts.GetAllImages()
	if !options.all {
		images, err = cmd.PromptImages(images)
		if err != nil {
			return fmt.Errorf("failed to prompt images to pull: %w", err)
		}
	}

	// Pull selected images
	for _, image := range images {
		fmt.Fprintf(os.Stderr, color.Ize(color.Green, "Pulling %s\n"), image)
		if err := docker.Pull(ctx, image); err != nil {
			return err
		}
	}

	return nil
}
