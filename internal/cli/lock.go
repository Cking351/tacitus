package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Cking351/tacitus/internal/crypto"
	"github.com/Cking351/tacitus/internal/helper"
)

func lockCmd() *cobra.Command {
	var deleteOriginal bool;
	var outputPath  string;
	var passphrase string;
	var useArmor bool;
	cmd := &cobra.Command{
		Use:   "lock <file>",
		Short: "Encrypt a file",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inputPath := args[0]
			outPath := outputPath
			if outPath == "" {
				outPath = inputPath + ".tct"
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

			if passphrase == "" {
				passphrase, err = helper.ReadPassphrase("Passphrase: ")
				if err != nil {
					return err
				}
				confirmed, err := helper.ReadPassphrase("Confirm Passphrase: ")
				if err != nil {
					return err
				}
				if passphrase != confirmed {
					os.Remove(outPath)
					return errors.New("passphrase does not match")
				}
			}

			if err := crypto.EncryptSymmetric(inFile, outFile, passphrase, useArmor); err != nil {
				return fmt.Errorf("encrypting: %w", err)
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

	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "output file path (default: <file>.tct)")
	cmd.Flags().BoolVarP(&deleteOriginal, "delete", "d", false, "delete the original file after successful encryption")
	cmd.Flags().StringVarP(&passphrase, "password", "p", "", "password for encryption (THIS IS INSECURE)")
	cmd.Flags().BoolVarP(&useArmor, "armor", "a", false, "Use PGP armor")

	return cmd
}
