package codex

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/pathbird/pbauthor/internal/api"
	"github.com/pathbird/pbauthor/internal/auth"
	"github.com/pathbird/pbauthor/internal/codex"
	"github.com/pathbird/pbauthor/internal/graphql"
	"github.com/pathbird/pbauthor/internal/prompt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"time"
)

// flag vars
var (
	skipConfirmation bool
	noWait           bool
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
				_, _ = fmt.Fprintln(os.Stderr, failf("Upload aborted."))
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
			_, _ = fmt.Fprintf(
				os.Stderr,
				"Failed to parse codex (%d issues):\n",
				len(parseErr.Errors),
			)
			for _, e := range parseErr.Errors {
				_, _ = fmt.Fprintf(os.Stderr, "- %s\n  (%s", failf(e.Message), blue(e.Error))
				if e.SourcePosition != "" {
					_, _ = fmt.Fprintf(os.Stderr, " at %s", cyan(e.SourcePosition))
				}

				_, _ = fmt.Fprintf(os.Stderr, ")\n")

				for _, line := range e.SourceInfo.SourceContext.Lines {
					_, _ = fmt.Fprintf(os.Stderr, "  %s %s\n", faint(">"), line)
				}
			}
			os.Exit(1)
		}

		detailsUrl := fmt.Sprintf("https://pathbird.com/codex/%s/details", res.CodexId)

		if !noWait {
			start := time.Now()
			timeoutCtx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
			defer cancel()
			log.Info("waiting for kernel build to complete (this may take up to 20 minutes)...")
			kernelStatus, err := codex.WaitForKernelBuildCompleted(
				timeoutCtx,
				graphql.NewClient(auth),
				res.CodexId,
			)
			if err != nil {
				log.WithError(
					err,
				).Error(
					"Something went wrong while trying to check the status of the codex.",
				)
				return err
			}
			log.Infof("waited %s for kernel build process", time.Since(start))
			if kernelStatus.BuildStatus != "built" {
				_, _ = fmt.Fprint(os.Stderr, failf(
					"Failed to build kernel (got status: %s): %s\n",
					kernelStatus.BuildStatus,
					detailsUrl,
				))
				return errors.Errorf(
					"failed to build kernel (got status: %s)",
					kernelStatus.BuildStatus,
				)
			}
		} else {
			log.Info("not waiting for kernel build to complete (--no-wait was set)")
		}

		fmt.Printf(successf("Successfully uploaded codex: %s", detailsUrl))

		return nil
	},
}

func init() {
	codexUploadCmd.Flags().BoolVarP(
		&skipConfirmation,
		"yes",
		"y",
		false,
		"don't ask for confirmation",
	)
	codexUploadCmd.Flags().BoolVar(
		&noWait,
		"no-wait",
		false,
		"don't wait for the kernel build process to complete",
	)
	Cmd.AddCommand(codexUploadCmd)
}

var (
	failf    = color.New(color.FgRed, color.Bold).SprintfFunc()
	successf = color.New(color.FgGreen, color.Bold).SprintfFunc()
	cyan     = color.New(color.FgCyan).SprintFunc()
	faint    = color.New(color.Faint).SprintFunc()
	blue     = color.New(color.FgBlue).SprintfFunc()
)
