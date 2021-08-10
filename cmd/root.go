package cmd

import (
	"fmt"
	"github.com/pathbird/pbauthor/cmd/auth"
	"github.com/pathbird/pbauthor/cmd/codex"
	"github.com/pathbird/pbauthor/internal/config"
	"github.com/pathbird/pbauthor/internal/version"
	"github.com/spf13/cobra"
	"os"

	log "github.com/sirupsen/logrus"
)

// flag vars
var (
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:           "pbauthor",
	SilenceErrors: true,
	SilenceUsage:  true,

	// Setup before running any other command
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := setUpLog(verbose); err != nil {
			panic(err)
		}
		version.CheckVersionAndPrintUpgradeNotice()
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		format := "error: %s\n"
		if verbose {
			format = "error: %#+v\n"
		}
		_, _ = fmt.Fprintf(os.Stderr, format, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(
		&verbose,
		"debug",
		false,
		"enable verbose logging (for debugging)",
	)
	rootCmd.PersistentFlags().StringVar(
		&config.PathbirdApiHost,
		"api-host",
		config.PathbirdApiHost,
		"Pathbird API host",
	)

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(auth.Cmd)
	rootCmd.AddCommand(codex.Cmd)
}

func setUpLog(verbose bool) error {
	if verbose {
		log.SetLevel(log.DebugLevel)
		log.Debugf("enabled debug logging")
	}
	return nil
}
