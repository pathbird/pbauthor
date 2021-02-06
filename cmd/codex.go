package cmd

import (
	"fmt"
	"github.com/mynerva-io/author-cli/internal/api"
	"github.com/mynerva-io/author-cli/internal/auth"
	"github.com/mynerva-io/author-cli/internal/codex"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"path/filepath"
)

var codexCmd = &cobra.Command{
	Use:   "codex",
	Short: "Work with codices",
}

var codexUploadCmd = &cobra.Command{
	Use: "upload <path>",

	RunE: func(cmd *cobra.Command, args []string) error {
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

		client := api.New(auth.ApiToken)
		res, err := codex.UploadCodex(client, &codex.UploadCodexOptions{
			Dir: dir,
		})
		if err != nil {
			return err
		}

		fmt.Printf("codexId: %s", res.CodexId)
		return nil
	},
}

func initCodex(cmd *cobra.Command) {
	codexCmd.AddCommand(codexUploadCmd)

	cmd.AddCommand(codexCmd)
}
