package cobraprint

import (
	"fmt"
	"slices"
	"strings"

	"github.com/eth-p/kubesel/internal/printer"
	tc "github.com/eth-p/kubesel/internal/textcomponent"
	"github.com/lithammer/dedent"
	"github.com/mattn/go-runewidth"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type HelpPrinterOptions struct {
	Indent                       string
	HeadingColor                 string
	CommandNameColor             string
	CommandShortDescriptionColor string
	FlagNameColor                string
	FlagValueColor               string
	ArgTypeColor                 string
	FlagDescriptionColor         string
}

type HelpPrinter struct {
	opts HelpPrinterOptions
}

func NewHelpPrinter(opts HelpPrinterOptions) *HelpPrinter {
	return &HelpPrinter{
		opts: opts,
	}
}

// PrintCommandHelp prints help documentation for the specified command.
func (p *HelpPrinter) PrintCommandHelp(cmd *cobra.Command, args []string) string {
	var root tc.Sequence
	print := helpPrintContext{
		opts:   &p.opts,
		cmd:    cmd,
		output: &root,
	}

	// Append the documentation.
	print.printCmdDescription()
	if cmd.Runnable() || cmd.HasSubCommands() {
		print.appendUsageSection()
		print.appendAliasesSection()
		print.appendExampleSection()
		print.appendSubcommandSection()
		print.appendFlagsSection()
	}

	// Render the text components.
	renderer := tc.NewRenderer()
	renderer.Render(&tc.Trim{
		Leading: true,
		Child:   &root,
	})

	return strings.TrimRight(renderer.String(), " ")
}

// PrintCommandUsage prints the usage documentation for the specified command.
func (p *HelpPrinter) PrintCommandUsage(cmd *cobra.Command) string {
	var root tc.Sequence
	print := helpPrintContext{
		opts:   &p.opts,
		cmd:    cmd,
		output: &root,
	}

	// Append the documentation.
	print.appendUsageSection()
	print.appendAliasesSection()
	print.appendExampleSection()
	print.appendSubcommandSection()
	print.appendFlagsSection()

	if cmd.HasAvailableSubCommands() {
		root.Append(&tc.Text{
			Text: fmt.Sprintf("\n\nUse \"%s [command] --help\" for more information about a command.", cmd.CommandPath()),
		})
	}

	// Render the text components.
	renderer := tc.NewRenderer()
	renderer.Render(&tc.Trim{
		Leading: true,
		Child:   &root,
	})

	return strings.TrimRight(renderer.String(), " ")
}

// FlagsBlock returns the flags help section for a command.
func (p *HelpPrinter) PrintCommandFlags(cmd *cobra.Command) string {
	var root tc.Sequence
	print := helpPrintContext{
		opts:   &p.opts,
		cmd:    cmd,
		output: &root,
	}

	// Append the documentation.
	print.appendFlagsSection()

	// Render the text components.
	renderer := tc.NewRenderer()
	renderer.Render(&tc.Trim{
		Leading: true,
		Child:   &root,
	})

	return strings.TrimRight(renderer.String(), " ")
}

// helpPrintContext contains the context of a help printer.
//
// The context may be derived to do things such as changing the
// command or wrapping the output in another text component.
type helpPrintContext struct {
	opts   *HelpPrinterOptions
	cmd    *cobra.Command
	output *tc.Sequence
}

// withIndent derives the helpPrintContext, creating a sub-context
// where the text is indented.
func (p helpPrintContext) withIndent() helpPrintContext {
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

func (p helpPrintContext) printSectionHeading(heading string) {
	p.output.Append(
		tc.Newline,
		&tc.Text{
			Text:  heading,
			Color: p.opts.HeadingColor,
		},
		tc.Newline,
	)
}

func (p helpPrintContext) printCmdDescription() {
	cmd := p.cmd
	description := cmd.Long
	if description == "" {
		description = cmd.Short
	}

	if description == "" {
		return
	}

	// Clean up the description string and append it.
	description = strings.Trim(dedent.Dedent(description), "\n")
	p.output.Append(
		&tc.Text{Text: description},
		tc.Newline,
	)
}

// appendUsageSection adds the `Usage:` heading and section.
func (p helpPrintContext) appendUsageSection() {
	if !p.cmd.Runnable() && !p.cmd.HasAvailableSubCommands() {
		return
	}

	p.printSectionHeading("Usage:")
	p.withIndent().appendUsageSectionContents()
}

// appendUsageSectionContents adds the `Usage:` section contents.
func (p helpPrintContext) appendUsageSectionContents() {
	if p.cmd.Runnable() {
		p.output.Append(
			&tc.Text{Text: p.cmd.UseLine()},
			tc.Newline,
		)
	}

	if p.cmd.HasAvailableSubCommands() {
		p.output.Append(&tc.Text{
			Text: fmt.Sprintf("%s [command]\n", p.cmd.CommandPath()),
		})
	}
}

// appendAliasesSection adds the `Aliases:` heading and section.
func (p helpPrintContext) appendAliasesSection() {
	cmd := p.cmd
	if len(cmd.Aliases) == 0 {
		return
	}

	p.printSectionHeading("Aliases:")
	p.withIndent().appendAliasesSectionContents()
}

// appendAliasesSectionContents adds the `Aliases:` section contents.
func (p helpPrintContext) appendAliasesSectionContents() {
	p.output.Append(
		&tc.Text{Text: p.cmd.NameAndAliases()},
		tc.Newline,
	)
}

// appendExampleSection adds the `Example:` heading and section.
func (p helpPrintContext) appendExampleSection() {
	cmd := p.cmd
	if !cmd.HasExample() {
		return
	}

	p.printSectionHeading("Examples:")
	p.withIndent().appendExampleSectionContents()
}

// appendExampleSectionContents adds the `Example:` section contents.
func (p helpPrintContext) appendExampleSectionContents() {
	example := FixDescriptionWhitespace(p.cmd.Example)
	p.output.Append(
		&tc.Text{Text: example},
		tc.Newline,
	)
}

// appendSubcommandSection adds the command's subcommand groups.
func (p helpPrintContext) appendSubcommandSection() {
	if !p.cmd.HasAvailableSubCommands() {
		return
	}

	groups := p.cmd.Groups()

	// If there are no command groups, print everything at once.
	if len(groups) == 0 {
		p.printSectionHeading("Available Commands:")
		p.withIndent().appendSubcommandSectionGroupContents(nil)
		return
	}

	// Otherwise, print them on a group-by-group basis.
	for _, group := range groups {
		p.printSectionHeading(group.Title)
		p.withIndent().appendSubcommandSectionGroupContents(group)
	}

	if !p.cmd.AllChildCommandsHaveGroup() {
		p.printSectionHeading("Additional Commands:")
		p.withIndent().appendSubcommandSectionGroupContents(nil)
	}
}

// appendSubcommandSectionGroupContents adds the commands within a subcommand
// group.
func (p helpPrintContext) appendSubcommandSectionGroupContents(group *cobra.Group) {
	groupID := ""
	if group != nil {
		groupID = group.ID
	}

	for _, subcmd := range p.cmd.Commands() {
		if subcmd.GroupID == groupID && subcmd.IsAvailableCommand() {
			p.appendSubcommand(subcmd)
		}
	}
}

// appendSubcommandSectionGroupContents adds a line describing the subcommand.
func (p helpPrintContext) appendSubcommand(subcmd *cobra.Command) {
	name := subcmd.Name()
	p.output.Append(
		&tc.Text{
			Text:  name,
			Color: p.opts.CommandNameColor,
		},
		&tc.Text{
			Text: printer.MakePadding(name, subcmd.NamePadding()),
		},
		&tc.Text{
			Text: " ",
		},
		&tc.Text{
			Text:  subcmd.Short,
			Color: p.opts.CommandShortDescriptionColor,
		},
		tc.Newline,
	)
}

// appendFlagsSection adds the `Flags:` heading and section.
func (p helpPrintContext) appendFlagsSection() {
	cmd := p.cmd
	if !cmd.HasAvailableLocalFlags() && !cmd.HasAvailableInheritedFlags() {
		return
	}

	p.printSectionHeading("Flags:")
	p.withIndent().appendFlagsSectionContents()
}

// appendAliasesSectionContents adds the `Flags:` section contents.
func (p helpPrintContext) appendFlagsSectionContents() {
	flags := gatherFlagsInfo(p.cmd)
	out := p.output

	// Pre-calculate common strings.
	noFlagShorthandSpacing := strings.Repeat(" ", flags.MaxShorthandWidth+2)
	maxSecondColWidth := flags.MaxNameWidth + flags.MaxVarNameWidth + 1

	// Iterate the flags and print them.
	for _, flag := range flags.Flags {

		// Flag shorthand: `-x, `
		if flag.Shorthand != "" {
			out.Append(
				&tc.Text{
					Text: printer.MakePadding(flag.Shorthand, flags.MaxShorthandWidth),
				},
				&tc.Text{
					Text:  flag.Shorthand,
					Color: p.opts.FlagNameColor,
				},
				&tc.Text{
					Text: ", ",
				},
			)
		} else if flags.MaxShorthandWidth > 0 {
			out.Append(&tc.Text{
				Text: noFlagShorthandSpacing,
			})
		}

		// Flag name: `--flag`
		width := flag.NameWidth + 1 + flag.VarNameWidth
		out.Append(
			&tc.Text{
				Text:  flag.Name,
				Color: p.opts.FlagNameColor,
			},
			&tc.Text{
				Text: " ",
			},
			&tc.Text{
				Text:  flag.VarName,
				Color: p.opts.FlagNameColor,
			},
		)

		// Description.
		out.Append(
			&tc.Text{
				Text:  strings.Repeat(" ", maxSecondColWidth-width),
				Color: p.opts.FlagNameColor,
			},
			&tc.Text{
				Text: "   ",
			},
			&tc.Text{
				Text:  flag.Usage,
				Color: p.opts.FlagDescriptionColor,
			},
		)

		// Default.
		if flag.Default != "" {
			out.Append(
				&tc.Text{
					Text: fmt.Sprintf(" (default %s)", flag.Default),
				},
			)
		}

		out.Append(tc.Newline)
	}
}

func gatherFlagsInfo(cmd *cobra.Command) flagsInfo {
	var flags flagsInfo

	visitFlag := func(f *pflag.Flag) {
		if f.Hidden {
			return
		}

		name := "--" + f.Name
		shorthand := ""
		if f.Shorthand != "" {
			shorthand = "-" + f.Shorthand
		}

		varname, usage := pflag.UnquoteUsage(f)
		flagInfo := flagInfo{
			Flag:           f,
			Name:           name,
			NameWidth:      runewidth.StringWidth(name),
			VarName:        varname,
			VarNameWidth:   runewidth.StringWidth(varname),
			Usage:          usage,
			UsageWidth:     runewidth.StringWidth(usage),
			Shorthand:      shorthand,
			ShorthandWidth: runewidth.StringWidth(shorthand),
		}

		if f.Name != "help" {
			flagInfo.Default = f.DefValue
			flagInfo.DefaultWidth = runewidth.StringWidth(f.DefValue)
		}

		flags.Flags = append(flags.Flags, flagInfo)
		flags.MaxNameWidth = maxInt(flags.MaxNameWidth, flagInfo.NameWidth)
		flags.MaxShorthandWidth = maxInt(flags.MaxShorthandWidth, flagInfo.ShorthandWidth)
		flags.MaxVarNameWidth = maxInt(flags.MaxVarNameWidth, flagInfo.VarNameWidth)
		flags.MaxUsageWidth = maxInt(flags.MaxUsageWidth, flagInfo.UsageWidth)
	}

	cmd.LocalFlags().VisitAll(visitFlag)
	cmd.InheritedFlags().VisitAll(visitFlag)

	// Sort the flags alphabetically.
	slices.SortFunc(flags.Flags, func(a, b flagInfo) int {
		return strings.Compare(a.Name, b.Name)
	})

	return flags
}

func maxInt(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

type flagsInfo struct {
	Flags                []flagInfo
	MaxNameWidth         int
	MaxVarNameWidth      int
	MaxUsageWidth        int
	MaxShorthandWidth    int
	MaxDefaultValueWidth int
}

type flagInfo struct {
	Flag              *pflag.Flag
	Name              string
	NameWidth         int
	VarName           string
	VarNameWidth      int
	Usage             string
	UsageWidth        int
	Shorthand         string
	ShorthandWidth    int
	NoOptDefault      string
	NoOptDefaultWidth int
	Default           string
	DefaultWidth      int
}

// FixDescriptionWhitespace removes indentation and leading/trailing newlines
// from a string. This cleans up the leftover whitespace caused by using
// a backtick string in the [cobra.Command] fields.
func FixDescriptionWhitespace(s string) string {
	s = dedent.Dedent(s)
	s = strings.Trim(s, "\n")
	return s
}
