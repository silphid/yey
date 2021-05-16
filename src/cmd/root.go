package cmd

import (
	"github.com/silphid/yey/cli/src/cmd/get"
	"github.com/silphid/yey/cli/src/cmd/set"
	"github.com/silphid/yey/cli/src/cmd/start"
	"github.com/silphid/yey/cli/src/cmd/use"
	"github.com/silphid/yey/cli/src/cmd/versioning"
	"github.com/silphid/yey/cli/src/internal/logging"
	"github.com/spf13/cobra"
)

// NewRoot creates the root cobra command
func NewRoot(version string) *cobra.Command {
	c := &cobra.Command{
		Use:          "yey",
		Short:        "A DevOps & CI/CD & Kubernetes-oriented general purpose Docker container with CLI launcher",
		Long:         `A DevOps & CI/CD & Kubernetes-oriented general purpose Docker container with CLI launcher`,
		SilenceUsage: true,
	}

	// var options internal.Options
	c.PersistentFlags().BoolVarP(&logging.Verbose, "verbose", "v", false, "display verbose messages")
	c.AddCommand(get.New())
	c.AddCommand(set.New())
	c.AddCommand(use.New())
	c.AddCommand(start.New())
	c.AddCommand(versioning.New(version))
	return c
}