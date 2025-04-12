package cli

import (
	"iter"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
)

// ListItemGenerator creates an iterator that returns structs of items to be
// displayed by the `kubesel list` subcommand.
type ListItemGenerator[I any] func() (iter.Seq[I], error)

type lister struct {
	identifier       string
	typeNameSingular string
	typeNamePlural   string
	itemType         reflect.Type
	itemGenerator    ListItemGenerator[any]
	cmd              *cobra.Command
}

// CreateListerFor adds a new type to the `kubsel list` command.
//
// This uses the provided [cobra.Command]'s name and aliases as identifiers
// for the type, generating the CLI such that `kubesel clusters` would add
// `kubesel list {clusters|cluster|cl}` as valid options.
//
// This also adds a `--list` flag to the provided [cobra.Command], which
// redirects execution from `kubesel <type>` to `kubesel list <type>` if
// set by the user.
func CreateListerFor[I any](cmd *cobra.Command, generator ListItemGenerator[I]) {
	cmdName := cmd.Name()

	typeName := "item"
	if v, ok := getCobraCommandAnnotation(cmd, TypeNameAnnotation); ok {
		typeName = v
	}

	typeNamePlural := "items"
	if v, ok := getCobraCommandAnnotation(cmd, PluralTypeNameAnnotation); ok {
		typeNamePlural = v
	}

	// Create the lister.
	newLister := &lister{
		identifier:       cmdName,
		typeNameSingular: typeName,
		typeNamePlural:   typeNamePlural,
		itemGenerator:    upcastListerItemGenerator(generator),
		itemType:         reflect.TypeFor[I](),
		cmd: &cobra.Command{
			Use:     cmdName,
			Aliases: cmd.Aliases,

			Short: "List available " + typeNamePlural,
			Long:  strings.ReplaceAll(listCommand.Long, "info", typeNamePlural),

			Args: cobra.NoArgs,
		},
	}

	newLister.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return listCommandImpl(newLister, cmd, args)
	}

	// Update commands.
	listCommand.AddCommand(newLister.cmd)
	injectListFlag(cmd, newLister)
}

func injectListFlag(target *cobra.Command, newLister *lister) {
	const flagName = "list"

	// Add the `--list` flag and set it as hidden.
	var printList bool
	target.Flags().BoolVar(&printList, flagName, false, "")
	flag := target.Flag(flagName)
	flag.Hidden = true
	flag.Usage = "Print the available " + newLister.typeNamePlural

	// Wrap the target command's `RunE` function.
	realRunE := target.RunE
	target.RunE = func(cmd *cobra.Command, args []string) error {
		if printList {
			UseListOutput("", &ListCommandOptions.OutputFormat)
			return newLister.cmd.RunE(newLister.cmd, []string{})
		}

		return realRunE(cmd, args)
	}
}

// upcastListerItemGenerator converts an I-returning [ListItemGenerator] into an
// any-returning [ListItemGenerator].
func upcastListerItemGenerator[I any](generator ListItemGenerator[I]) ListItemGenerator[any] {
	return func() (iter.Seq[any], error) {
		realGenerator, err := generator()
		if err != nil {
			return nil, err
		}

		return func(yield func(any) bool) {
			for value := range realGenerator {
				if !yield(value) {
					break
				}
			}
		}, nil
	}
}
