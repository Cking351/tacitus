package helper

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func ReadPassphrase(prompt string) (string, error) {
	fmt.Print(prompt)
	line, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return "", fmt.Errorf("reading passphrase: %v", err)
	}
	return string(line), nil
}