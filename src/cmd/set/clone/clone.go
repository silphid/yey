package clone

import (
	"github.com/silphid/yey/cli/src/internal/core"
	"github.com/spf13/cobra"
)

// New creates a cobra command
func New() *cobra.Command {
	return &cobra.Command{
		Use:   "clone",
		Short: "Sets 'clone' yey state variable (clone directory) to given value",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return run(args[0])
		},
	}
}

func run(cloneDir string) error {
	c, err := core.New()
	if err != nil {
		return err
	}
	state, err := c.GetState()
	if err != nil {
		return err
	}
	state.CloneDir = cloneDir
	return state.Save()
}
