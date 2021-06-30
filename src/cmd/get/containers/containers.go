package containers

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/TwinProduction/go-color"
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
	if len(names) == 0 {
		fmt.Fprintln(os.Stderr, color.Ize(color.Green, "no yey containers found"))
		return nil
	}
	fmt.Println(strings.Join(names, "\n"))
	return nil
}
