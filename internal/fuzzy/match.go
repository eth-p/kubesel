package fuzzy

import (
	"errors"
	"slices"

	"github.com/junegunn/fzf/src/algo"
	"github.com/junegunn/fzf/src/util"
)

var ErrNoMatch = errors.New("no match")

// Matches runs a fuzzy find query against a list of strings, returning
// the strings which can be matched by the query.
func Matches(items []string, query string) []string {
	var result []string

	// Run every item through fzf's match algorithm with the query.
	queryRunes := []rune(query)
	slab := util.MakeSlab(1024, 8192)
	for _, item := range items {
		chars := util.ToChars([]byte(item))
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

		result = append(result, item)
	}

	return result
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
		items = Matches(items, query)
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
