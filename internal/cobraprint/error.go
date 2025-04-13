package cobraprint

import (
	"fmt"
	"io"

	"github.com/eth-p/kubesel/internal/cobraerr"
	tc "github.com/eth-p/kubesel/internal/textcomponent"
	"github.com/spf13/cobra"
)

type ErrorPrinterOptions struct {
	HelpPrinter       *HelpPrinter
	Indent            string
	ErrorCommandColor string
	ErrorTextColor    string
	TipColor          string
}

// ErrorPrinter is a utility for pretty-printing the error returned by
// [cobra.Command.Execute].
type ErrorPrinter struct {
	opts ErrorPrinterOptions
}

func NewErrorPrinter(opts ErrorPrinterOptions) *ErrorPrinter {
	return &ErrorPrinter{
		opts: opts,
	}
}

func (p *ErrorPrinter) PrintCommandError(w io.Writer, cmd *cobra.Command, err error) {
	err = cobraerr.Parse(err) // try to parse cobra's unstructured errors

	var root tc.Sequence
	print := errorPrintContext{
		opts: &p.opts,
		cmd:  cmd,
		err:  err,
		text: &root,
	}

	print.prependCommandName() // `kubesel subcmd: `
	print.appendError()

	// Render the text components.
	renderer := tc.NewRenderer()
	renderer.Render(print.text)
	_, _ = io.WriteString(w, renderer.String())
}

// errorPrintContext is a scope of the error printer.
// This is used to isolate nested errors from the outermost errors.
type errorPrintContext struct {
	opts *ErrorPrinterOptions
	cmd  *cobra.Command
	err  error
	text *tc.Sequence
}

func (p *errorPrintContext) prependCommandName() {
	cmdName := "kubectl"
	if p.cmd != nil {
		cmdName = p.cmd.CommandPath()
	}

	p.text.Append(&tc.Text{
		Color: p.opts.ErrorCommandColor,
		Text:  fmt.Sprintf("%s: ", cmdName),
	})
}

func (p *errorPrintContext) appendError() {
	// Handle invalid flag, unknown flag, and unknown command errors specially.
	switch refined := p.err.(type) {
	case *cobraerr.InvalidFlagError:
		p.appendInvalidFlagError(refined)
		return

	case *cobraerr.UnknownFlagError:
		p.appendUnknownFlagError(refined)
		return

	case *cobraerr.UnknownCommandError:
		p.appendUnknownCommandError(refined)
		return
	}

	// Unknown
	p.text.Append(
		&tc.Text{
			Color: "",
			Text:  p.err.Error(),
		},
		tc.Newline,
	)
}

func (p *errorPrintContext) appendInvalidFlagError(err *cobraerr.InvalidFlagError) {
	p.text.Append(&tc.Text{
		Color: p.opts.ErrorTextColor,
		Text:  fmt.Sprintf("%q is not a valid value for --%s\n", err.Value, err.Flag),
	})

	if err.Cause != "" {
		p.text.Append(
			&tc.Text{
				Text: err.Cause,
			},
			tc.Newline,
		)
	}

	// TODO: Show flags
}

func (p *errorPrintContext) appendUnknownFlagError(err *cobraerr.UnknownFlagError) {
	if err.IsShorthandFlag {
		p.text.Append(&tc.Text{
			Color: p.opts.ErrorTextColor,
			Text:  fmt.Sprintf("-%s (in -%s) is not a valid flag\n", err.Flag, err.FlagSet),
		})
	} else {
		p.text.Append(&tc.Text{
			Color: p.opts.ErrorTextColor,
			Text:  fmt.Sprintf("--%s is not a valid flag\n", err.Flag),
		})
	}

	// TODO: Show flags
}

func (p *errorPrintContext) appendUnknownCommandError(err *cobraerr.UnknownCommandError) {
	p.text.Append(&tc.Text{
		Color: p.opts.ErrorTextColor,
		Text:  fmt.Sprintf("'%s' is not a valid subcommand\n", err.Command),
	})

	if p.cmd == nil || p.cmd.Name() != err.ParentCommand {
		return
	}

	// Get suggestions.
	suggestions := p.cmd.SuggestionsFor(err.Command)
	if len(suggestions) < 1 {
		return
	}

	// Create a suggestions box.
	contents := &tc.Sequence{}
	for _, suggestion := range suggestions {
		contents.Append(&tc.Text{
			Text: fmt.Sprintf("%s\n", suggestion),
		})
	}

	// Append the suggestions box.
	p.text.Append(
		tc.Newline,
		&tc.Text{
			Color: p.opts.TipColor,
			Text:  "Did you mean: \n",
		},
		&tc.LinePrefix{
			Child: &tc.Trim{
				Trailing: true,
				Child:    contents,
			},
			Prefix: &tc.Text{
				Color: p.opts.TipColor,
				Text:  "▏",
			},
		},
		tc.Newline,
	)
}
