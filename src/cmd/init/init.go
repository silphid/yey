package tidy

import (
	"context"

	"github.com/spf13/cobra"
)

type options struct {
	Parent string
}

// New creates a cobra command
func New() *cobra.Command {
	var opt options

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Creates a .yeyrc.yaml file in current directory, if doesn't exist",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.Context(), opt)
		},
	}

	cmd.Flags().StringVar(&opt.Parent, "parent", "", "sets parent property to given URL or path")

	return cmd
}

func run(ctx context.Context, opt options) error {
}
