package main

import (
	"github.com/spf13/cobra"
)

const (
	// TypeNameAnnotation is a [cobra.Command] annotation for describing
	// what type of item a subcommand is working with.
	TypeNameAnnotation = "lister-item-singular"

	// ListerItemNameAnnotation is a [cobra.Command] annotation for describing
	// the plural name of the item type a subcommand is working with.
	PluralTypeNameAnnotation = "lister-item-plural"
)

func getCobraCommandAnnotation(cmd *cobra.Command, name string) (string, bool) {
	if cmd.Annotations == nil {
		return "", false
	}

	value, ok := cmd.Annotations[name]
	return value, ok
}
