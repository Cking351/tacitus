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
var ErrWrongKey = errors.New("file is not encrypted to this personal key")

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

	body, err := messageBody(r)
	if err != nil {
		return translateError(err)
	}

	md, err := openpgp.ReadMessage(body, nil, prompt, nil)
	if err != nil {
		return translateError(err)
	}

	_, err = io.Copy(w, md.UnverifiedBody)
	return translateError(err)
}

func EncryptPublic(r io.Reader, w io.Writer, recipient *openpgp.Entity, useArmor bool) error {
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
	plaintextWriter, err := openpgp.Encrypt(out, []*openpgp.Entity{recipient}, nil, nil, nil)
	if err != nil {
		if armorWriter != nil {
			armorWriter.Close()
		}
		return err
	}
	if _, err := io.Copy(plaintextWriter, r); err != nil {
		plaintextWriter.Close()
		if armorWriter != nil {
			armorWriter.Close()
		}
		return err
	}
	if err := plaintextWriter.Close(); err != nil {
		if armorWriter != nil {
			armorWriter.Close()
		}
		return err
	}
	if armorWriter != nil {
		return armorWriter.Close()
	}
	return nil
}

func DecryptPrivate(r io.Reader, w io.Writer, identity *openpgp.Entity) error {
	body, err := messageBody(r)
	if err != nil {
		return translateError(err)
	}
	md, err := openpgp.ReadMessage(body, openpgp.EntityList{identity}, nil, nil)
	if err != nil {
		return translateError(err)
	}
	_, err = io.Copy(w, md.UnverifiedBody)
	return translateError(err)
}

func messageBody(r io.Reader) (io.Reader, error) {
	const pgpOpening = "-----BEGIN"
	br := bufio.NewReader(r)
	peek, _ := br.Peek(len(pgpOpening))
	if !bytes.HasPrefix(peek, []byte(pgpOpening)) {
		return br, nil
	}
	block, err := armor.Decode(br)
	if err != nil {
		return nil, err
	}
	return block.Body, nil
}

// Translate so we don't expose raw crypto returns
func translateError(err error) error {
	if errors.Is(err, pgperrors.ErrKeyIncorrect) {
		return ErrWrongKey
	}
	var sigErr pgperrors.SignatureError
	var structErr pgperrors.StructuralError
	if errors.As(err, &sigErr) || errors.As(err, &structErr) {
		return ErrCorruptFile
	}
	return err
}
