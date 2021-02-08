package auth

import (
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with the Mynerva API",
	Long:  "Authenticate with the Mynerva API.",
}

func InitAuthCmd(cmd *cobra.Command) {
	cmd.AddCommand(authCmd)

	authCmd.AddCommand(authLoginCmd, authStatusCmd)
}
