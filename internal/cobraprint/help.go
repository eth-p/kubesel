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
		opts: &p.opts,
		cmd:  cmd,
		text: &root,
	}

	// Append the documentation.
	print.printCmdDescription()
	if cmd.Runnable() || cmd.HasSubCommands() {
		print.printCmdUsage()
		print.printCommandAliases()
		print.printCommandExample()
		print.printCommandSubcommands()
		print.printCommandFlagsSection()
	}

	// Render the text components.
	renderer := tc.NewRenderer()
	renderer.Render(&root)
	return renderer.String()
}

// PrintCommandUsage prints the usage documentation for the specified command.
func (p *HelpPrinter) PrintCommandUsage(cmd *cobra.Command) string {
	var root tc.Sequence
	print := helpPrintContext{
		opts: &p.opts,
		cmd:  cmd,
		text: &root,
	}

	// Append the documentation.
	print.printCmdUsage()
	print.printCommandAliases()
	print.printCommandExample()
	print.printCommandSubcommands()
	print.printCommandFlagsSection()

	if cmd.HasAvailableSubCommands() {
		root.Append(&tc.Text{
			Text: fmt.Sprintf("\n\nUse \"%s [command] --help\" for more information about a command.", cmd.CommandPath()),
		})
	}

	// Render the text components.
	renderer := tc.NewRenderer()
	renderer.Render(&root)
	return renderer.String()
}

// helpPrintContext contains the context of a help printer.
type helpPrintContext struct {
	opts *HelpPrinterOptions
	cmd  *cobra.Command
	text *tc.Sequence
}

func (p *helpPrintContext) withIndent() *helpPrintContext {
	newCtx := *p

	output := &tc.Sequence{}
	p.text.Append(&tc.LinePrefix{
		Prefix: &tc.Text{Text: p.opts.Indent},
		Child:  output,
	})

	newCtx.text = output
	return &newCtx
}

func (p *helpPrintContext) printSectionHeading(heading string) {
	p.text.Append(
		tc.Newline,
		&tc.Text{
			Text:  heading,
			Color: p.opts.HeadingColor,
		},
		tc.Newline,
	)
}

func (p *helpPrintContext) printCmdDescription() {
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
	p.text.Append(
		&tc.Text{Text: description},
		tc.Newline,
	)
}

// printCmdUsage prints the `Usage:` section.
func (p *helpPrintContext) printCmdUsage() {
	cmd := p.cmd
	if !cmd.Runnable() && !cmd.HasAvailableSubCommands() {
		return
	}

	p.printSectionHeading("Usage:")

	if cmd.Runnable() {
		p.text.Append(&tc.Text{
			Text: fmt.Sprintf("%s%s\n", p.opts.Indent, cmd.UseLine()),
		})
	}

	if cmd.HasAvailableSubCommands() {
		p.text.Append(&tc.Text{
			Text: fmt.Sprintf("%s%s [command]\n", p.opts.Indent, cmd.CommandPath()),
		})
	}
}

// printCommandAliases prints the `Aliases:` section.
func (p *helpPrintContext) printCommandAliases() {
	cmd := p.cmd
	if len(cmd.Aliases) == 0 {
		return
	}

	p.printSectionHeading("Aliases:")
	p.text.Append(
		&tc.Text{
			Text: fmt.Sprintf("%s%s\n", p.opts.Indent, cmd.NameAndAliases()),
		},
	)
}

// printCommandExample prints the `Example:` section.
func (p *helpPrintContext) printCommandExample() {
	cmd := p.cmd
	if !cmd.HasExample() {
		return
	}

	// Clean up the example string and indent it.
	example := strings.Trim(dedent.Dedent(cmd.Example), "\n")
	example = strings.ReplaceAll(example, "\n", "\n"+p.opts.Indent)

	p.printSectionHeading("Examples:")
	p.text.Append(
		&tc.Text{
			Text: fmt.Sprintf("%s%s\n", p.opts.Indent, example),
		},
	)
}

// printCommandExample prints the command's subcommands.
func (p *helpPrintContext) printCommandSubcommands() {
	cmd := p.cmd
	if !cmd.HasAvailableSubCommands() {
		return
	}

	cmds := cmd.Commands()

	// When there are no command groups.
	if len(cmd.Groups()) == 0 {
		p.printSectionHeading("Available Commands:")
		for _, subcmd := range cmds {
			if subcmd.IsAvailableCommand() {
				p.printSubcommandLine(subcmd)
			}
		}

		return
	}

	// When there are command groups.
	for _, group := range cmd.Groups() {
		p.printSectionHeading(group.Title)
		for _, subcmd := range cmds {
			if subcmd.GroupID == group.ID && subcmd.IsAvailableCommand() {
				p.printSubcommandLine(subcmd)
			}
		}
	}

	if !cmd.AllChildCommandsHaveGroup() {
		p.printSectionHeading("Additional Commands:")
		for _, subcmd := range cmds {
			if subcmd.GroupID == "" && (subcmd.IsAvailableCommand()) {
				p.printSubcommandLine(subcmd)
			}
		}
	}
}

func (p *helpPrintContext) printSubcommandLine(subcmd *cobra.Command) {
	name := subcmd.Name()
	p.text.Append(
		&tc.Text{
			Text: p.opts.Indent,
		},
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

func (p *helpPrintContext) printCommandFlagsSection() {
	cmd := p.cmd
	if !cmd.HasAvailableLocalFlags() && !cmd.HasAvailableInheritedFlags() {
		return
	}

	p.printSectionHeading("Flags:")
	p.withIndent().printCommandFlags()
}

func (p *helpPrintContext) printCommandFlags() {
	flags := gatherFlagsInfo(p.cmd)
	out := p.text

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
