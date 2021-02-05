package cmd

import (
	"github.com/mynerva-io/author-cli/internal/auth"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var codexCmd = &cobra.Command{
	Use: "codex",
	Short: "Work with codices",

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		auth, err := auth.GetAuth()
		if err != nil {
			return err
		}
		if auth == nil {
			return errors.New("not authenticated")
		}
		return nil
	},
}

var codexUploadCmd = &cobra.Command{
	Use: "upload",

	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("not implemented")
	},
}

func initCodex(cmd *cobra.Command) {
	codexCmd.AddCommand(codexUploadCmd)

	cmd.AddCommand(codexCmd)
}