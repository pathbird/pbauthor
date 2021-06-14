package auth

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with the Pathbird API",
	Long:  "Authenticate with the Pathbird API.",
}
