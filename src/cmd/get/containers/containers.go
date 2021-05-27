package containers

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/silphid/yey/src/internal/docker"
)

// New creates a cobra command
func New() *cobra.Command {
	return &cobra.Command{
		Use:   "containers",
		Short: "Lists running containers",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.Context())
		},
	}
}

func run(ctx context.Context) error {
	names, err := docker.ListContainers(ctx)
	if err != nil {
		return err
	}
	fmt.Println(strings.Join(names, "\n"))
	return nil
}
