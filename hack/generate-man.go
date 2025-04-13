package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/eth-p/kubesel/internal/cli"
	"github.com/eth-p/kubesel/internal/cobraprint"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var (
	outputDir string
)

func main() {
	flag.StringVar(&outputDir, "outdir", "", "the output path")
	flag.Parse()
	if outputDir == "" {
		fmt.Fprintf(os.Stderr, "-outdir must be specified\n")
		os.Exit(1)
	}

	// Fix inconsistencies.
	cli.RootCommand.Use = "kubesel"
	fixCommandDescriptions(&cli.RootCommand)

	// Generate the manpages.
	err := doc.GenManTreeFromOpts(&cli.RootCommand, doc.GenManTreeOptions{
		Header:           &doc.GenManHeader{},
		Path:             outputDir,
		CommandSeparator: "-",
	})

	if err != nil {
		panic(err)
	}
}

func fixCommandDescriptions(cmd *cobra.Command) {
	cmd.Short = cobraprint.FixDescriptionWhitespace(cmd.Short)
	cmd.Long = cobraprint.FixDescriptionWhitespace(cmd.Long)
	cmd.Example = cobraprint.FixDescriptionWhitespace(cmd.Example)
	for _, subcmd := range cmd.Commands() {
		fixCommandDescriptions(subcmd)
	}
}
