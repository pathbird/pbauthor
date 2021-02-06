package cmd

import (
	"fmt"
	"github.com/mynerva-io/author-cli/internal/config"
	"github.com/spf13/cobra"
	"os"

	log "github.com/sirupsen/logrus"
)

var verbose bool

var rootCmd = &cobra.Command{
	Use:   "mynerva-author",
	SilenceErrors: true,
	SilenceUsage: true,

	// Setup before running any other command
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := setUpLog(verbose); err != nil {
			panic(err)
		}
		return nil
	},
}

func Execute() {
	initFlags()

	initAuth(rootCmd)
	initCodex(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		format := "error: %s\n"
		if verbose {
			format = "error: %#+v\n"
		}
		_, _ = fmt.Fprintf(os.Stderr, format, err)
		os.Exit(1)
	}
}

func initFlags() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVar(&config.MynervaApiHost, "api-host", config.MynervaApiHost, "Mynerva API host")
}

func setUpLog(verbose bool) error {
	println("setUpLog")
	if verbose {
		log.SetLevel(log.DebugLevel)
		log.Debugf("enabled debug logging")
	}
	return nil
}
