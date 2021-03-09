package cmd

import (
	"fmt"
	"github.com/mynerva-io/author-cli/cmd/auth"
	"github.com/mynerva-io/author-cli/cmd/codex"
	"github.com/mynerva-io/author-cli/internal/config"
	"github.com/mynerva-io/author-cli/internal/version"
	"github.com/spf13/cobra"
	"os"

	log "github.com/sirupsen/logrus"
)

// flag vars
var (
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:           "mynerva-author",
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
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVar(&config.MynervaApiHost, "api-host", config.MynervaApiHost, "Mynerva API host")

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
