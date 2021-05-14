package contexts

import (
	"fmt"
	"strings"

	"github.com/silphid/yey/cli/src/internal/core"
	"github.com/spf13/cobra"
)

// New creates a cobra command
func New() *cobra.Command {
	return &cobra.Command{
		Use:   "contexts",
		Short: "Lists available contexts",
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
	names, err := c.GetContextNames()
	if err != nil {
		return err
	}
	fmt.Println(strings.Join(names, "\n"))
	return nil
}
