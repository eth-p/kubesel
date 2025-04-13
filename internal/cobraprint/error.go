package cobraprint

import (
	"errors"
	"fmt"
	"io"

	"github.com/eth-p/kubesel/internal/cobraerr"
	tc "github.com/eth-p/kubesel/internal/textcomponent"
	"github.com/eth-p/kubesel/pkg/kubesel"
	"github.com/spf13/cobra"
)

type ErrorPrinterOptions struct {
	HelpPrinter       *HelpPrinter
	BlockquoteIndent  string
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
		opts:   &p.opts,
		cmd:    cmd,
		err:    err,
		output: &root,
	}

	print.appendCommandName() // `kubesel subcmd: `
	print.appendError()

	// Render the text components.
	renderer := tc.NewRenderer()
	renderer.Render(print.output)
	_, _ = io.WriteString(w, renderer.String())
}

// errorPrintContext contains the context of an error printer.
//
// The context may be derived to do things such as changing the
// error or wrapping the output in another text component.
type errorPrintContext struct {
	opts   *ErrorPrinterOptions
	cmd    *cobra.Command
	output *tc.Sequence
	err    error
}

// withIndent derives the errorPrintContext, creating a sub-context
// where the text is indented.
func (p errorPrintContext) withIndent() errorPrintContext {
	newCtx := p // shallow copy

	// Create a new Sequence component to store the to-be-indented components.
	// Use it as the child for a LinePrefix component, and add the LinePrefix
	// component to this helpPrintContext's childOutput.
	childOutput := &tc.Sequence{}
	p.output.Append(&tc.LinePrefix{
		Prefix: &tc.Text{Text: p.opts.Indent},
		Child:  childOutput,
	})

	newCtx.output = childOutput
	return newCtx
}

// withBlockquote derives the errorPrintContext, creating a sub-context
// where the text is indented.
func (p errorPrintContext) withBlockquote(color string, title string) errorPrintContext {
	newCtx := p // shallow copy

	if title != "" {
		p.output.Append(
			&tc.Text{
				Text:  title,
				Color: color,
			},
			tc.Newline,
		)
	}

	// Create a new Sequence component to store the to-be-indented components.
	// Use it as the child for a LinePrefix component, and add the LinePrefix
	// component to this helpPrintContext's childOutput.
	childOutput := &tc.Sequence{}
	p.output.Append(
		&tc.LinePrefix{
			Prefix: &tc.Text{
				Text:  p.opts.BlockquoteIndent,
				Color: color,
			},
			Child: &tc.Trim{
				Trailing: true,
				Child:    childOutput,
			},
		},
		tc.Newline,
	)

	newCtx.output = childOutput
	return newCtx
}

func (p errorPrintContext) appendCommandName() {
	cmdName := "kubectl"
	if p.cmd != nil {
		cmdName = p.cmd.CommandPath()
	}

	p.output.Append(&tc.Text{
		Color: p.opts.ErrorCommandColor,
		Text:  fmt.Sprintf("%s: ", cmdName),
	})
}

func (p errorPrintContext) appendError() {
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

	// Kubesel errors.
	if errors.Is(p.err, kubesel.ErrUnmanaged) {
		p.appendErrorText("kubesel is not initialized\n")
		p.withBlockquote(p.opts.ErrorTextColor, "").
			appendText("Use `kubesel init` to set up kubesel in the current shell.")
		return
	}

	// Unknown error.
	p.appendErrorText("unexpected error\n")
	p.withBlockquote(p.opts.ErrorTextColor, "").
		appendText(p.err.Error())
}

func (p errorPrintContext) appendInvalidFlagError(err *cobraerr.InvalidFlagError) {
	p.appendErrorText(fmt.Sprintf("%q is not a valid value for --%s\n", err.Value, err.Flag))
	if err.Cause != "" {
		p.withBlockquote(p.opts.ErrorTextColor, "").
			appendText(err.Cause)
	}

	if p.opts.HelpPrinter != nil {
		p.output.Append(
			tc.Newline,
			&tc.Text{
				Text: p.opts.HelpPrinter.PrintCommandFlags(p.cmd),
			},
		)
	}
}

func (p errorPrintContext) appendUnknownFlagError(err *cobraerr.UnknownFlagError) {
	if err.IsShorthandFlag {
		p.appendErrorText(fmt.Sprintf("-%s (in -%s) is not a valid flag\n", err.Flag, err.FlagSet))
	} else {
		p.appendErrorText(fmt.Sprintf("--%s is not a valid flag\n", err.Flag))
	}

	if p.opts.HelpPrinter != nil {
		p.output.Append(
			tc.Newline,
			&tc.Text{
				Text: p.opts.HelpPrinter.PrintCommandFlags(p.cmd),
			},
		)
	}
}

func (p errorPrintContext) appendUnknownCommandError(err *cobraerr.UnknownCommandError) {
	p.appendErrorText(fmt.Sprintf("'%s' is not a valid subcommand\n", err.Command))
	if p.cmd == nil || p.cmd.Name() != err.ParentCommand {
		return
	}

	// Append suggestions.
	suggestions := p.cmd.SuggestionsFor(err.Command)
	if len(suggestions) > 0 {
		p.withBlockquote(p.opts.ErrorTextColor, "Did you mean:").
			appendCommandSuggestions(suggestions)
	}
}

func (p errorPrintContext) appendCommandSuggestions(suggestions []string) {
	for _, suggestion := range suggestions {
		p.output.Append(
			&tc.Text{Text: suggestion},
			tc.Newline,
		)
	}
}

func (p errorPrintContext) appendText(text string) {
	p.output.Append(&tc.Text{Text: text})
}

func (p errorPrintContext) appendErrorText(text string) {
	p.output.Append(&tc.Text{
		Color: p.opts.ErrorTextColor,
		Text:  text,
	})
}
