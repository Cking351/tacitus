// Package crypto wraps github.com/ProtonMail/go-crypto/openpgp for tacitus.
package crypto


import (
	"errors"
	"io"
	"github.com/ProtonMail/go-crypto/openpgp"
)


func EncryptSymmetric(r io.Reader, w io.Writer, passphrase string) error {
	plaintextWriter, err := openpgp.SymmetricallyEncrypt(w, []byte(passphrase), nil, nil)
	if err != nil {
		return err
	}
	defer plaintextWriter.Close()
	_, err = io.Copy(plaintextWriter, r)
	return err
}

func DecryptSymmetric(r io.Reader, w io.Writer, passphrase string) error {
	tried := false

	prompt := func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		if tried {
			// Second call
			return nil, errors.New("crypto: incorrect passphrase")
		}
		tried = true
		return []byte(passphrase), nil
	}

	
	md, err := openpgp.ReadMessage(r, nil, prompt, nil)
	if err != nil {
		return err
	}


	_, err = io.Copy(w, md.UnverifiedBody)
	return err
}
