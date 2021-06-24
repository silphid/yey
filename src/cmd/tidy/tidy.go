package tidy

import (
	"context"
	"strings"

	"github.com/spf13/cobra"

	yey "github.com/silphid/yey/src/internal"
	"github.com/silphid/yey/src/internal/docker"
)

// New creates a cobra command
func New() *cobra.Command {
	var options docker.RemoveOptions

	cmd := &cobra.Command{
		Use:   "tidy",
		Short: "Removes unreferenced containers",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.Context(), options)
		},
	}

	cmd.Flags().BoolVarP(&options.Force, "force", "f", false, "removes containers forcibly")

	return cmd
}

func run(ctx context.Context, options docker.RemoveOptions) error {
	contexts, err := yey.LoadContexts()
	if err != nil {
		return err
	}

	validNames := make(map[string]struct{})
	forEachPossibleNameCombination(contexts.GetNamesInAllLayers(), nil, func(names []string) error {
		ctx, err := contexts.GetContext(names)
		if err != nil {
			return err
		}
		validNames[yey.ContainerName(contexts.Path, ctx)] = struct{}{}
		return nil
	})

	prefix := yey.ContainerPathPrefix(contexts.Path)

	names, err := docker.ListContainers(ctx)
	if err != nil {
		return err
	}

	var unreferencedContainers []string
	for _, container := range names {
		if !strings.HasPrefix(container, prefix) {
			continue
		}
		if _, ok := validNames[container]; ok {
			continue
		}
		unreferencedContainers = append(unreferencedContainers, container)
	}

	return docker.RemoveMany(ctx, unreferencedContainers, options)
}

// forEachPossibleNameCombination calls given callback function with each possible combination of names from all layers
func forEachPossibleNameCombination(namesInLayers [][]string, baseCombo []string, fn func([]string) error) error {
	currentDepth := len(baseCombo)
	for _, name := range namesInLayers[currentDepth] {
		currentCombo := append(baseCombo, name)
		if currentDepth == len(namesInLayers)-1 {
			fn(currentCombo)
		} else {
			if err := forEachPossibleNameCombination(namesInLayers, currentCombo, fn); err != nil {
				return err
			}
		}
	}
	return nil
}
