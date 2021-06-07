package run

import (
	_ "embed"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/arsham/figurine/figurine"
	_ "github.com/arsham/figurine/statik"
	"github.com/arsham/rainbow/rainbow"
)

//go:embed "banner.txt"
var banner []byte

func ShowBanner(contextName string) error {
	// Show "yey!" banner
	rand.Seed(time.Now().UTC().UnixNano())
	l := &rainbow.Light{
		Writer: os.Stdout,
		Seed:   rand.Int63n(256),
	}
	if _, err := l.Write(banner); err != nil {
		return err
	}

	// Show context name
	if contextName != "base" {
		if err := figurine.Write(os.Stdout, contextName, "Small.flf"); err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}
