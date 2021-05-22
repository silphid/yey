package containers

import (
	"context"
	"fmt"

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
	api, err := docker.NewAPI()
	if err != nil {
		return fmt.Errorf("failed to connect to docker client: %w", err)
	}
	return api.ListContainers(ctx)
}
