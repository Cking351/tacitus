package cli

import (
	"errors"

	"github.com/spf13/cobra"
)

func unlockCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unlock <file>",
		Short: "Decrypt a file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: milestone 2 — wire up internal/crypto.DecryptSymmetric
			return errors.New("unlock: not yet implemented")
		},
	}

	// TODO: --output flag (milestone 2)

	return cmd
}
