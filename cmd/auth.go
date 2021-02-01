package cmd

import (
	"errors"
	"fmt"
	"github.com/mynerva-io/author-cli/internal/auth"
	"github.com/spf13/cobra"
	"time"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with the Mynerva API",
	Long:  "Authenticate with the Mynerva API.",
}

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

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check authentication status",

	RunE: func(cmd *cobra.Command, args []string) error {
		auth, err := auth.GetAuth()
		if err != nil {
			return err
		}
		if auth == nil {
			fmt.Printf("Not authenticated.")
			return nil
		}
		fmt.Printf("Authenticated (until %s\n", auth.Expiration.Format(time.RFC1123))
		return nil
	},
}

func initAuth(cmd *cobra.Command) {
	cmd.AddCommand(authCmd)

	authCmd.AddCommand(authLoginCmd, authStatusCmd)
}
