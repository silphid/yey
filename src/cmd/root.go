package cmd

import (
	yey "github.com/silphid/yey/src/internal"
	"github.com/spf13/cobra"
)

// NewRoot creates the root cobra command
func NewRoot() *cobra.Command {
	c := &cobra.Command{
		Use:          "yey",
		Short:        "An interactive, human-friendly docker launcher for dev and devops",
		Long:         "An interactive, human-friendly docker launcher for dev and devops",
		SilenceUsage: true,
	}

	c.PersistentFlags().BoolVarP(&yey.IsVerbose, "verbose", "v", false, "output verbose messages to stderr")
	c.PersistentFlags().BoolVar(&yey.IsDryRun, "dry-run", false, "output docker command to stdout instead of executing it")

	return c
}
