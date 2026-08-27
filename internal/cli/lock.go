package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Cking351/tacitus/internal/crypto"
)

func lockCmd() *cobra.Command {
	var deleteOriginal bool;
	var outputPath  string;
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

			passphrase, err := readPassphrase("Passphrase: ")
			if err != nil {
				return err
			}

			if err := crypto.EncryptSymmetric(inFile, outFile, passphrase); err != nil {
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

	// TODO: --to <recipient>, --armor flags (milestones 2-4)
	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "output file path (default: <file>.tct)")
	cmd.Flags().BoolVarP(&deleteOriginal, "delete", "d", false, "delete the original file after successful encryption")

	return cmd
}

// readPassphrase prints prompt, reads a line from stdin, and trims the
// trailing newline. Plain text for now — no-echo terminal input is a
// milestone 4 upgrade.
func readPassphrase(prompt string) (string, error) {
	fmt.Print(prompt)
	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("reading passphrase: %w", err)
	}
	return strings.TrimSpace(line), nil
}
