package printer

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
	"github.com/mattn/go-runewidth"
)

var ansiResetCode = ansi.SGR(ansi.ResetAttr)

// MakePaddingWithChar returns a string that will pad out the specified string
// to be at least `width` characters wide when displayed on a terminal.
func MakePaddingWithChar(str string, width int, padChar string) string {
	strWidth := runewidth.StringWidth(str)
	if strWidth >= width {
		return ""
	}

	return strings.Repeat(padChar, width-strWidth)
}

// MakePadding returns a string that will pad out the specified string
// to be at least `width` characters wide when displayed on a terminal.
func MakePadding(str string, width int) string {
	return MakePaddingWithChar(str, width, " ")
}

// ApplyColor applies a color to the provided string.
//
// NOTE: This does not work with nested color-applied strings.
func ApplyColor(color, str string) string {
	if color == "" || str == "" {
		return str
	}

	var sb strings.Builder
	sb.Grow(len(str) + len(color) + len(ansiResetCode))
	sb.WriteString(color)
	sb.WriteString(str)
	sb.WriteString(ansiResetCode)
	return sb.String()
}
