package context

import (
	"fmt"

	"github.com/silphid/yey/cli/src/internal/core"
	"github.com/spf13/cobra"
)

// New creates a cobra command
func New() *cobra.Command {
	return &cobra.Command{
		Use:   "context",
		Short: "Displays resolved values of given context",
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
	context, err := c.GetContext(name)
	if err != nil {
		return err
	}
	fmt.Println(context.String())
	return nil
}
