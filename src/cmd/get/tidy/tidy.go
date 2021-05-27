package tidy

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	yey "github.com/silphid/yey/src/internal"
	"github.com/silphid/yey/src/internal/docker"
)

// New creates a cobra command
func New() *cobra.Command {
	return &cobra.Command{
		Use:   "tidy",
		Short: "cleans unreferenced project containers",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.Context())
		},
	}
}

func run(ctx context.Context) error {
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

	prefix := fmt.Sprintf(
		"yey-%s-%s",
		filepath.Base(filepath.Dir(contexts.Path)),
		yey.Hash(contexts.Path),
	)

	names, err := docker.ListContainers(ctx)
	if err != nil {
		return err
	}

	for _, container := range names {
		if !strings.HasPrefix(container, prefix) {
			continue
		}
		if _, ok := validNames[container]; ok {
			continue
		}
		if err := docker.Remove(ctx, container); err != nil {
			return err
		}
	}

	return nil
}
