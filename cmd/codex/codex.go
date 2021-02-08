package codex

import (
	"github.com/spf13/cobra"
)

var codexCmd = &cobra.Command{
	Use:   "codex",
	Short: "Work with codices",
}

func InitCodexCmd(cmd *cobra.Command) {
	codexCmd.AddCommand(codexUploadCmd)

	cmd.AddCommand(codexCmd)
}
