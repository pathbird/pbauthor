package cmd

import (
	"fmt"
	"github.com/pathbird/pbauthor/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print version information",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("current version: %s\n", version.Version)
		return nil
	},
}
