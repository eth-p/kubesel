package cobraprint

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/eth-p/kubesel/internal/printer"
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
	out  io.Writer
}

func NewHelpPrinter(out io.Writer, opts HelpPrinterOptions) *HelpPrinter {
	return &HelpPrinter{
		out:  out,
		opts: opts,
	}
}

// PrintCommandHelp prints help documentation for the specified command.
// This should be given to [cobra.Command]'s `SetHelpFunc`.
func (p *HelpPrinter) PrintCommandHelp(cmd *cobra.Command, args []string) {
	w := p.out

	usage := cmd.Long
	if usage == "" {
		usage = cmd.Short
	}

	if usage != "" {
		fmt.Fprintln(w, strings.Trim(dedent.Dedent(usage), "\n"))
		fmt.Fprintln(w)
	}

	if cmd.Runnable() || cmd.HasSubCommands() {
		p.printCmdUsage(cmd)
		p.printCommandAliases(cmd)
		p.printCommandExample(cmd)
		p.printCommandSubcommands(cmd)
		p.printCommandFlags(cmd)
	}
}

// PrintCommandUsage prints the usage documentation for the specified command.
// This should be given to [cobra.Command]'s `SetUsageFunc`.
func (p *HelpPrinter) PrintCommandUsage(cmd *cobra.Command) error {
	w := cmd.OutOrStdout()
	p.printCmdUsage(cmd)
	p.printCommandAliases(cmd)
	p.printCommandExample(cmd)
	p.printCommandSubcommands(cmd)
	p.printCommandFlags(cmd)

	if cmd.HasAvailableSubCommands() {
		fmt.Fprintf(w, "\n\nUse \"%s [command] --help\" for more information about a command.", cmd.CommandPath())
	}

	fmt.Fprintln(w)
	return nil
}

// printCmdUsage prints the `Usage:` section.
func (p *HelpPrinter) printCmdUsage(cmd *cobra.Command) {
	w := p.out

	fmt.Fprintf(w, "%s\n", printer.ApplyColor(p.opts.HeadingColor, "Usage:"))

	if cmd.Runnable() {
		fmt.Fprintf(w, "%s%s\n", p.opts.Indent, cmd.UseLine())
	}

	if cmd.HasAvailableSubCommands() {
		fmt.Fprintf(w, "%s%s [command]\n", p.opts.Indent, cmd.CommandPath())
	}
}

// printCommandAliases prints the `Aliases:` section.
func (p *HelpPrinter) printCommandAliases(cmd *cobra.Command) {
	if len(cmd.Aliases) == 0 {
		return
	}

	w := p.out
	fmt.Fprintf(w, "\n%s\n", printer.ApplyColor(p.opts.HeadingColor, "Aliases:"))
	fmt.Fprintf(w, "%s%s\n", p.opts.Indent, cmd.NameAndAliases())
}

// printCommandExample prints the `Example:` section.
func (p *HelpPrinter) printCommandExample(cmd *cobra.Command) {
	if !cmd.HasExample() {
		return
	}

	// Clean up the example string and indent it.
	example := strings.Trim(dedent.Dedent(cmd.Example), "\n")
	example = strings.ReplaceAll(example, "\n", "\n"+p.opts.Indent)

	w := p.out
	fmt.Fprintf(w, "\n%s\n", printer.ApplyColor(p.opts.HeadingColor, "Examples:"))
	fmt.Fprintf(w, "%s%s\n", p.opts.Indent, example)
}

// printCommandExample prints the command's subcommands.
func (p *HelpPrinter) printCommandSubcommands(cmd *cobra.Command) {
	if !cmd.HasAvailableSubCommands() {
		return
	}

	w := p.out
	cmds := cmd.Commands()

	// When there are no command groups.
	if len(cmd.Groups()) == 0 {
		fmt.Fprintf(w, "\n%s\n", printer.ApplyColor(p.opts.HeadingColor, "Available Commands:"))
		for _, subcmd := range cmds {
			if subcmd.IsAvailableCommand() {
				p.printSubcommandLine(subcmd)
			}
		}

		return
	}

	// When there are command groups.
	for _, group := range cmd.Groups() {
		fmt.Fprintf(w, "\n%s\n", printer.ApplyColor(p.opts.HeadingColor, group.Title))
		for _, subcmd := range cmds {
			if subcmd.GroupID == group.ID && subcmd.IsAvailableCommand() {
				p.printSubcommandLine(subcmd)
			}
		}
	}

	if !cmd.AllChildCommandsHaveGroup() {
		fmt.Fprintf(w, "\n%s\n", printer.ApplyColor(p.opts.HeadingColor, "Additional Commands:"))
		for _, subcmd := range cmds {
			if subcmd.GroupID == "" && (subcmd.IsAvailableCommand()) {
				p.printSubcommandLine(subcmd)
			}
		}
	}
}

func (p *HelpPrinter) printSubcommandLine(subcmd *cobra.Command) {
	w := p.out
	name := subcmd.Name()
	fmt.Fprintf(w,
		"%s%s%s %s\n",
		p.opts.Indent,
		printer.ApplyColor(p.opts.CommandNameColor, name),
		printer.MakePadding(name, subcmd.NamePadding()),
		printer.ApplyColor(p.opts.CommandShortDescriptionColor, subcmd.Short),
	)
}

func (p *HelpPrinter) printCommandFlags(cmd *cobra.Command) {
	if !cmd.HasAvailableLocalFlags() && !cmd.HasAvailableInheritedFlags() {
		return
	}

	w := p.out
	flags := gatherFlagsInfo(cmd)

	// Pre-calculate common strings.
	noFlagShorthandSpacing := strings.Repeat(" ", flags.MaxShorthandWidth+2)
	maxSecondColWidth := flags.MaxNameWidth + flags.MaxVarNameWidth + 1

	// Iterate the flags and print them.
	fmt.Fprintf(w, "\n%s\n", printer.ApplyColor(p.opts.HeadingColor, "Flags:"))
	for _, flag := range flags.Flags {
		fmt.Fprint(w, p.opts.Indent)

		// Flag shorthand: `-x, `
		if flag.Shorthand != "" {
			fmt.Fprintf(w, "%s%s, ",
				printer.MakePadding(flag.Shorthand, flags.MaxShorthandWidth),
				printer.ApplyColor(p.opts.FlagNameColor, flag.Shorthand),
			)
		} else if flags.MaxShorthandWidth > 0 {
			fmt.Fprint(w, noFlagShorthandSpacing)
		}

		// Flag name: `--flag`
		width := flag.NameWidth + 1 + flag.VarNameWidth
		fmt.Fprintf(w, "%s %s",
			printer.ApplyColor(p.opts.FlagNameColor, flag.Name),
			printer.ApplyColor(p.opts.ArgTypeColor, flag.VarName),
		)

		// Description.
		fmt.Fprint(w, strings.Repeat(" ", maxSecondColWidth-width))
		fmt.Fprintf(w, "   %s",
			printer.ApplyColor(p.opts.FlagDescriptionColor, flag.Usage),
		)

		// Default.
		if flag.Default != "" {
			fmt.Fprintf(w, " (default %s)", flag.Default)
		}

		fmt.Fprintln(w)
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
