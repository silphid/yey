package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/silphid/yey/src/cmd"

	"github.com/silphid/yey/src/cmd/get"

	getcontainers "github.com/silphid/yey/src/cmd/get/containers"
	getcontext "github.com/silphid/yey/src/cmd/get/context"
	getcontexts "github.com/silphid/yey/src/cmd/get/contexts"

	"github.com/silphid/yey/src/cmd/start"
	"github.com/silphid/yey/src/cmd/versioning"
)

var version string

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		<-ctx.Done()
		stop()
	}()

	rootCmd := cmd.NewRoot()
	rootCmd.AddCommand(start.New())
	rootCmd.AddCommand(versioning.New(version))

	getCMD := get.New()
	getCMD.AddCommand(getcontext.New())
	getCMD.AddCommand(getcontexts.New())
	getCMD.AddCommand(getcontainers.New())

	rootCmd.AddCommand(getCMD)

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(-1)
	}
}
