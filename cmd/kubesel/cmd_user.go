package main

import (
	"fmt"
	"iter"
	"net/url"
	"slices"
	"strings"

	"github.com/eth-p/kubesel/pkg/kubeconfig"
	"github.com/spf13/cobra"
)

var UserCommand = cobra.Command{
	RunE: UserCommandMain,

	Aliases: []string{
		"users",
		"usr",
	},

	Use:     "user [name]",
	GroupID: "Kubeconfig",

	Short: "Change the current user",
	Long: `
	`,
	Example: `
	`,

	Annotations: map[string]string{
		TypeNameAnnotation:       "user",
		PluralTypeNameAnnotation: "users",
	},

	Args:              cobra.RangeArgs(0, 1),
	ValidArgsFunction: nil,
}

var UserCommandOptions struct {
}

func init() {
	Command.AddCommand(&UserCommand)
	CreateListerFor(&UserCommand, UserListItemIter)
}

func UserCommandMain(cmd *cobra.Command, args []string) error {
	ksel, err := Kubesel()
	if err != nil {
		return err
	}

	session, err := ksel.CurrentSession()
	if err != nil {
		return err
	}

	knownUsers := ksel.GetAuthInfoNames()
	desiredUser := ""

	// Select the user.
	if len(args) == 0 {
		// TODO: picker
	} else {
		desiredUser = args[0]
	}

	if !slices.Contains(knownUsers, desiredUser) {
		return fmt.Errorf("unknown user: %v", desiredUser)
	}

	session.SetAuthInfoName(desiredUser)
	err = session.Save()
	if err != nil {
		return fmt.Errorf("error saving session: %w", err)
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
