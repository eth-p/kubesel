package main

import (
	"os"
	"runtime/debug"
	"sync"

	"github.com/eth-p/kubesel/pkg/kubesel"
	"github.com/spf13/cobra"
)

// Command is the root `kubesel` command.
var Command = cobra.Command{
	Use: "kubesel",
}

var GlobalOptions struct {
	Color bool
}

// _VERSION is defined by a go build tag.
// If not empty, this will be returned as the version of kubesel.
var _VERSION string

var Kubesel = sync.OnceValues(kubesel.NewKubesel)

func init() {
	Command.AddGroup(&cobra.Group{
		ID:    "Info",
		Title: "Informational Commands:",
	})

	Command.AddGroup(&cobra.Group{
		ID:    "Kubeconfig",
		Title: "Kubeconfig Commands:",
	})

	Command.AddGroup(&cobra.Group{
		ID:    "Kubesel",
		Title: "Kubesel Commands:",
	})

	Command.PersistentFlags().BoolVar(
		&GlobalOptions.Color,
		"color",
		true, // TODO: auto
		"Print with colors",
	)

	Command.Version = _VERSION
	if Command.Version == "" {
		if buildinfo, ok := debug.ReadBuildInfo(); ok {
			Command.Version = buildinfo.Main.Version
		} else {
			Command.Version = "unknown"
		}
	}
}

func main() {
	cmd, err := Command.ExecuteC()
	if err == nil {
		return
	}

	_ = cmd
	os.Exit(1)
}
