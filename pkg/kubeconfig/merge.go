package kubeconfig

// MergeConfig merges two kubeconfig [Config]s together using the same
// first-definition-wins and never-merge-values rules as kubectl.
//
// REF: https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/#merging-kubeconfig-files
func MergeConfig(first, other *Config) *Config {
	return &Config{
		ApiVersion:     firstNonNil(first.ApiVersion, other.ApiVersion),
		Kind:           firstNonNil(first.Kind, other.Kind),
		CurrentContext: firstNonNil(first.CurrentContext, other.CurrentContext),
		Preferences:    mergePreferences(first.Preferences, other.Preferences),
		Clusters:       mergeNamedSlice(first.Clusters, other.Clusters),
		Contexts:       mergeNamedSlice(first.Contexts, other.Contexts),
		AuthInfos:      mergeNamedSlice(first.AuthInfos, other.AuthInfos),
		Extensions:     mergeNamedSlice(first.Extensions, other.Extensions),
		Remaining:      mergeRemaining(first.Remaining, other.Remaining),
	}
}

func mergePreferences(first, other *Preferences) *Preferences {
	if other == nil {
		return first
	}

	if first == nil {
		return other
	}

	return &Preferences{
		Colors:     firstNonNil(first.Colors, other.Colors),
		Extensions: mergeNamedSlice(first.Extensions, other.Extensions),
		Remaining:  mergeRemaining(first.Remaining, other.Remaining),
	}
}

func mergeNamedSlice[T interface{ key() *string }](first, other []T) []T {
	merged := make([]T, 0, len(first))
	seenKeys := make(map[string]bool)

	for _, item := range first {
		key := item.key()
		if key != nil {
			merged = append(merged, item)
			seenKeys[*key] = true
		}
	}

	for _, item := range other {
		key := item.key()
		if key != nil && !seenKeys[*key] {
			merged = append(merged, item)
			seenKeys[*key] = true
		}
	}

	return merged
}

// mergeRemaining attempts to merge any remaining, unknown fields.
//
// This is a shallow merge and cannot account for fields that are a slice
// of named items, such as the [Config.Clusters] field and similar.
func mergeRemaining(first, other map[string]any) map[string]any {
	merged := make(map[string]any)

	for key, val := range other {
		merged[key] = val
	}

	for key, val := range first {
		merged[key] = val // duplicate keys are replaced by the first map
	}

	return merged
}
