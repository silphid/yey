package set

import (
	"github.com/silphid/yey/cli/src/cmd/set/clone"
	"github.com/spf13/cobra"
)

// New creates a cobra command
func New() *cobra.Command {
	c := &cobra.Command{
		Use:   "set",
		Short: "Sets a yey state variable to given value",
	}
	c.AddCommand(clone.New())
	return c
}
