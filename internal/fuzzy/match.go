package fuzzy

import (
	"errors"
	"slices"

	"github.com/junegunn/fzf/src/algo"
	"github.com/junegunn/fzf/src/util"
)

var ErrNoMatch = errors.New("no match")

type MatchResult[T any] struct {
	Item  T
	Score int
}

// MatchesFunc runs a fuzzy find query against a list of items, returning
// the [MatchResults] for all items matched by the query.
func MatchesFunc[T any](items []T, query string, getName func(T) string) []MatchResult[T] {
	var result []MatchResult[T]

	// Run every item through fzf's match algorithm with the query.
	queryRunes := []rune(query)
	slab := util.MakeSlab(1024, 8192)
	for _, item := range items {
		chars := util.ToChars([]byte(getName(item)))
		match, _ := algo.FuzzyMatchV2(
			false,      // not case sensitive
			true,       // normalized
			false,      // ??
			&chars,     // what to match against
			queryRunes, // the query
			false,      // do not return positions
			slab,       // slab allocator
		)

		if match.Score < 1 {
			continue
		}

		result = append(result, MatchResult[T]{
			Item:  item,
			Score: match.Score,
		})
	}

	return result
}

// Matches runs a fuzzy find query against a list of strings, returning
// the [MatchResult]s for all strings matched by the query.
func Matches(items []string, query string) []MatchResult[string] {
	return MatchesFunc(items, query, func(s string) string { return s })
}

// Matches runs a fuzzy find query against a list of strings, returning
// the strings which can be matched by the query.
func StringMatches(items []string, query string) []string {
	matches := Matches(items, query)
	if len(matches) == 0 {
		return []string{}
	}

	results := make([]string, len(matches))
	for i, match := range matches {
		results[i] = match.Item
	}

	return results
}

func MatchOneOrPick(items []string, query string) (string, error) {
	slices.Sort(items)

	if query != "" {
		// Exact match?
		_, found := slices.BinarySearch(items, query)
		if found {
			return query, nil
		}

		// Fuzzy match?
		items = StringMatches(items, query)
		switch len(items) {
		case 0:
			return "", ErrNoMatch

		case 1:
			return items[0], nil
		}
	}

	// No query or more than 1 item. Use fzf TUI as a picker.
	pickOpts := PickOptions{
		Query: query,
	}

	return Pick(items, &pickOpts)
}

// SortedMatchesFunc runs a fuzzy find query against a slice, returning the
// items which can be matched by the query. The returned slice is sorted by
// score.
func SortedMatchesFunc[T any](
	items []T,
	query string,
	getName func(*T) string,
	compare func(a, b MatchResult[T]) int,
) []T {
	var matches []MatchResult[T]

	// Run every item through fzf's match algorithm with the query.
	queryRunes := []rune(query)
	slab := util.MakeSlab(1024, 8192)
	for _, item := range items {
		chars := util.ToChars([]byte(getName(&item)))
		match, _ := algo.FuzzyMatchV2(
			false,      // not case sensitive
			true,       // normalized
			false,      // ??
			&chars,     // what to match against
			queryRunes, // the query
			false,      // do not return positions
			slab,       // slab allocator
		)

		if match.Score < 1 {
			continue
		}

		matches = append(matches, MatchResult[T]{
			Item:  item,
			Score: match.Score,
		})
	}

	// Sort the matches.
	slices.SortFunc(matches, compare)

	// Return the matched items.
	result := make([]T, len(matches))
	for i, match := range matches {
		result[i] = match.Item
	}

	return result
}
