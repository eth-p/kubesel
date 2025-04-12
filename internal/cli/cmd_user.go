package cli

import (
	"fmt"
	"iter"
	"net/url"
	"slices"
	"strings"

	"github.com/eth-p/kubesel/internal/fuzzy"
	"github.com/eth-p/kubesel/pkg/kubeconfig"
	"github.com/spf13/cobra"
)

var userCommand = cobra.Command{
	RunE: userCommandMain,

	Aliases: []string{
		"users",
		"usr",
	},

	Use:     "user [name]",
	GroupID: "Kubeconfig",

	Short: "Change to a different user",
	Long: `
		Change to a different Kubernetes auth user in the current shell.

		When selecting a user, you can use its full name as it
		appears in 'kubesel list users' or a fuzzy match of its
		name. If no user is specified or if the specified name
		fuzzily matches multiple users, a fzf picker will be
		opened.
	`,
	Example: `
		kubesel cluster my.cluster.example  # full name
		kubesel cluster myclstr             # fuzzy match
		kubesel cluster                     # fzf picker
	`,

	Args:              cobra.RangeArgs(0, 1),
	ValidArgsFunction: nil,
}

var UserCommandOptions struct {
}

func init() {
	RootCommand.AddCommand(&userCommand)
	createManagedPropertyCommands(&userCommand, managedProperty[UserListItem]{
		PropertyNameSingular: "user",
		PropertyNamePlural:   "users",
		ListGenerator:        UserListItemIter,
	})
}

func userCommandMain(cmd *cobra.Command, args []string) error {
	ksel, err := Kubesel()
	if err != nil {
		return err
	}

	managedConfig, err := ksel.GetManagedKubeconfig()
	if err != nil {
		return err
	}

	query := ""
	if len(args) > 0 {
		query = args[0]
	}

	available := ksel.GetAuthInfoNames()
	desired, err := fuzzy.MatchOneOrPick(available, query)
	if err != nil {
		return err
	}

	// Safeguard.
	if !slices.Contains(available, desired) {
		return fmt.Errorf("unknown user: %v", desired)
	}

	// Apply the user.
	managedConfig.SetAuthInfoName(desired)
	err = managedConfig.Save()
	if err != nil {
		return fmt.Errorf("error updating kubeconfig: %w", err)
	}

	return nil
}

type UserListItem struct {
	Name         *string `yaml:"name" printer:"Name,order=0"`
	AuthProvider string  `yaml:"auth-provider" printer:"Auth Provider,order=1"`
}

func UserListItemIter() (iter.Seq[UserListItem], error) {
	kubesel, err := Kubesel()
	if err != nil {
		return nil, err
	}

	return func(yield func(UserListItem) bool) {
		for _, user := range kubesel.GetMergedKubeconfig().AuthInfos {
			userInfo := user.User
			if userInfo == nil {
				userInfo = &kubeconfig.AuthInfo{}
			}

			item := UserListItem{
				Name:         user.Name,
				AuthProvider: summarizeAuthProvider(userInfo),
			}

			if !yield(item) {
				return
			}
		}
	}, nil
}

func summarizeAuthProvider(userInfo *kubeconfig.AuthInfo) string {
	if userInfo == nil {
		return ""
	}

	// Auth provider:
	if userInfo.AuthProvider != nil && userInfo.AuthProvider.Name != nil {
		authType := *userInfo.AuthProvider.Name
		switch authType {

		// OIDC:
		case "oidc":
			issuer, ok := userInfo.AuthProvider.Config["idp-issuer-url"]
			if ok {
				issuerUrl, err := url.Parse(issuer)
				if err != nil {
					return authType + " (" + issuer + ")"
				}

				return authType + " (" + issuerUrl.Hostname() + ")"
			}

		// Unknown:
		default:
			return authType
		}
	}

	// Plug-in auth type:
	if userInfo.Exec != nil && userInfo.Exec.Command != nil {
		command, _, _ := strings.Cut(*userInfo.Exec.Command, " ")
		return command
	}

	// Basic auth type:
	if userInfo.Username != nil {
		return "basic (" + *userInfo.Username + ")"
	}

	return ""
}
