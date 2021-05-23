package cmd

import (
	"github.com/silphid/yey/src/internal/logging"
	"github.com/spf13/cobra"
)

// NewRoot creates the root cobra command
func NewRoot() *cobra.Command {
	c := &cobra.Command{
		Use:          "yey",
		Short:        "A DevOps & CI/CD & Kubernetes-oriented general purpose Docker container with CLI launcher",
		Long:         `A DevOps & CI/CD & Kubernetes-oriented general purpose Docker container with CLI launcher`,
		SilenceUsage: true,
	}

	// var options internal.Options
	c.PersistentFlags().BoolVarP(&logging.Verbose, "verbose", "v", false, "display verbose messages")
	// c.AddCommand(get.New())
	// c.AddCommand(start.New())
	// c.AddCommand(versioning.New(version))
	return c
}
