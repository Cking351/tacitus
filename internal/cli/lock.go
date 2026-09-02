package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/spf13/cobra"

	"github.com/Cking351/tacitus/internal/crypto"
	"github.com/Cking351/tacitus/internal/helper"
	"github.com/Cking351/tacitus/internal/keystore"
)

func lockCmd() *cobra.Command {
	var deleteOriginal bool
	var force bool
	var outputPath string
	var passphrase string
	var useArmor bool
	var useKey bool
	cmd := &cobra.Command{
		Use:   "lock <file>",
		Short: "Encrypt a file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inputPath := args[0]
			outPath := outputPath
			if outPath == "" {
				outPath = inputPath + ".tct"
			}
			if outPath == inputPath {
				return errors.New("output path must differ from the input file")
			}

			inFile, err := openInput(inputPath)
			if err != nil {
				return err
			}
			defer inFile.Close()

			var recipient *openpgp.Entity
			if useKey {
				store, err := keystore.DefaultStore()
				if err != nil {
					return err
				}
				recipient, err = store.LoadPublic()
				if err != nil {
					return err
				}
			} else if passphrase == "" {
				passphrase, err = helper.ReadPassphrase("Passphrase: ")
				if err != nil {
					return err
				}
				confirmed, err := helper.ReadPassphrase("Confirm Passphrase: ")
				if err != nil {
					return err
				}
				if passphrase != confirmed {
					return errors.New("passphrase does not match")
				}
			}

			outFile, err := createOutput(outPath, force)
			if err != nil {
				return err
			}
			defer outFile.Close()

			if useKey {
				if err := crypto.EncryptPublic(inFile, outFile, recipient, useArmor); err != nil {
					cleanupOutput(outFile, outPath)
					return fmt.Errorf("encrypting: %w", err)
				}
			} else if err := crypto.EncryptSymmetric(inFile, outFile, passphrase, useArmor); err != nil {
				cleanupOutput(outFile, outPath)
				return fmt.Errorf("encrypting: %w", err)
			}
			if err := outFile.Close(); err != nil {
				cleanupOutput(outFile, outPath)
				return fmt.Errorf("writing %s: %w", outPath, err)
			}

			fmt.Printf("Wrote %s\n", outPath)

			if deleteOriginal {
				inFile.Close()
				if err := os.Remove(inputPath); err != nil {
					return fmt.Errorf("remove %s: %w", inputPath, err)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "output file path (default: <file>.tct)")
	cmd.Flags().BoolVarP(&deleteOriginal, "delete", "d", false, "delete the original file after successful encryption")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "overwrite the output file if it already exists")
	cmd.Flags().StringVarP(&passphrase, "password", "p", "", "password for encryption (THIS IS INSECURE)")
	cmd.Flags().BoolVarP(&useArmor, "armor", "a", false, "Use PGP armor")
	cmd.Flags().BoolVarP(&useKey, "key", "k", false, "encrypt with your personal key")
	cmd.MarkFlagsMutuallyExclusive("key", "password")

	return cmd
}
