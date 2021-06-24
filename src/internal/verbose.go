package yey

import (
	"fmt"
	"os"
)

var (
	IsVerbose bool
)

func Log(format string, a ...interface{}) {
	if IsVerbose {
		fmt.Fprintf(os.Stderr, format+"\n", a...)
	}
}
