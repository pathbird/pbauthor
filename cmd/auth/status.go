package auth

import (
	"fmt"
	"github.com/pathbird/pbauthor/internal/auth"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"time"
)

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check authentication status",

	RunE: func(cmd *cobra.Command, args []string) error {
		auth, err := auth.GetAuth()
		if err != nil {
			return err
		}
		if auth == nil {
			return errors.New("not authenticated")
		}

		fmt.Printf("✅ Authenticated (until %s)\n", auth.Expiration.Format(time.RFC1123))
		return nil
	},
}

func init() {
	Cmd.AddCommand(authStatusCmd)
}
