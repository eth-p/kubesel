package printer

import (
	"fmt"
	"reflect"
)

type formattedFlag int
type FormatterFunc func(reflect.Value) formatted

const (
	ffOK      formattedFlag = 0
	ffMissing formattedFlag = 1 << iota
)

type formatted struct {
	value string
	flag  formattedFlag
}

func formatterForType(typ reflect.Type) FormatterFunc {
	ptrToType := reflect.Zero(typ)
	switch ptrToType.Interface().(type) {
	case *string:
		return formatStringPointer

	case string:
		return formatString

	default:
		panic(fmt.Sprintf("no formatter for %v", typ.String()))
	}
}

func formatStringPointer(value reflect.Value) formatted {
	if value.IsNil() {
		return formatted{
			value: "",
			flag:  ffMissing,
		}
	}

	return formatted{
		value: *(value.Interface().(*string)),
	}
}

func formatString(value reflect.Value) formatted {
	return formatted{
		value: value.Interface().(string),
	}
}
