package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/eth-p/kubesel/internal/cli"
	"github.com/shirou/gopsutil/v4/host"
)

// THIS FILE IS ONLY A PROGRAM ENTRYPOINT.
//
// The kubesel command-line is implemented by the `internal/cli` package.
// We can't implement it here as part of `main` because we need to be able
// to access the cobra.Command structs for manpage generation.

func main() {
	host.EnableBootTimeCache(true)
	cli.DetectTerminal()
	exitcode, _ := cli.Run(os.Args[1:])
	os.Exit(exitcode)
}

func init() {
	cli.RootCommand.Version = getVersion()
}

// VERSION is the most recent git tag.
var VERSION string = "v0.0.1"

// GIT_REVISION is defined by a go build flag.
// If empty, this will be read from the build info instead.
var GIT_REVISION string

func getVersion() string {
	if GIT_REVISION == "" {
		if buildinfo, ok := debug.ReadBuildInfo(); ok {
			for _, buildsetting := range buildinfo.Settings {
				switch buildsetting.Key {
				case "vcs.revision":
					GIT_REVISION = buildsetting.Value
				}
			}
		}
	}

	return fmt.Sprintf("%s (%s)", VERSION, GIT_REVISION)
}
