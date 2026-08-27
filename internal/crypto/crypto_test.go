package crypto

import "testing"

// See TODO.md milestone 1.

func TestEncryptDecryptRoundTrip(t *testing.T) {
	t.Skip("TODO: encrypt a buffer, decrypt it, assert output == input")
}

func TestDecryptWrongPassphraseFails(t *testing.T) {
	t.Skip("TODO: encrypt with one passphrase, decrypt with another, assert error")
}

func TestDecryptTamperedCiphertextFails(t *testing.T) {
	t.Skip("TODO: encrypt, flip a byte in the ciphertext, assert decrypt errors")
}
