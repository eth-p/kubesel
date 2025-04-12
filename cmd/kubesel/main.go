package main

import (
	"os"
	"runtime/debug"

	"github.com/eth-p/kubesel/internal/cli"
)

// THIS FILE IS ONLY A PROGRAM ENTRYPOINT.
//
// The kubesel command-line is implemented by the `internal/cli` package.
// We can't implement it here as part of `main` because we need to be able
// to access the cobra.Command structs for manpage generation.

func main() {
	cli.DetectTerminal()
	exitcode, _ := cli.Run(os.Args[1:])
	os.Exit(exitcode)
}

func init() {
	cli.RootCommand.Version = VERSION
}

// VERSION is defined by a go build flag.
// If empty, the build info provided by the Go compiler will be used instead.
var VERSION string

func GetVersion() string {
	if VERSION != "" {
		return VERSION
	}

	if buildinfo, ok := debug.ReadBuildInfo(); ok {
		VERSION = buildinfo.Main.Version
	} else {
		VERSION = "unknown"
	}

	return VERSION
}
