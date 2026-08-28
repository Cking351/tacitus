// Package crypto wraps github.com/ProtonMail/go-crypto/openpgp for tacitus.
package crypto


import (
	"bufio"
	"bytes"
	"errors"
	"io"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	pgperrors "github.com/ProtonMail/go-crypto/openpgp/errors"
)

var ErrCorruptFile = errors.New("file is corrupted or not a valid tacitus file")


func EncryptSymmetric(r io.Reader, w io.Writer, passphrase string, useArmor bool) error {
	out := w
	var armorWriter io.WriteCloser
	if useArmor {
		var err error
		armorWriter, err = armor.Encode(w, "PGP MESSAGE", nil)
		if err != nil {
			return err
		}
		out = armorWriter
	}
	plaintextWriter, err := openpgp.SymmetricallyEncrypt(out, []byte(passphrase), nil, nil)
	if err != nil {
		return err
	}
	if _, err := io.Copy(plaintextWriter, r); err != nil {
		plaintextWriter.Close()
		return err
	}
	// Close, don't defer: these writers flush the final packets, and a dropped
	// close error means a silently truncated ciphertext.
	if err := plaintextWriter.Close(); err != nil {
		return err
	}
	if armorWriter != nil {
		return armorWriter.Close()
	}
	return nil
}

func DecryptSymmetric(r io.Reader, w io.Writer, passphrase string) error {
	tried := false

	prompt := func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		if tried {
			// Second call
			return nil, errors.New("incorrect passphrase")
		}
		tried = true
		return []byte(passphrase), nil
	}

	pgpOpening := "-----BEGIN"
	br := bufio.NewReader(r)
	peek, _ := br.Peek(len(pgpOpening))

	body := io.Reader(br)
	if bytes.HasPrefix(peek, []byte(pgpOpening)) {
		block, err := armor.Decode(br)
		if err != nil {
			return translateError(err)
		}
		body = block.Body
	}
	
	md, err := openpgp.ReadMessage(body, nil, prompt, nil)
	if err != nil {
		return translateError(err)
	}

	_, err = io.Copy(w, md.UnverifiedBody)
	return translateError(err)
}

// Translate so we don't expose raw crypto returns
func translateError(err error) error {
	var sigErr pgperrors.SignatureError
	var structErr pgperrors.StructuralError
	if errors.As(err, &sigErr) || errors.As(err, &structErr) {
		return ErrCorruptFile
	}
	return err
}