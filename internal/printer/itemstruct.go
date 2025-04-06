package printer

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

// ItemType contains information on how a printer should emit data from a
// specific struct type.
type ItemType struct {
	Fields ItemStructFieldList
}

// ItemTypeOf analyzes the provided [reflect.Type], returning an [ItemType]
// for it.
func ItemTypeOf(typ reflect.Type) (*ItemType, error) {
	if typ.Kind() != reflect.Struct {
		return nil, &InvalidTypeError{
			Expected: "struct",
			Actual:   typ,
		}
	}

	fields, err := analyzeFields(typ)
	if err != nil {
		return nil, err
	}

	return &ItemType{
		Fields: fields,
	}, nil
}

// ItemStructField holds information about a struct field and how to present it
// when printing it.
type ItemStructField struct {
	ReflectFieldIndex []int
	ReflectFormatter  FormatterFunc

	Name     string
	Order    int
	OnlyWide bool
}

// analyzeFields returns an [ItemStructFieldList] created from the provided
// struct's fields. This uses the `printer` struct tag to find additional
// information about fields.
//
// Tag format:
//
//	`printer:"<name>"`
//	`printer:"<name>,<opt>,..."`
//	`printer:"<name>,<opt>=<val>,..."`
//
// Supported tag options:
func analyzeFields(typ reflect.Type) ([]ItemStructField, error) {
	result := make([]ItemStructField, 0, typ.NumField())
	for i := range typ.NumField() {
		field := typ.Field(i)

		info := ItemStructField{
			ReflectFieldIndex: field.Index,
			Order:             i,
			Name:              field.Name,
			ReflectFormatter:  formatterForType(field.Type),
		}

		tag, ok := field.Tag.Lookup("printer")
		if ok {
			applyPrinterFieldTag(&info, tag)
		}

		result = append(result, info)
	}

	// Sort the struct fields.
	slices.SortStableFunc(result, func(a, b ItemStructField) int {
		return a.Order - b.Order
	})

	return ItemStructFieldList(result), nil
}

func applyPrinterFieldTag(target *ItemStructField, tag string) {
	name, opts, hasOpts := strings.Cut(tag, ",")
	target.Name = name

	if !hasOpts {
		return
	}

	for _, opt := range strings.Split(opts, ",") {
		optName, optVal, _ := strings.Cut(opt, "=")
		err := applyPrinterFieldTagOption(target, optName, optVal)
		if err != nil {
			panic(fmt.Sprintf("invalid printer struct tag: %v", err))
		}
	}

	_ = opts
}

func applyPrinterFieldTagOption(target *ItemStructField, optName string, optVal string) error {
	switch optName {
	case "wide":
		target.OnlyWide = true

	case "order":
		order, err := strconv.ParseInt(optVal, 10, strconv.IntSize)
		if err != nil {
			return fmt.Errorf("order option must be number, got %q", optVal)
		}

		target.Order = int(order)
	}

	return nil
}

// ItemStructFieldList is a list of [ItemStructField] objects.
type ItemStructFieldList []ItemStructField

func normalizeName(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", ""))
}

// Pick creates a copy of this [ItemStructFieldList], retaining only the fields
// named by the provided string array. This array also indicates the order
// of returned fields.
func (f ItemStructFieldList) Pick(pick []string) (ItemStructFieldList, error) {
	byName := make(map[string]*ItemStructField, len(f))
	for _, field := range f {
		byName[normalizeName(field.Name)] = &field
	}

	result := make([]ItemStructField, 0, len(f))
	for _, name := range pick {
		field, ok := byName[normalizeName(strings.ToLower(name))]
		if !ok {
			return nil, newUnknownFieldError(f, name)
		}

		result = append(result, *field)
	}

	return result, nil
}

// FilterOut creates a copy of this [ItemStructFieldList], retaining only the
// fields which the predicate function returns `false` for.
// named by the provided string array. This array also indicates the order
// of returned fields.
func (f ItemStructFieldList) FilterOut(predicate func(f *ItemStructField) bool) ItemStructFieldList {
	result := make([]ItemStructField, 0, len(f))
	for _, field := range f {
		if !predicate(&field) {
			result = append(result, field)
		}
	}

	return result
}
