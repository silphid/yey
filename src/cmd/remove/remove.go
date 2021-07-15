package remove

import (
	"context"
	"fmt"
	"os"

	"github.com/TwinProduction/go-color"
	"github.com/spf13/cobra"

	"github.com/silphid/yey/src/cmd"
	yey "github.com/silphid/yey/src/internal"
	"github.com/silphid/yey/src/internal/docker"
)

type RemoveOptions struct {
	All   bool
	Force bool
}

// New creates a cobra command
func New() *cobra.Command {
	var options RemoveOptions

	cmd := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"rm"},
		Short:   "Removes context container",
		Args:    cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.Context(), args, options)
		},
	}

	cmd.Flags().BoolVarP(&options.All, "all", "a", false, "include all yey containers")
	cmd.Flags().BoolVarP(&options.Force, "force", "f", false, "force removes container")

	return cmd
}

func run(ctx context.Context, names []string, options RemoveOptions) error {
	contexts, err := yey.LoadContexts()
	if err != nil {
		return err
	}

	containers, err := docker.ListContainers(ctx)
	if err != nil {
		return fmt.Errorf("failed to list containers to prompt for removal: %w", err)
	}

	// Compute all valid contexts
	combos := contexts.GetCombos()
	var validContexts []yey.Context
	for _, combo := range combos {
		context, err := contexts.GetContext(combo)
		if err != nil {
			return fmt.Errorf("failed to get context to prompt for removal: %w", err)
		}
		validContexts = append(validContexts, context)
	}

	// Compute all valid containers
	var validContainers []string
	for _, validContext := range validContexts {
		container := yey.ContainerName(contexts.Path, validContext)

		// Found in list of containers?
		for i := range containers {
			if containers[i] == container {
				validContainers = append(validContainers, container)
				// Remove from list of containers
				containers = append(containers[:i], containers[i+1:]...)
				break
			}
		}
	}

	// Include all containers?
	var otherContainers []string
	if options.All {
		otherContainers = containers
	}

	// Abort if no containers to remove
	if len(validContainers) == 0 && len(otherContainers) == 0 {
		fmt.Fprintln(os.Stderr, color.Ize(color.Green, "no yey containers found to remove"))
		return nil
	}

	// Prompt
	selectedContainers, err := cmd.PromptContainers(validContainers, otherContainers)
	if err != nil {
		return fmt.Errorf("failed to prompt for containers: %w", err)
	}

	// Remove selected containers
	for _, container := range selectedContainers {
		opt := docker.RemoveOptions{Force: options.Force}
		if err := remove(ctx, container, opt); err != nil {
			return err
		}
	}
	return nil
}

func remove(ctx context.Context, container string, options docker.RemoveOptions) error {
	yey.Log("Removing %s", container)
	return docker.Remove(ctx, container, options)
}
