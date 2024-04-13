package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/bufbuild/protoplugin"
	"github.com/mfridman/buildversion"
	"github.com/mfridman/protoc-gen-go-json/internal/plugin"
)

var version string

func main() {
	runArgs(os.Args[1:])

	ctx, stop := newContext()
	defer stop()

	go func() {
		defer stop()
		if err := protoplugin.Run(
			ctx,
			nil,
			os.Stdin,
			os.Stdout,
			os.Stderr,
			protoplugin.HandlerFunc(plugin.Handle),
		); err != nil {
			fmt.Fprintf(os.Stderr, "protoc-gen-go-json: %v\n", err)
			os.Exit(1)
		}
	}()

	select {
	case <-ctx.Done():
		stop()
	}
}

func runArgs(args []string) error {
	if len(args) > 0 {
		for _, arg := range args {
			switch arg {
			case "--version", "-version":
				fmt.Fprintf(os.Stdout, "protoc-gen-go-json version: %s\n", buildversion.New(version))
				os.Exit(0)
			default:
				fmt.Fprintf(os.Stderr, "protoc-gen-go-json: unknown argument: %s\n", arg)
				os.Exit(1)
			}
		}
	}
	return nil
}

func newContext() (context.Context, context.CancelFunc) {
	signals := []os.Signal{os.Interrupt}
	if runtime.GOOS != "windows" {
		signals = append(signals, syscall.SIGTERM)
	}
	return signal.NotifyContext(context.Background(), signals...)
}
