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
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.Context(), options)
		},
	}

	cmd.Flags().BoolVarP(&options.all, "all", "a", false, "pull all images")

	return cmd
}

func run(ctx context.Context, options pullOptions) error {
	contexts, err := yey.LoadContexts()
	if err != nil {
		return err
	}

	// Determine which images to pull
	imagesAndPlatforms := contexts.GetAllImagesAndPlatforms("")
	if !options.all {
		imagesAndPlatforms, err = cmd.PromptImagesAndPlatforms(imagesAndPlatforms)
		if err != nil {
			return fmt.Errorf("failed to prompt images to pull: %w", err)
		}
	}

	// Pull selected images
	for _, item := range imagesAndPlatforms {
		fmt.Fprintf(os.Stderr, color.Ize(color.Green, "Pulling %s\n"), item)
		if err := docker.Pull(ctx, item.Image, item.Platform); err != nil {
			return err
		}
	}

	return nil
}
