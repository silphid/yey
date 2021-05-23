package main

import (
	"os"

	"github.com/silphid/yey/src/cmd"
	"github.com/silphid/yey/src/cmd/get"
	"github.com/silphid/yey/src/cmd/start"
	"github.com/silphid/yey/src/cmd/versioning"
)

var version string

func main() {
	rootCmd := cmd.NewRoot()
	rootCmd.AddCommand(start.New())
	rootCmd.AddCommand(get.New())
	rootCmd.AddCommand(versioning.New(version))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
