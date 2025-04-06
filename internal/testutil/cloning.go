package testutil

import (
	"fmt"
	"reflect"
	"sync/atomic"
)

type BasicTypeGenerator struct {
	nextString atomic.Int64
	nextBool   atomic.Int64
}

func (c *BasicTypeGenerator) NewString() string {
	id := c.nextString.Add(1)
	return fmt.Sprintf("string #%v", id)
}

func (c *BasicTypeGenerator) NewBool() bool {
	id := c.nextBool.Add(1)
	return id%2 == 0
}

// GenerateCloneTestdata uses reflection to recursively fill a struct with
// example data. This trades simplicity for safety, ensuring that newly-added
// fields are not accidentally omitted from tests.
//
// BUG: This cannot handle recursive types.
func GenerateCloneTestdata(btg *BasicTypeGenerator, ty reflect.Type) reflect.Value {
	switch ty.Kind() {
	case reflect.Slice:
		slice := reflect.MakeSlice(ty, 1, 1)
		slice.Index(0).Set(GenerateCloneTestdata(btg, ty.Elem()))
		return slice

	case reflect.Map:
		dict := reflect.MakeMap(ty)
		key := GenerateCloneTestdata(btg, ty.Key())
		val := GenerateCloneTestdata(btg, ty.Elem())
		dict.SetMapIndex(key, val)
		return dict

	case reflect.String:
		str := reflect.ValueOf(btg.NewString())
		return str

	case reflect.Bool:
		str := reflect.ValueOf(btg.NewBool())
		return str

	case reflect.Pointer:
		ptr := reflect.New(ty.Elem())
		ptr.Elem().Set(GenerateCloneTestdata(btg, ty.Elem()))
		return ptr

	case reflect.Struct:
		inst := reflect.New(ty).Elem()
		for i, end := 0, ty.NumField(); i < end; i++ {
			inst.Field(i).Set(GenerateCloneTestdata(btg, ty.Field(i).Type))
		}
		return inst

	case reflect.Interface:
		if reflect.TypeFor[any]().AssignableTo(ty) {
			return reflect.ValueOf(btg.NewString())
		}
	}

	panic(fmt.Errorf("Unsupported type for generating clone testdata: %v", ty))
}

// VerifyClone uses reflection to recursively compare that two structs
// are identical in all ways except pointer addresses.
//
// If the latter struct differs, this will return an error with the path.
//
// BUG: This cannot handle recursive types.
// BUG: This cannot handle pointer map keys.
func VerifyClone(expected, actual reflect.Value, path string) error {
	if expected.Type() != actual.Type() {
		return fmt.Errorf(
			"mismatched types, expected %v but got %v (at %s)",
			expected.Type(),
			actual.Type(),
			path,
		)
	}

	if expected.Kind() != actual.Kind() {
		return fmt.Errorf(
			"mismatched kinds, expected %v but got %v (at %s)",
			expected.Kind(),
			actual.Kind(),
			path,
		)
	}

	switch expected.Kind() {
	case reflect.Slice:
		return verifyClonedSlice(expected, actual, path)
	case reflect.Map:
		return verifyClonedMap(expected, actual, path)
	case reflect.String, reflect.Bool:
		return verifyClonedPrimitive(expected, actual, path)
	case reflect.Pointer:
		return verifyClonedPointer(expected, actual, path)
	case reflect.Struct:
		return verifyClonedStruct(expected, actual, path)
	case reflect.Interface:
		if reflect.TypeFor[any]().AssignableTo(expected.Type()) {
			return VerifyClone(expected.Elem(), actual.Elem(), path)
		}
	}

	panic(fmt.Errorf("Unsupported type for validating cloned testdata: %v (at %s)", expected.Type(), path))
}

// verifyClonedSlice compares the two provided slices.
// The lengths and contents must be equal.
func verifyClonedSlice(expected, actual reflect.Value, path string) error {
	if expected.Len() != actual.Len() {
		return fmt.Errorf(
			"mismatched slice lengths, expected %v but got %v (at %s)",
			expected.Len(),
			actual.Len(),
			path,
		)
	}

	for i, end := 0, expected.Len(); i < end; i++ {
		err := VerifyClone(expected.Index(i), actual.Index(i), fmt.Sprintf("%s[%v]", path, i))
		if err != nil {
			return err
		}
	}

	return nil
}

// verifyClonedMap compares the two provided maps.
// The lengths and contents must be equal.
//
// BUG: This cannot handle pointer map keys.
func verifyClonedMap(expected, actual reflect.Value, path string) error {
	if expected.Len() != actual.Len() {
		return fmt.Errorf(
			"mismatched map sizes, expected %v but got %v (at %s)",
			expected.Len(),
			actual.Len(),
			path,
		)
	}

	for _, key := range expected.MapKeys() {
		err := VerifyClone(expected.MapIndex(key), actual.MapIndex(key), fmt.Sprintf("%s[%q]", path, key.String()))
		if err != nil {
			return err
		}
	}

	return nil
}

// verifyClonedPrimitive compares the two provided basic-type values.
// They must be equal.
func verifyClonedPrimitive(expected, actual reflect.Value, path string) error {
	if !expected.Equal(actual) {
		return fmt.Errorf(
			"mismatched %ss, expected %v but got %v (at %s)",
			expected.Kind().String(),
			expected.Len(),
			actual.Len(),
			path,
		)
	}

	return nil
}

// verifyClonedPointer compares the two provided pointers.
// They must not be equal in identity (address), but they must be equal
// in value.
func verifyClonedPointer(expected, actual reflect.Value, path string) error {
	if expected.IsNil() && actual.IsNil() {
		return nil
	}

	if expected.Elem().IsValid() != actual.Elem().IsValid() {
		return fmt.Errorf(
			"mismatched pointers, expected nil to be %v but got %v (at %s)",
			!expected.Elem().IsValid(),
			!actual.Elem().IsValid(),
			path,
		)
	}

	if expected.Equal(actual) {
		return fmt.Errorf(
			"identical pointers, expected a deep copy (at %s)",
			path,
		)
	}

	return VerifyClone(expected.Elem(), actual.Elem(), path)
}

// verifyClonedStruct compares the two provided structs.
// Each field must be equal.
func verifyClonedStruct(expected, actual reflect.Value, path string) error {
	for i, end := 0, expected.NumField(); i < end; i++ {
		err := VerifyClone(expected.Field(i), actual.Field(i), fmt.Sprintf("%s.%v", path, expected.Type().Field(i).Name))
		if err != nil {
			return err
		}
	}

	return nil
}
