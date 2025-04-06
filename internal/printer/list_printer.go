package printer

import (
	"fmt"
	"io"
	"reflect"
)

// List returns a [Printer] that writes a list of items.
func List(items ItemType, w io.Writer) (Printer, error) {
	return &listPrinter{
		field:  items.Fields[0],
		writer: w,
	}, nil
}

type listPrinter struct {
	itemType reflect.Type
	field    ItemStructField
	writer   io.Writer
}

func (p *listPrinter) Add(item any) {
	value := reflect.ValueOf(item).FieldByIndex(p.field.ReflectFieldIndex)
	fmt.Fprintln(p.writer, p.field.ReflectFormatter(value).value)
}

func (p *listPrinter) Close() {
}
