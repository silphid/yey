package get

import (
	"github.com/spf13/cobra"
)

// New creates a cobra command
func New() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Displays value(s) of entity or variable",
	}
}
