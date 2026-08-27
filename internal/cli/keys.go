package cli

import (
	"errors"

	"github.com/spf13/cobra"
)

func keysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keys",
		Short: "Manage stored keys",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List stored keys",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: milestone 3 — wire up internal/keystore
			return errors.New("keys list: not yet implemented")
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "import <path>",
		Short: "Import a public key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: milestone 3 — wire up internal/keystore
			return errors.New("keys import: not yet implemented")
		},
	})

	return cmd
}
