package get

import (
	"github.com/silphid/yey/cli/src/cmd/get/clone"
	"github.com/silphid/yey/cli/src/cmd/get/containers"
	"github.com/silphid/yey/cli/src/cmd/get/context"
	"github.com/silphid/yey/cli/src/cmd/get/contexts"
	"github.com/silphid/yey/cli/src/cmd/get/tag"
	"github.com/spf13/cobra"
)

// New creates a cobra command
func New() *cobra.Command {
	c := &cobra.Command{
		Use:   "get",
		Short: "Displays value(s) of entity or variable",
	}
	c.AddCommand(clone.New())
	c.AddCommand(containers.New())
	c.AddCommand(context.New())
	c.AddCommand(contexts.New())
	c.AddCommand(tag.New())
	return c
}
