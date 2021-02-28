package auth

import (
	"github.com/manifoldco/promptui"
	"github.com/mynerva-io/author-cli/internal/auth"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to Mynerva",

	RunE: func(cmd *cobra.Command, args []string) error {
		credentials, err := authLoginPrompt()
		if err != nil {
			return err
		}
		authResult, err := auth.AuthenticateWithPassword(credentials.email, credentials.password)
		if err != nil {
			return err
		}
		if authResult == nil {
			return errors.New("failed to authenticate")
		}
		return nil
	},
}

func init() {
	Cmd.AddCommand(authLoginCmd)
}

type authLoginPromptResult struct {
	email    string
	password string
}

func authLoginPrompt() (*authLoginPromptResult, error) {
	emailPrompt := promptui.Prompt{
		Label: "Email",
	}
	email, err := emailPrompt.Run()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to prompt email")
	}

	passwordPrompt := promptui.Prompt{
		Label: "Password",
		Mask:  '*',
	}
	password, err := passwordPrompt.Run()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to prompt password")
	}

	return &authLoginPromptResult{email: email, password: password}, nil
}
