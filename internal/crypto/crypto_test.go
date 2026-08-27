package crypto

import (
	"bytes"
	"testing"
)

// See TODO.md milestone 1.

func TestEncryptDecryptRoundTrip(t *testing.T) {
	plaintext := []byte("hello, tacitus")
	passphrase := "correct horse battery staple"

	var ciphertext bytes.Buffer
	if err := EncryptSymmetric(bytes.NewReader(plaintext), &ciphertext, passphrase); err != nil {
		t.Fatalf("EncryptSymmetric: %v", err)
	}

	var decrypted bytes.Buffer
	if err := DecryptSymmetric(&ciphertext, &decrypted, passphrase); err != nil {
		t.Fatalf("DecryptSymmetric: %v", err)
	}

	if !bytes.Equal(decrypted.Bytes(), plaintext) {
		t.Fatalf("got %q, want %q", decrypted.Bytes(), plaintext)
	}
}

func TestDecryptWrongPassphraseFails(t *testing.T) {
	plaintext := []byte("hello, tacitus")

	var ciphertext bytes.Buffer
	if err := EncryptSymmetric(bytes.NewReader(plaintext), &ciphertext, "right passphrase"); err != nil {
		t.Fatalf("EncryptSymmetric: %v", err)
	}

	var decrypted bytes.Buffer
	if err := DecryptSymmetric(&ciphertext, &decrypted, "wrong passphrase"); err == nil {
		t.Fatal("expected an error decrypting with the wrong passphrase, got nil")
	}
}

func TestDecryptTamperedCiphertextFails(t *testing.T) {
	t.Skip("TODO: encrypt, flip a byte in the ciphertext, assert decrypt errors")
}
