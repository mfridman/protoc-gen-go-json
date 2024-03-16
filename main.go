package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"strings"
	"syscall"
	"time"

	"github.com/bufbuild/protoplugin"
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
				var version string
				if version == "" {
					version = getVersionFromBuildInfo()
				}
				fmt.Fprintf(os.Stdout, "protoc-gen-go-json version: %s\n", strings.TrimSpace(version))
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

// getVersionFromBuildInfo returns the version string from the build info, if available. It will
// always return a non-empty string.
//
//   - If the build info is not available, it returns "devel".
//   - If the main version is set, it returns the string as is.
//   - If building from source, it returns "devel" followed by the first 12 characters of the VCS
//     revision, followed by ", dirty" if the working directory was dirty. For example,
//     "devel (abcdef012345, dirty)" or "devel (abcdef012345)". If the VCS revision is not available,
//     "unknown revision" is used instead.
func getVersionFromBuildInfo() string {
	const defaultVersion = "devel"

	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		// Should only happen if -buildvcs=false is set or using a really old version of Go.
		return defaultVersion
	}
	// The (devel) string is not documented, but it is the value used by the Go toolchain. See
	// https://github.com/golang/go/issues/29228
	if s := buildInfo.Main.Version; s != "" && s != "(devel)" {
		return buildInfo.Main.Version
	}
	var vcs struct {
		revision string
		time     time.Time
		modified bool
	}
	for _, setting := range buildInfo.Settings {
		switch setting.Key {
		case "vcs.revision":
			vcs.revision = setting.Value
		case "vcs.time":
			vcs.time, _ = time.Parse(time.RFC3339, setting.Value)
		case "vcs.modified":
			vcs.modified = (setting.Value == "true")
		}
	}

	var b strings.Builder
	b.WriteString(defaultVersion)
	b.WriteString(" (")
	if vcs.revision == "" || len(vcs.revision) < 12 {
		b.WriteString("unknown revision")
	} else {
		b.WriteString(vcs.revision[:12])
	}
	if vcs.modified {
		b.WriteString(", dirty")
	}
	b.WriteString(")")
	return b.String()
}
