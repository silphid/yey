package containers

import (
	"github.com/silphid/yey/cli/src/internal/docker/api"
	"github.com/spf13/cobra"
)

// New creates a cobra command
func New() *cobra.Command {
	return &cobra.Command{
		Use:   "containers",
		Short: "Lists running containers",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			return run()
		},
	}
}

func run() error {
	return api.ListContainers()
}
