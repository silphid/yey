package containers

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/TwinProduction/go-color"
	"github.com/spf13/cobra"

	yey "github.com/silphid/yey/src/internal"
	"github.com/silphid/yey/src/internal/docker"
)

type Options struct {
	All bool
}

// New creates a cobra command
func New() *cobra.Command {
	var options Options

	cmd := &cobra.Command{
		Use:   "containers",
		Short: "Lists running containers",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.Context(), options)
		},
	}

	cmd.Flags().BoolVarP(&options.All, "all", "a", false, "include all yey containers")

	return cmd
}

func run(ctx context.Context, options Options) error {
	containers, err := docker.ListContainers(ctx, true)
	if err != nil {
		return err
	}
	totalCount := len(containers)

	if !options.All {
		contexts, err := yey.LoadContexts()
		if err != nil {
			return err
		}
		prefix := yey.ContainerPathPrefix(contexts.Path)

		var filteredContainers []string
		for _, container := range containers {
			if strings.HasPrefix(container, prefix) {
				filteredContainers = append(filteredContainers, container)
			}
		}
		containers = filteredContainers
	}

	if len(containers) == 0 {
		if totalCount > 0 {
			fmt.Fprintln(os.Stderr, color.Ize(color.Green, fmt.Sprintf("no project-specific yey containers found, but %d other(s) were found that you could include with --all flag", totalCount)))
			return nil
		}
		fmt.Fprintln(os.Stderr, color.Ize(color.Green, "no yey containers found"))
		return nil
	}
	fmt.Println(strings.Join(containers, "\n"))
	return nil
}
