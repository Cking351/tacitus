package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/spf13/cobra"

	"github.com/Cking351/tacitus/internal/crypto"
	"github.com/Cking351/tacitus/internal/helper"
	"github.com/Cking351/tacitus/internal/keystore"
)

func unlockCmd() *cobra.Command {
	var deleteOriginal bool
	var force bool
	var outputPath string
	var passphrase string
	var useKey bool
	cmd := &cobra.Command{
		Use:   "unlock <file>",
		Short: "Decrypt a file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inputPath := args[0]
			outPath := outputPath
			if outPath == "" {
				outPath = strings.TrimSuffix(inputPath, ".tct")
				if outPath == inputPath {
					outPath = inputPath + ".decrypted"
				}
			}
			if outPath == inputPath {
				return errors.New("output path must differ from the input file")
			}

			inFile, err := openInput(inputPath)
			if err != nil {
				return err
			}
			defer inFile.Close()

			var identity *openpgp.Entity
			if useKey {
				store, err := keystore.DefaultStore()
				if err != nil {
					return err
				}
				identity, err = store.LoadPrivate()
				if err != nil {
					return err
				}
			} else if passphrase == "" {
				passphrase, err = helper.ReadPassphrase("Passphrase: ")
				if err != nil {
					return err
				}
			}

			outFile, err := createOutput(outPath, force)
			if err != nil {
				return err
			}
			defer outFile.Close()

			if useKey {
				if err := crypto.DecryptPrivate(inFile, outFile, identity); err != nil {
					cleanupOutput(outFile, outPath)
					return fmt.Errorf("decrypting: %w", err)
				}
			} else if err := crypto.DecryptSymmetric(inFile, outFile, passphrase); err != nil {
				cleanupOutput(outFile, outPath)
				return fmt.Errorf("decrypting: %w", err)
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

	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "output file path (default: strip .tct suffix)")
	cmd.Flags().BoolVarP(&deleteOriginal, "delete", "d", false, "delete the encrypted file after successful decryption")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "overwrite the output file if it already exists")
	cmd.Flags().StringVarP(&passphrase, "password", "p", "", "password for decryption (THIS IS INSECURE)")
	cmd.Flags().BoolVarP(&useKey, "key", "k", false, "decrypt with your personal key")
	cmd.MarkFlagsMutuallyExclusive("key", "password")

	return cmd
}
