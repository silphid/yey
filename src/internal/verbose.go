package yey

import (
	"fmt"
	"os"

	"github.com/TwinProduction/go-color"
)

var (
	IsVerbose bool
)

func Log(format string, a ...interface{}) {
	if IsVerbose {
		fmt.Fprintf(os.Stderr, color.Ize(color.Yellow, format+"\n"), a...)
	}
}
