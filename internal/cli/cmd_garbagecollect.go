package cli

import (
	"fmt"

	"github.com/eth-p/kubesel/pkg/kubesel"
	"github.com/spf13/cobra"
)

var gcCommand = cobra.Command{
	Use:     "garbage-collect",
	GroupID: "Kubesel",
	// Hidden:  true,

	Short: "Remove kubesel-managed files for defunct shells",
	Long: `
		Search for and remove any files created by kubesel which
		belong to processes that are no longer alive.
	`,
	Example: `
		kubesel garbage-collect
	`,

	RunE: gcCommandMain,
}

var GCCommandOptions struct {
}

func init() {
	RootCommand.AddCommand(&gcCommand)
}

func gcCommandMain(cmd *cobra.Command, args []string) error {
	ksel, err := Kubesel()
	if err != nil {
		return err
	}

	// Run the GC function until everything is checked.
	results, err := ksel.GarbageCollect(&kubesel.GarbageCollectOptions{})
	if err != nil {
		return fmt.Errorf("garbage collection failed: %w", err)
	}

	// Print the results.
	w := cmd.OutOrStdout()
	fmt.Fprintf(w, "Files checked: %d\n", len(results.FilesChecked))
	fmt.Fprintf(w, "Files deleted: %d\n", len(results.FilesDeleted))
	fmt.Fprintf(w, "Errors: %d\n", len(results.Errors))
	for _, err := range results.Errors {
		fmt.Fprintf(w, " - %v\n", err)
	}

	return nil
}
