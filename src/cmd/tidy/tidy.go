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
	var options options

	cmd := &cobra.Command{
		Use:   "tidy",
		Short: "Removes unreferenced containers",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.Context(), options)
		},
	}

	cmd.Flags().BoolVarP(&options.force, "force", "f", false, "removes containers forcibly")

	return cmd
}

type options struct {
	force bool
}

func run(ctx context.Context, options options) error {
	contexts, err := yey.LoadContexts()
	if err != nil {
		return err
	}

	validNames := make(map[string]struct{})
	for _, name := range contexts.GetNames() {
		ctx, err := contexts.GetContext(name)
		if err != nil {
			return err
		}
		validNames[yey.ContainerName(contexts.Path, ctx)] = struct{}{}
	}

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

	var removeOptions []docker.RemoveOption
	if options.force {
		removeOptions = append(removeOptions, docker.ForceRemove)
	}

	return docker.RemoveMany(ctx, unreferencedContainers, removeOptions...)
}
