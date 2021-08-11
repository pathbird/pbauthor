package codex

import (
	"github.com/pathbird/pbauthor/internal/auth"
	"github.com/pathbird/pbauthor/internal/codex"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"path/filepath"
)

var codexInitConfig struct {
	systemPackages []string
	codexName      string
}

var codexInitCmd = &cobra.Command{
	Use:   "init [<path>]",
	Short: "initialize a codex configuration file",

	RunE: func(cmd *cobra.Command, args []string) error {
		// If no dir is specified, use current directory.
		if len(args) == 0 {
			args = append(args, ".")
		}

		if len(args) != 1 {
			return cmd.Usage()
		}

		dir, err := filepath.Abs(args[0])
		if err != nil {
			return errors.Wrap(err, "invalid codex directory")
		}

		auth, err := auth.GetAuth()
		if err != nil {
			return err
		}
		if auth == nil {
			return errors.New("not authenticated")
		}

		config, err := codex.InitConfig(dir)
		if err != nil {
			return errors.Wrap(err, "failed to initialize codex config")
		}

		if len(codexInitConfig.systemPackages) != 0 {
			config.Kernel.SystemPackages = codexInitConfig.systemPackages
		}

		if codexInitConfig.codexName != "" {
			config.Upload.Name = codexInitConfig.codexName
		}

		if err := config.Save(); err != nil {
			return errors.Wrap(err, "failed to save codex config")
		}
		return nil
	},
}

func init() {
	codexInitCmd.Flags().StringArrayVar(
		&codexInitConfig.systemPackages,
		"system-packages",
		[]string{},
		"a list of additional system packages to install",
	)
	codexInitCmd.Flags().StringVar(
		&codexInitConfig.codexName,
		"name",
		"",
		"the name of the codex as displayed in the Pathbird UI",
	)
	Cmd.AddCommand(codexInitCmd)
}
