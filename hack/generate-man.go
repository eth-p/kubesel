package main

import (
	"github.com/eth-p/kubesel/internal/cli"
	"github.com/eth-p/kubesel/internal/cobraprint"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func main() {
	cli.RootCommand.Use = "kubesel"
	fixCommandDescriptions(&cli.RootCommand)

	// Generate the manpages.
	err := doc.GenManTreeFromOpts(&cli.RootCommand, doc.GenManTreeOptions{
		Header:           &doc.GenManHeader{},
		Path:             ".",
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
