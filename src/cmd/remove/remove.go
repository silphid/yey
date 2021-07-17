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

	containers, err := docker.ListContainers(ctx, true)
	if err != nil {
		return fmt.Errorf("failed to list containers to prompt for removal: %w", err)
	}
	totalCount := len(containers)

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
		if totalCount > 0 {
			fmt.Fprintln(os.Stderr, color.Ize(color.Green, fmt.Sprintf("no project-specific yey containers found, but %d other(s) were found that you could include with --all flag", totalCount)))
			return nil
		}
		fmt.Fprintln(os.Stderr, color.Ize(color.Green, "no yey containers found to remove"))
		return nil
	}

	// Prompt
	selectedContainers, err := cmd.PromptContainers(validContainers, otherContainers, "Select containers to remove")
	if err != nil {
		return fmt.Errorf("failed to prompt for containers: %w", err)
	}

	// Prompt user to confirm force removing currently running containers
	if !options.Force {
		runningContainers, err := getRunningContainers(ctx, selectedContainers)
		if err != nil {
			return err
		}

		if len(runningContainers) > 0 {
			forceRemoveContainers, err := cmd.PromptContainers(runningContainers, nil, color.Ize(color.Red, "Select which of the following currently running containers you really want to force remove"))
			if err != nil {
				return err
			}
			selectedContainers = subtractStrings(selectedContainers, runningContainers)

			// Force remove user-confirmed running containers
			for _, container := range forceRemoveContainers {
				opt := docker.RemoveOptions{Force: true}
				if err := remove(ctx, container, opt); err != nil {
					return err
				}
			}
		}
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

func getRunningContainers(ctx context.Context, containers []string) ([]string, error) {
	runningContainers, err := docker.ListContainers(ctx, true)
	if err != nil {
		return nil, fmt.Errorf("failed to list running containers: %w", err)
	}

	intersection := intersectStrings(containers, runningContainers)
	return intersection, nil
}

func subtractStrings(superset []string, subset []string) []string {
	var results []string
	for _, value := range superset {
		if !stringIsInStrings(value, subset) {
			results = append(results, value)
		}
	}
	return results
}

func intersectStrings(set1 []string, set2 []string) []string {
	var results []string
	for _, value1 := range set1 {
		if stringIsInStrings(value1, set2) {
			results = append(results, value1)
		}
	}
	return results
}

func stringIsInStrings(candidate string, values []string) bool {
	for _, value := range values {
		if value == candidate {
			return true
		}
	}
	return false
}

func remove(ctx context.Context, container string, options docker.RemoveOptions) error {
	yey.Log("Removing %s", container)
	return docker.Remove(ctx, container, options)
}
