package crypto

import (
	"bytes"
	"errors"
	"testing"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
)

func TestEncryptDecryptRoundTrip(t *testing.T) {
	plaintext := []byte("hello, tacitus")
	passphrase := "correct horse battery staple"

	var ciphertext bytes.Buffer
	if err := EncryptSymmetric(bytes.NewReader(plaintext), &ciphertext, passphrase, false); err != nil {
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
	if err := EncryptSymmetric(bytes.NewReader(plaintext), &ciphertext, "right passphrase", false); err != nil {
		t.Fatalf("EncryptSymmetric: %v", err)
	}

	var decrypted bytes.Buffer
	if err := DecryptSymmetric(&ciphertext, &decrypted, "wrong passphrase"); err == nil {
		t.Fatal("expected an error decrypting with the wrong passphrase, got nil")
	}
}

func TestDecryptTamperedCiphertextFails(t *testing.T) {
	var ciphertext bytes.Buffer
	if err := EncryptSymmetric(bytes.NewReader([]byte("hello, tacitus")), &ciphertext, "passphrase", false); err != nil {
		t.Fatalf("EncryptSymmetric: %v", err)
	}
	data := ciphertext.Bytes()
	data[len(data)-1] ^= 1

	var decrypted bytes.Buffer
	if err := DecryptSymmetric(bytes.NewReader(data), &decrypted, "passphrase"); err == nil {
		t.Fatal("expected an error decrypting tampered ciphertext, got nil")
	}
}

func TestEncryptDecryptPublicRoundTrip(t *testing.T) {
	entity := testEntity(t)
	plaintext := []byte("personal-key encryption")

	var ciphertext bytes.Buffer
	if err := EncryptPublic(bytes.NewReader(plaintext), &ciphertext, entity, true); err != nil {
		t.Fatalf("EncryptPublic: %v", err)
	}
	var decrypted bytes.Buffer
	if err := DecryptPrivate(&ciphertext, &decrypted, entity); err != nil {
		t.Fatalf("DecryptPrivate: %v", err)
	}
	if !bytes.Equal(decrypted.Bytes(), plaintext) {
		t.Fatalf("got %q, want %q", decrypted.Bytes(), plaintext)
	}
}

func TestDecryptPrivateRejectsWrongIdentity(t *testing.T) {
	plaintext := []byte("personal-key encryption")
	var ciphertext bytes.Buffer
	if err := EncryptPublic(bytes.NewReader(plaintext), &ciphertext, testEntity(t), false); err != nil {
		t.Fatalf("EncryptPublic: %v", err)
	}
	var decrypted bytes.Buffer
	err := DecryptPrivate(&ciphertext, &decrypted, testEntity(t))
	if !errors.Is(err, ErrWrongKey) {
		t.Fatalf("DecryptPrivate error = %v, want ErrWrongKey", err)
	}
}

func testEntity(t *testing.T) *openpgp.Entity {
	t.Helper()
	entity, err := openpgp.NewEntity("test", "", "", &packet.Config{
		Algorithm: packet.PubKeyAlgoEdDSA,
		Curve:     packet.Curve25519,
	})
	if err != nil {
		t.Fatalf("NewEntity: %v", err)
	}
	return entity
}
