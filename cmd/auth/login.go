package auth

import (
	"github.com/mynerva-io/author-cli/internal/auth"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to Mynerva",

	RunE: func(cmd *cobra.Command, args []string) error {
		authResult, err := auth.AuthenticateFromUserInput()
		if err != nil {
			return err
		}
		if authResult == nil {
			return errors.New("failed to authenticate")
		}
		return nil
	},
}
