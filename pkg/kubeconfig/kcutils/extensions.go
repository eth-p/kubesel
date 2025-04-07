package kcutils

import (
	"fmt"

	"github.com/eth-p/kubesel/pkg/kubeconfig"
	"github.com/go-viper/mapstructure/v2"
)

type HasExtension interface {
	kubeconfig.Config | kubeconfig.Cluster |
		kubeconfig.AuthInfo | kubeconfig.Context |
		kubeconfig.Preferences
}

// ExtensionsFrom returns the [kubeconfig.NamedExtension] attached to the
// provided kubeconfig struct type.
func ExtensionsFrom[T HasExtension](extensible *T) []kubeconfig.NamedExtension {
	switch refined := any(*extensible).(type) {
	case kubeconfig.Config:
		return refined.Extensions
	case kubeconfig.Cluster:
		return refined.Extensions
	case kubeconfig.AuthInfo:
		return refined.Extensions
	case kubeconfig.Context:
		return refined.Extensions
	case kubeconfig.Preferences:
		return refined.Extensions
	}

	panic(fmt.Sprintf("unsupported type: %T", extensible))
}

// DecodeExtension decodes a [kubeconfig.Extension] into a struct, excluding
// the `ApiVersion` and `Kind` fields. This uses the `json` struct tag when
// decoding.
func DecodeExtension[T any](extension *kubeconfig.Extension, target *T) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook:           nil,
		ZeroFields:           true,
		Squash:               true,
		Result:               target,
		TagName:              "json",
		SquashTagOption:      "inline",
		IgnoreUntaggedFields: false,
	})

	if err != nil {
		return err
	}

	return decoder.Decode(extension.Remaining)
}

// EncodeExtension encodes a struct into a [kubeconfig.Extension], excluding
// the `ApiVersion` and `Kind` fields. This uses the `json` struct tag when
// decoding.
func EncodeExtension[T any](extension *T, target *kubeconfig.Extension) error {
	target.Remaining = make(map[string]any)

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook:           nil,
		ZeroFields:           true,
		Squash:               true,
		Result:               &target.Remaining,
		TagName:              "json",
		SquashTagOption:      "inline",
		IgnoreUntaggedFields: false,
	})

	if err != nil {
		return err
	}

	return decoder.Decode(extension)
}
