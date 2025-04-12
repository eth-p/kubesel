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
	createManagedPropertyCommands(&userCommand, managedProperty[userInfo]{
		PropertyNameSingular: "user",
		PropertyNamePlural:   "users",
		GetItemInfos:         userInfoIter,
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

type userInfo struct {
	Name         *string `yaml:"name" printer:"Name,order=0"`
	AuthProvider string  `yaml:"auth-provider" printer:"Auth Provider,order=1"`
}

func userInfoIter() (iter.Seq[userInfo], error) {
	kubesel, err := Kubesel()
	if err != nil {
		return nil, err
	}

	return func(yield func(userInfo) bool) {
		for _, kcNamedUser := range kubesel.GetMergedKubeconfig().AuthInfos {
			kcUser := kcNamedUser.User
			if kcUser == nil {
				kcUser = &kubeconfig.AuthInfo{}
			}

			item := userInfo{
				Name:         kcNamedUser.Name,
				AuthProvider: summarizeAuthProvider(kcUser),
			}

			if !yield(item) {
				return
			}
		}
	}, nil
}

func summarizeAuthProvider(kcUser *kubeconfig.AuthInfo) string {
	if kcUser == nil {
		return ""
	}

	// Auth provider:
	if kcUser.AuthProvider != nil && kcUser.AuthProvider.Name != nil {
		authType := *kcUser.AuthProvider.Name
		switch authType {

		// OIDC:
		case "oidc":
			issuer, ok := kcUser.AuthProvider.Config["idp-issuer-url"]
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
	if kcUser.Exec != nil && kcUser.Exec.Command != nil {
		command, _, _ := strings.Cut(*kcUser.Exec.Command, " ")
		return command
	}

	// Basic auth type:
	if kcUser.Username != nil {
		return "basic (" + *kcUser.Username + ")"
	}

	return ""
}
