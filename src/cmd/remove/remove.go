package remove

import (
	"context"
	"fmt"
	"os"
	"strings"

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
		fmt.Fprintln(os.Stderr, color.Ize(color.Green, "no containers to remove"))
		return nil
	}

	// Parse container names to slices
	containerNames := make([][]string, len(containers))
	for i, container := range containers {
		// TODO: improve this logic to support context names with dashes in them
		// We need a more deterministic way to trace back a container name to its context names
		containerNames[i] = strings.Split(container, "-")
	}

	// Predicate to determine whether context name in given layer has a corresponding container
	predicate := func(name string, layer int) bool {
		for _, containerName := range containerNames {
			skipContainerPrefixes := 3
			if containerName[skipContainerPrefixes+layer] == name {
				return true
			}
		}
		return false
	}

	// Prompt
	selectedNames, err := cmd.GetOrPromptMultipleContextNames(contexts, names, predicate)
	if err != nil {
		return fmt.Errorf("failed to prompt for context: %w", err)
	}

	return removeRecursively(ctx, contexts, selectedNames, []string{}, 0, options)
}

func removeRecursively(ctx context.Context, contexts yey.Contexts, selectedNames [][]string, names []string, layer int, options docker.RemoveOptions) error {
	for _, name := range selectedNames[layer] {
		currentNames := append(names, name)
		var err error
		if layer == len(selectedNames)-1 {
			err = remove(ctx, contexts, currentNames, options)
		} else {
			// Recurse
			err = removeRecursively(ctx, contexts, selectedNames, currentNames, layer+1, options)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func remove(ctx context.Context, contexts yey.Contexts, names []string, options docker.RemoveOptions) error {
	context, err := contexts.GetContext(names)
	if err != nil {
		return fmt.Errorf("failed to get context: %w", err)
	}

	container := yey.ContainerName(contexts.Path, context)

	yey.Log("Removing %s", container)
	return docker.Remove(ctx, container, options)
}
