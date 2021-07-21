package contexts

import (
	"fmt"
	"strings"

	yey "github.com/silphid/yey/src/internal"

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
	contexts, err := yey.LoadContexts()
	if err != nil {
		return err
	}

	combos := contexts.GetCombos()
	for _, combo := range combos {
		fmt.Println(strings.Join(combo, " "))
	}
	return nil
}
