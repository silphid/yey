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

	containers, err := docker.ListContainers(ctx)
	if err != nil {
		return fmt.Errorf("failed to list containers to prompt for removal: %w", err)
	}

	// Abort if no containers to remove
	if len(containers) == 0 {
		fmt.Fprintln(os.Stderr, color.Ize(color.Green, "no yey containers found to remove"))
		return nil
	}

	// Func to determine which containers match given context names
	getMatchingContainers := func(names [][]string) []string {
		pattern := yey.ContainerNamePattern(names)
		matchingContainers := make([]string, 0, len(containers))
		for _, container := range containers {
			if pattern.MatchString(container) {
				matchingContainers = append(matchingContainers, container)
			}
		}
		return matchingContainers
	}

	// Prompt
	selectedNames, err := cmd.GetOrPromptMultipleContextNames(contexts, names, getMatchingContainers)
	if err != nil {
		return fmt.Errorf("failed to prompt for context: %w", err)
	}

	// Remove all containers matching pattern
	pattern := yey.ContainerNamePattern(selectedNames)
	for _, container := range containers {
		if pattern.MatchString(container) {
			if err := remove(ctx, container, options); err != nil {
				return err
			}
		}
	}
	return nil
}

func remove(ctx context.Context, container string, options docker.RemoveOptions) error {
	yey.Log("Removing %s", container)
	return docker.Remove(ctx, container, options)
}
