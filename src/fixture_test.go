package main

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
	"github.com/stretchr/testify/require"
)

var rootDir = "_test_assets_"

type ConsoleTester struct {
	c *expect.Console
	r *require.Assertions
}

func (ct ConsoleTester) ExpectString(value string) {
	_, err := ct.c.ExpectString(value)
	ct.r.NoError(err)
}

func (ct ConsoleTester) ExpectEOF() {
	_, err := ct.c.ExpectEOF()
	ct.r.NoError(err)
}

func (ct ConsoleTester) SendLine(value string) {
	_, err := ct.c.SendLine(value)
	ct.r.NoError(err)
}

func ConsoleTest(t *testing.T, ctx context.Context, cmdArgs []string, test func(ct ConsoleTester)) ([]byte, error) {
	t.Helper()
	r := require.New(t)

	buf := new(bytes.Buffer)
	console, _, err := vt10x.NewVT10XConsole(expect.WithStdout(buf))
	require.NoError(t, err)
	ct := ConsoleTester{console, r}

	cmd := exec.CommandContext(ctx, cmdArgs[0], cmdArgs[1:]...)
	cmd.Dir = rootDir

	pr, pw := io.Pipe()
	tty := console.Tty()

	go func() {
		r := io.TeeReader(pr, os.Stdout)
		io.Copy(tty, r)
	}()

	cmd.Stdout, cmd.Stderr, cmd.Stdin = pw, pw, console.Tty()

	donec := make(chan struct{})
	go func() {
		defer close(donec)
		if test != nil {
			test(ct)
		}
	}()

	err = cmd.Run()

	// Close the slave end of the tty, and read the remaining bytes from the master end.
	console.Tty().Close()
	<-donec

	return buf.Bytes(), err
}
