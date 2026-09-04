package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Cking351/tacitus/internal/keystore"
)

func exportCmd() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "export <path>",
		Short: "Back up your personal private key to a file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			destPath := args[0]

			store, err := keystore.DefaultStore()
			if err != nil {
				return err
			}
			if destPath == store.PrivatePath() {
				return errors.New("destination must differ from the managed private key")
			}

			entity, err := store.LoadPrivate()
			if err != nil {
				return err
			}
			data, err := os.ReadFile(store.PrivatePath())
			if err != nil {
				return fmt.Errorf("reading private key: %w", err)
			}

			outFile, err := createOutput(destPath, force)
			if err != nil {
				return err
			}
			defer outFile.Close()

			if _, err := outFile.Write(data); err != nil {
				cleanupOutput(outFile, destPath)
				return fmt.Errorf("writing %s: %w", destPath, err)
			}
			if err := outFile.Close(); err != nil {
				cleanupOutput(outFile, destPath)
				return fmt.Errorf("writing %s: %w", destPath, err)
			}

			fmt.Printf("Exported private key to %s\n", destPath)
			fmt.Printf("Fingerprint: %X\n", entity.PrimaryKey.Fingerprint)
			fmt.Println("Store this somewhere safe and separate from your cloud backups — anyone with it can decrypt your key-encrypted files, and losing it permanently loses access to them.")
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "overwrite the destination if it already exists")

	return cmd
}
