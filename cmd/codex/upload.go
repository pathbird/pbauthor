package codex

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/mynerva-io/author-cli/internal/api"
	"github.com/mynerva-io/author-cli/internal/auth"
	"github.com/mynerva-io/author-cli/internal/codex"
	"github.com/mynerva-io/author-cli/internal/config"
	"github.com/mynerva-io/author-cli/internal/prompt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

// flag vars
var (
	skipConfirmation bool
)

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

		// This is a placeholder for now, but in the future, we'd like to confirm
		// the name of the codex and the course it's being uploaded to.
		if !skipConfirmation {
			if !prompt.Confirm("Upload codex?") {
				_, _ = fmt.Fprintln(os.Stderr, red("Upload aborted."))
				os.Exit(1)
			}
		}

		client := api.New(auth.ApiToken)
		res, parseErr, err := codex.UploadCodex(client, &codex.UploadCodexOptions{
			Dir: dir,
		})
		if err != nil {
			return err
		}
		if parseErr != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to parse codex (%d issues):\n", len(parseErr.Errors))
			for _, e := range parseErr.Errors {
				_, _ = fmt.Fprintf(os.Stderr, "  - %s\n    (%s at %s)\n", red(e.Message), faint(e.Error), cyan(e.SourcePosition))
			}
			os.Exit(1)
		}

		fmt.Printf("codexId: %s\n", res.CodexId)
		fmt.Printf("url: %s/codex/%s\n", config.MynervaApiHost, res.CodexId)

		return nil
	},
}

func init() {
	codexUploadCmd.Flags().BoolVarP(&skipConfirmation, "yes", "y", false, "don't ask for confirmation")
}

var (
	red   = color.New(color.FgRed, color.Bold).SprintFunc()
	cyan  = color.New(color.FgCyan).SprintFunc()
	faint = color.New(color.Faint).SprintFunc()
)
