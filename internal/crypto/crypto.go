// Package crypto wraps github.com/ProtonMail/go-crypto/openpgp for tacitus.
package crypto

import (
	"errors"
	"io"
)

// EncryptSymmetric encrypts r into w using passphrase-based symmetric
// encryption. See PLAN.md milestone 1.
func EncryptSymmetric(r io.Reader, w io.Writer, passphrase string) error {
	return errors.New("crypto: EncryptSymmetric not yet implemented")
}

// DecryptSymmetric decrypts r into w using passphrase-based symmetric
// decryption. See PLAN.md milestone 1.
func DecryptSymmetric(r io.Reader, w io.Writer, passphrase string) error {
	return errors.New("crypto: DecryptSymmetric not yet implemented")
}
