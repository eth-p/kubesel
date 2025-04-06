package kubeconfig

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// UnmarshalJSON implements [json.Unmarshaler].
//
// This is needed because `inline` tag isn't supported by `encoding/json`
// ([issue]), and we want to at least support decoding extensions properly.
// When [encoding/json/v2] becomes available, that would be preferred over
// this hack.
//
// [encoding/json/v2]: https://github.com/golang/go/issues/71497
// [issue]: https://github.com/golang/go/issues/6213
func (e *Extension) UnmarshalJSON(data []byte) error {
	fields := make(map[string]any)

	err := json.Unmarshal(data, &fields)
	if err != nil {
		return err
	}

	if apiVersionField, ok := fields["apiVersion"]; ok {
		apiVersion, ok := apiVersionField.(string)
		if !ok {
			return &json.UnmarshalTypeError{
				Value:  fmt.Sprintf("%v", apiVersionField),
				Type:   reflect.TypeFor[string](),
				Struct: "Extension",
				Field:  "ApiVersion",
			}
		}

		e.ApiVersion = &apiVersion
		delete(fields, "apiVersion")
	}

	if kindField, ok := fields["kind"]; ok {
		kind, ok := kindField.(string)
		if !ok {
			return &json.UnmarshalTypeError{
				Value:  fmt.Sprintf("%v", kindField),
				Type:   reflect.TypeFor[string](),
				Struct: "Extension",
				Field:  "Kind",
			}
		}

		e.Kind = &kind
		delete(fields, "kind")
	}

	e.Remaining = fields

	return nil
}

// MarshalJSON implements [json.Marshaler].
//
// This is needed because `inline` tag isn't supported by `encoding/json`
// ([issue]), and we want to at least support decoding extensions properly.
// When [encoding/json/v2] becomes available, that would be preferred over
// this hack.
//
// [encoding/json/v2]: https://github.com/golang/go/issues/71497
// [issue]: https://github.com/golang/go/issues/6213
func (e *Extension) MarshalJSON() ([]byte, error) {
	const numFixedKeys = 2
	fields := make(map[string]any, len(e.Remaining)+numFixedKeys)

	for key, val := range e.Remaining {
		fields[key] = val
	}

	if e.ApiVersion != nil {
		fields["apiVersion"] = e.ApiVersion
	}

	if e.Kind != nil {
		fields["kind"] = e.Kind
	}

	return json.Marshal(fields)
}
