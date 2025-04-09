package main

import (
	"github.com/eth-p/kubesel/internal/cli"
	"github.com/spf13/cobra/doc"
)

func main() {
	err := doc.GenManTreeFromOpts(&cli.Command, doc.GenManTreeOptions{
		Header: &doc.GenManHeader{},
		Path:   ".",
	})

	if err != nil {
		panic(err)
	}
}
