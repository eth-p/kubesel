package cli

import (
	"iter"
	"net/url"
	"strings"

	"github.com/eth-p/kubesel/pkg/kubeconfig"
	"github.com/eth-p/kubesel/pkg/kubesel"
	"github.com/spf13/cobra"
)

var userCommand = cobra.Command{
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
}

var UserCommandOptions struct {
}

func init() {
	RootCommand.AddCommand(&userCommand)
	createManagedPropertyCommands(&userCommand, managedProperty[UserListItem]{
		PropertyNameSingular: "user",
		PropertyNamePlural:   "users",
		ListGenerator:        UserListItemIter,
		GetItemNames:         userNames,
		Switch:               userSwitchImpl,
	})
}

func userSwitchImpl(ksel *kubesel.Kubesel, managedKc *kubesel.ManagedKubeconfig, target string) error {
	managedKc.SetAuthInfoName(target)
	return managedKc.Save()
}

func userNames() ([]string, error) {
	kubesel, err := Kubesel()
	if err != nil {
		return nil, err
	}

	return kubesel.GetAuthInfoNames(), nil
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
