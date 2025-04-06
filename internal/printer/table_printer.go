package printer

import (
	"io"
	"reflect"
	"slices"
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// TableOptions change how the [Table] printer behaves.
type TableOptions struct {
	ColumnSeparator string
	PickColumns     []string

	SortRows bool
	ShowWide bool

	HeaderTransform func(s string) string
	HeaderColor     string

	BorderLeft            string
	BorderRight           string
	BorderTopFill         string
	BorderTopSeparator    string
	BorderTopLeft         string
	BorderTopRight        string
	BorderMidFill         string
	BorderMidSeparator    string
	BorderMidLeft         string
	BorderMidRight        string
	BorderBottomFill      string
	BorderBottomSeparator string
	BorderBottomLeft      string
	BorderBottomRight     string
}

// Table returns a [Printer] that buffers its contents and writes everything
// out as a table when closed.
func Table(items ItemType, w io.Writer, opts TableOptions) (Printer, error) {
	var err error
	fields := items.Fields

	// If PrintColumns is not nil, print only specific columns.
	if opts.PickColumns != nil {
		fields, err = fields.Pick(opts.PickColumns)
		if err != nil {
			return nil, err
		}
	}

	// If ShowWide is false, remove "wide" columns.
	if !opts.ShowWide {
		fields = fields.FilterOut(func(f *ItemStructField) bool {
			return f.OnlyWide
		})
	}

	// Calculate the minimum column widths.
	nCols := len(fields)
	columnWidths := make([]int, nCols)
	for i, field := range fields {
		columnWidths[i] = len(field.Name)
	}

	// Set default options.
	if opts.ColumnSeparator == "" {
		opts.ColumnSeparator = " "
	}

	// Create the printer.
	return &tablePrinter{
		options: opts,
		writer:  w,
		fields:  fields,

		columnCount:  nCols,
		columnWidths: columnWidths,
	}, nil
}

type tablePrinter struct {
	options  TableOptions
	itemType reflect.Type
	writer   io.Writer
	fields   []ItemStructField

	columnCount  int
	columnWidths []int
	cells        []formatted
	rowIndex     []int // index (in `cells`) of 1st column for row i
}

func (p *tablePrinter) Add(item any) {
	p.rowIndex = append(p.rowIndex, len(p.cells))
	for i, field := range p.fields {
		// Format the field.
		value := reflect.ValueOf(item).FieldByIndex(field.ReflectFieldIndex)
		formatted := field.ReflectFormatter(value)

		// Update column widths.
		width := len(formatted.value)
		if width > p.columnWidths[i] {
			p.columnWidths[i] = width
		}

		// Append it as a cell.
		p.cells = append(p.cells, formatted)
	}
}

func (p *tablePrinter) Close() {
	var (
		sb       strings.Builder
		nColumns = p.columnCount
		nRows    = len(p.cells) / nColumns
	)

	// Sort the table by the first column.
	slices.SortFunc(p.rowIndex, func(a, b int) int {
		return strings.Compare(p.cells[a].value, p.cells[b].value)
	})

	// Print the top border.
	p.printBorder(
		&sb,
		p.options.BorderTopLeft, p.options.BorderTopFill,
		p.options.BorderTopSeparator, p.options.BorderTopRight,
	)

	// Print the header.
	p.prepareRow(&sb)
	for col, field := range p.fields {
		headerCell := field.Name
		if p.options.HeaderTransform != nil {
			headerCell = p.options.HeaderTransform(headerCell)
		}
		p.appendCell(&sb, col, headerCell, len(headerCell), p.options.HeaderColor)
	}
	p.finishRow(&sb)
	p.flushBuffer(&sb)

	// Print the border between the header.
	p.printBorder(
		&sb,
		p.options.BorderMidLeft, p.options.BorderMidFill,
		p.options.BorderMidSeparator, p.options.BorderMidRight,
	)

	// Print the cells.
	for row := range nRows {
		i := p.rowIndex[row]

		// Squeeze empty lines.
		if nColumns == 1 && p.cells[i].value == "" {
			i++
			continue
		}

		// Print the row.
		p.prepareRow(&sb)
		for col := range nColumns {
			cell := p.cells[i]
			i++
			p.appendCell(&sb, col, cell.value, len(cell.value), "")
		}
		p.finishRow(&sb)
		p.flushBuffer(&sb)
	}

	// Print the bottom border.
	p.printBorder(
		&sb,
		p.options.BorderBottomLeft, p.options.BorderBottomFill,
		p.options.BorderBottomSeparator, p.options.BorderBottomRight,
	)
}

// printBorder prints a horizontal table border.
func (p *tablePrinter) printBorder(sb *strings.Builder, left, fill, sep, right string) {
	if left == "" && right == "" && fill == "" && sep == "" {
		return
	}

	sb.WriteString(left)
	for col, width := range p.columnWidths {
		if col > 0 {
			sb.WriteString(sep)
		}

		sb.WriteString(strings.Repeat(fill, width))
	}
	sb.WriteString(right)
	sb.WriteRune('\n')

	p.flushBuffer(sb)
}

// appendCell appends a cell to the row.
func (p *tablePrinter) appendCell(sb *strings.Builder, col int, text string, textWidth int, color string) {
	if col > 0 {
		sb.WriteString(p.options.ColumnSeparator)
	}

	if color != "" {
		sb.WriteString(color)
	}

	sb.WriteString(text)
	sb.WriteString(strings.Repeat(" ", p.columnWidths[col]-textWidth))

	if color != "" {
		sb.WriteString(ansi.SGR(ansi.ResetAttr))
	}
}

// prepareRow resets the buffer.
// This must be called at the start of every row.
func (p *tablePrinter) prepareRow(sb *strings.Builder) {
	sb.Reset()
	if p.options.BorderLeft != "" {
		sb.WriteString(p.options.BorderLeft)
	}
}

// finishRow flushes the buffer to the destination.
// This must be called after every row.
func (p *tablePrinter) finishRow(sb *strings.Builder) {
	if p.options.BorderRight != "" {
		sb.WriteString(p.options.BorderRight)
	}

	sb.WriteRune('\n')
}

func (p *tablePrinter) flushBuffer(sb *strings.Builder) {
	io.WriteString(p.writer, sb.String())
	sb.Reset()
}
