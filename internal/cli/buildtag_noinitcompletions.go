//go:build no_init_completions

package cli

func init() {
	initScriptLoadsCompletions = false
}
