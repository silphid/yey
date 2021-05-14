package use

import (
	"github.com/silphid/yey/cli/src/internal/core"
	"github.com/spf13/cobra"
)

// New creates a cobra command
func New() *cobra.Command {
	return &cobra.Command{
		Use:   "use",
		Short: "Sets current context to use by default",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(_ *cobra.Command, args []string) error {
			name := ""
			if len(args) == 1 {
				name = args[0]
			}
			return run(name)
		},
	}
}

func run(name string) error {
	c, err := core.New()
	if err != nil {
		return err
	}
	return c.UseContext(name)
}
