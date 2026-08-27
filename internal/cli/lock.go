package cli

import (
	"errors"

	"github.com/spf13/cobra"
)

func lockCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lock <file>",
		Short: "Encrypt a file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: milestone 2 — wire up internal/crypto.EncryptSymmetric
			return errors.New("lock: not yet implemented")
		},
	}

	// TODO: --to <recipient>, --output, --delete, --armor flags (milestones 2-4)

	return cmd
}
