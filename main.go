package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/bufbuild/protoplugin"
	"github.com/mfridman/protoc-gen-go-json/internal/plugin"
)

func main() {
	// Handle the version flag here to keep the existing "-version" alias
	// and output format. protoplugin.Main handles everything else: signals,
	// unknown arguments, and the exit code.
	if len(os.Args) == 2 && (os.Args[1] == "--version" || os.Args[1] == "-version") {
		_, _ = fmt.Fprintf(os.Stdout, "protoc-gen-go-json version: %s\n", version())
		return
	}
	protoplugin.Main(protoplugin.HandlerFunc(plugin.Handle))
}

func version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok || info.Main.Version == "" {
		return "(devel)"
	}
	return info.Main.Version
}
