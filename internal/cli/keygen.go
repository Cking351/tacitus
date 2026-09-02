package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Cking351/tacitus/internal/keystore"
)

func keygenCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "keygen",
		Short: "Generate a personal encryption key",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := keystore.DefaultStore()
			if err != nil {
				return err
			}
			fingerprint, err := store.Generate()
			if err != nil {
				return fmt.Errorf("generating personal key: %w", err)
			}
			fmt.Println("Created personal encryption key.")
			fmt.Printf("Fingerprint: %s\n", fingerprint)
			fmt.Printf("Private key: %s\n", store.PrivatePath())
			fmt.Println("Back up the private key securely. Anyone with it can decrypt your key-encrypted files.")
			return nil
		},
	}
}
