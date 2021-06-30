package remove

import (
	"context"
	"fmt"
	"os"
	"regexp"
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
		fmt.Fprintln(os.Stderr, color.Ize(color.Green, "no yey containers found to remove"))
		return nil
	}

	// Predicate to determine whether context name in given layer has a corresponding container
	predicate := func(name string) bool {
		for _, container := range containers {
			if strings.Contains(container, "-"+name+"-") {
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

	// Find regex patterns for all selected names
	patterns := getPatternsRecursively(ctx, contexts, selectedNames, []string{}, 0, options)

	// Remove all containers matching patterns
	for _, container := range containers {
		for _, pattern := range patterns {
			if pattern.MatchString(container) {
				if err := remove(ctx, container, options); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func getPatternsRecursively(ctx context.Context, contexts yey.Contexts, selectedNames [][]string, names []string, layer int, options docker.RemoveOptions) []*regexp.Regexp {
	patterns := []*regexp.Regexp{}
	for _, name := range selectedNames[layer] {
		currentNames := append(names, name)
		if layer == len(selectedNames)-1 {
			patterns = append(patterns, yey.ContainerNamePattern(currentNames))
		} else {
			// Recurse
			patterns = append(patterns, getPatternsRecursively(ctx, contexts, selectedNames, currentNames, layer+1, options)...)
		}
	}
	return patterns
}

func remove(ctx context.Context, container string, options docker.RemoveOptions) error {
	yey.Log("Removing %s", container)
	return docker.Remove(ctx, container, options)
}
