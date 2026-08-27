package cli

import (
	"errors"

	"github.com/spf13/cobra"
)

func keygenCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "keygen",
		Short: "Generate a new keypair",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: milestone 3 — wire up internal/keystore
			return errors.New("keygen: not yet implemented")
		},
	}
}
