package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Cking351/tacitus/internal/crypto"
)

func unlockCmd() *cobra.Command {
	var deleteOriginal bool;
	var outputPath string;
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

			inFile, err := os.Open(inputPath)
			if err != nil {
				return fmt.Errorf("opening %s: %w", inputPath, err)
			}
			defer inFile.Close()

			outFile, err := os.Create(outPath)
			if err != nil {
				return fmt.Errorf("creating %s: %w", outPath, err)
			}
			defer outFile.Close()

			passphrase, err := readPassphrase("Passphrase: ")
			if err != nil {
				return err
			}

			if err := crypto.DecryptSymmetric(inFile, outFile, passphrase); err != nil {
				return fmt.Errorf("decrypting: %w", err)
			}

			fmt.Printf("Wrote %s\n", outPath)

			if deleteOriginal {
				if err := os.Remove(inputPath); err != nil {
					return fmt.Errorf("remove %s: %w", inputPath, err)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "output file path (default: strip .tct suffix)")
	cmd.Flags().BoolVarP(&deleteOriginal, "delete", "d", false, "delete the encrypted file after successful decryption")

	return cmd
}
