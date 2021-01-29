package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var authCmd = &cobra.Command{
	Use:                        "auth",
	Short:                      "Authenticate with the Mynerva API",
	Long:                       "Authenticate with the Mynerva API.",

	Run: func(cmd *cobra.Command, args []string) {
		_, _ = fmt.Fprintf(os.Stderr, "Not implemented.")
	},
}