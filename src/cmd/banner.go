package cmd

import (
	_ "embed"
	"math/rand"
	"os"
	"time"

	"github.com/arsham/rainbow/rainbow"
)

//go:embed "banner.txt"
var banner []byte

func ShowBanner(contextName string) error {
	rand.Seed(time.Now().UTC().UnixNano())
	l := &rainbow.Light{
		Writer: os.Stdout,
		Seed:   rand.Int63n(256),
	}
	if _, err := l.Write(banner); err != nil {
		return err
	}
	return nil
}
