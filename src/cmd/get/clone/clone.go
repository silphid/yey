package clone

import (
	"fmt"

	"github.com/silphid/yey/cli/src/internal/core"
	"github.com/spf13/cobra"
)

// New creates a cobra command
func New() *cobra.Command {
	return &cobra.Command{
		Use:   "clone",
		Short: "Displays configured clone directory",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			return run()
		},
	}
}

func run() error {
	c, err := core.New()
	if err != nil {
		return err
	}
	state, err := c.GetState()
	if err != nil {
		return err
	}
	fmt.Println(state.CloneDir)
	return nil
}
