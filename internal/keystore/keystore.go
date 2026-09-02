package keystore

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
)

const (
	privateKeyFile = "identity.private.asc"
	publicKeyFile  = "identity.public.asc"
)

var ErrIdentityExists = errors.New("a personal key already exists")

type Store struct {
	Dir string
}

func DefaultStore() (Store, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return Store{}, fmt.Errorf("finding user configuration directory: %w", err)
	}
	return Store{Dir: filepath.Join(configDir, "tacitus")}, nil
}

func (s Store) PrivatePath() string { return filepath.Join(s.Dir, privateKeyFile) }
func (s Store) PublicPath() string  { return filepath.Join(s.Dir, publicKeyFile) }

func (s Store) Generate() (string, error) {
	if s.Dir == "" {
		return "", errors.New("identity directory is empty")
	}
	if err := os.MkdirAll(s.Dir, 0o700); err != nil {
		return "", fmt.Errorf("creating identity directory: %w", err)
	}
	if err := os.Chmod(s.Dir, 0o700); err != nil {
		return "", fmt.Errorf("protecting identity directory: %w", err)
	}
	privateExists := pathExists(s.PrivatePath())
	publicExists := pathExists(s.PublicPath())

	switch {
	case privateExists && publicExists:
		return "", ErrIdentityExists
	case privateExists && !publicExists:
		return s.rebuildPublic()
	case !privateExists && publicExists:
		if err := os.Remove(s.PublicPath()); err != nil {
			return "", fmt.Errorf("removing orphaned public key: %w", err)
		}
	}

	entity, err := openpgp.NewEntity("tacitus", "", "", &packet.Config{
		Algorithm: packet.PubKeyAlgoEdDSA,
		Curve:     packet.Curve25519,
	})
	if err != nil {
		return "", fmt.Errorf("generating OpenPGP identity: %w", err)
	}

	privateTemp, err := s.writeTemporary(privateKeyFile, 0o600, func(w io.Writer) error {
		return writeArmored(w, "PGP PRIVATE KEY BLOCK", func(aw io.Writer) error {
			return entity.SerializePrivate(aw, nil)
		})
	})
	if err != nil {
		return "", err
	}
	defer os.Remove(privateTemp)

	publicTemp, err := s.writeTemporary(publicKeyFile, 0o644, func(w io.Writer) error {
		return writeArmored(w, "PGP PUBLIC KEY BLOCK", entity.Serialize)
	})
	if err != nil {
		return "", err
	}
	defer os.Remove(publicTemp)

	if err := publish(privateTemp, s.PrivatePath()); err != nil {
		return "", err
	}
	if err := publish(publicTemp, s.PublicPath()); err != nil {
		_ = os.Remove(s.PrivatePath())
		return "", err
	}

	return fmt.Sprintf("%X", entity.PrimaryKey.Fingerprint), nil
}

func (s Store) rebuildPublic() (string, error) {
	entity, err := s.LoadPrivate()
	if err != nil {
		return "", fmt.Errorf("repairing public key: %w", err)
	}

	publicTemp, err := s.writeTemporary(publicKeyFile, 0o644, func(w io.Writer) error {
		return writeArmored(w, "PGP PUBLIC KEY BLOCK", entity.Serialize)
	})
	if err != nil {
		return "", err
	}
	defer os.Remove(publicTemp)

	if err := publish(publicTemp, s.PublicPath()); err != nil {
		return "", err
	}

	return fmt.Sprintf("%X", entity.PrimaryKey.Fingerprint), nil
}

func (s Store) LoadPublic() (*openpgp.Entity, error) {
	return loadEntity(s.PublicPath(), "public")
}

func (s Store) LoadPrivate() (*openpgp.Entity, error) {
	entity, err := loadEntity(s.PrivatePath(), "private")
	if err != nil {
		return nil, err
	}
	if entity.PrivateKey == nil {
		return nil, errors.New("managed private key contains no private material")
	}
	return entity, nil
}

func loadEntity(path, kind string) (*openpgp.Entity, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("no personal key found; run tacitus keygen")
	}
	if err != nil {
		return nil, fmt.Errorf("opening %s key: %w", kind, err)
	}
	defer f.Close()

	entities, err := openpgp.ReadArmoredKeyRing(f)
	if err != nil {
		return nil, fmt.Errorf("reading %s key: %w", kind, err)
	}
	if len(entities) != 1 {
		return nil, fmt.Errorf("managed %s key must contain exactly one identity", kind)
	}
	if _, ok := entities[0].EncryptionKey(entityTime()); !ok {
		return nil, fmt.Errorf("managed %s key has no usable encryption key", kind)
	}
	return entities[0], nil
}

func entityTime() time.Time { return time.Now() }

func (s Store) writeTemporary(name string, mode os.FileMode, write func(io.Writer) error) (string, error) {
	f, err := os.CreateTemp(s.Dir, "."+name+".*")
	if err != nil {
		return "", fmt.Errorf("creating temporary key file: %w", err)
	}
	path := f.Name()
	if err := f.Chmod(mode); err != nil {
		f.Close()
		os.Remove(path)
		return "", fmt.Errorf("protecting temporary key file: %w", err)
	}
	if err := write(f); err != nil {
		f.Close()
		os.Remove(path)
		return "", fmt.Errorf("writing key file: %w", err)
	}
	if err := f.Close(); err != nil {
		os.Remove(path)
		return "", fmt.Errorf("closing key file: %w", err)
	}
	return path, nil
}

func writeArmored(w io.Writer, blockType string, write func(io.Writer) error) error {
	aw, err := armor.Encode(w, blockType, nil)
	if err != nil {
		return err
	}
	if err := write(aw); err != nil {
		aw.Close()
		return err
	}
	return aw.Close()
}

func publish(tempPath, destination string) error {
	if err := os.Link(tempPath, destination); err != nil {
		if os.IsExist(err) {
			return ErrIdentityExists
		}
		return fmt.Errorf("saving key file: %w", err)
	}
	return os.Remove(tempPath)
}

func pathExists(path string) bool {
	_, err := os.Lstat(path)
	return err == nil
}
