package main

import (
	"errors"
	"io"
	"strings"

	"github.com/charmbracelet/x/ansi"
	"github.com/eth-p/kubesel/internal/printer"
)

var (
	ErrUnknownFormat      = errors.New("unknown output format")
	ErrFormatNoOptions    = errors.New("format does not take options")
	ErrFormatNeedsOptions = errors.New("format requires options")
)

// UseListOutput updates the [OutputFormat] to display items in a list.
func UseListOutput(opts string, target *OutputFormat) error {
	if opts != "" {
		return ErrFormatNoOptions
	}

	*target = OutputFormat{
		name:       "list",
		newPrinter: printer.List,
	}

	return nil
}

// UseTableOutput updates the [OutputFormat] to display items in a table.
func UseTableOutput(opts string, target *OutputFormat) error {
	var (
		columns []string = nil
		wide    bool     = false
	)

	if opts == "*" {
		wide = true
	} else if len(opts) > 0 {
		columns = strings.Split(opts, ",")
		wide = true
	}

	*target = OutputFormat{
		name: "table",
		newPrinter: func(typ printer.ItemType, out io.Writer) (printer.Printer, error) {
			opts := printer.TableOptions{
				PickColumns: columns,
				ShowWide:    wide,
				SortRows:    true,

				ColumnSeparator:       " │ ",
				BorderLeft:            "│ ",
				BorderRight:           " │",
				BorderTopLeft:         "┌─",
				BorderTopRight:        "─┐",
				BorderTopSeparator:    "─┬─",
				BorderTopFill:         "─",
				BorderMidLeft:         "├─",
				BorderMidRight:        "─┤",
				BorderMidSeparator:    "─┼─",
				BorderMidFill:         "─",
				BorderBottomLeft:      "└─",
				BorderBottomRight:     "─┘",
				BorderBottomSeparator: "─┴─",
				BorderBottomFill:      "─",
			}

			if GlobalOptions.Color {
				opts.HeaderColor = ansi.SGR(ansi.BoldAttr)
			}

			return printer.Table(typ, out, opts)
		},
	}

	return nil
}

// UseColumnOutput updates the [OutputFormat] to display items in columns.
func UseColumnOutput(opts string, target *OutputFormat) error {
	var (
		columns []string = nil
		wide    bool     = false
	)

	if opts == "*" {
		wide = true
	} else if len(opts) > 0 {
		columns = strings.Split(opts, ",")
		wide = true
	}

	*target = OutputFormat{
		name: "column",
		newPrinter: func(typ printer.ItemType, out io.Writer) (printer.Printer, error) {
			opts := printer.TableOptions{
				ColumnSeparator: "  ",
				HeaderTransform: strings.ToUpper,

				PickColumns: columns,
				ShowWide:    wide,
			}

			if GlobalOptions.Color {
				opts.HeaderColor = ansi.SGR(ansi.BoldAttr)
			}

			return printer.Table(typ, out, opts)
		},
	}

	return nil
}

// OutputFormat is the flag used by `kubesel list`.
type OutputFormat struct {
	name       string
	newPrinter func(item printer.ItemType, out io.Writer) (printer.Printer, error)
}

func (f *OutputFormat) DefaultIfUnset() {
	if f.name == "" {
		UseTableOutput("", f)
	}
}

// String implements [pflag.Value].
func (f *OutputFormat) String() string {
	f.DefaultIfUnset()
	return f.name
}

// Set implements [pflag.Value].
func (f *OutputFormat) Set(v string) error {
	format, opts, _ := strings.Cut(v, "=")

	switch format {
	case "list":
		return UseListOutput(opts, f)

	case "table":
		return UseTableOutput(opts, f)

	case "columns", "column", "cols", "col":
		return UseColumnOutput(opts, f)

	default:
		return ErrUnknownFormat
	}
}

// Type implements [pflag.Value].
func (f *OutputFormat) Type() string {
	return "format"
}
