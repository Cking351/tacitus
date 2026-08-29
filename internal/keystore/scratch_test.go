package keystore

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
)

func TestScratch(t *testing.T) {
	cfg := &packet.Config{
		Algorithm: packet.PubKeyAlgoEdDSA,
		Curve:     packet.Curve25519,
	}

	e, err := openpgp.NewEntity("tacitus vault", "", "chris@example.com", cfg)
	if err != nil {
		t.Fatalf("NewEntity: %v", err)
	}

	t.Logf("fingerprint		%X", e.PrimaryKey.Fingerprint)
	t.Logf("subkeys			%d", len(e.Subkeys))
	for id := range e.Identities {
		t.Logf("identity	%q", id)
	}
	t.Logf("primary algo %d", e.PrimaryKey.PubKeyAlgo)
	for i, sk := range e.Subkeys {
		t.Logf("subkey[%d] algo %d", i, sk.PublicKey.PubKeyAlgo)
	}

	dir, err := os.MkdirTemp("", "tacitus-")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	pubPath := filepath.Join(dir, "vault.pub")
	t.Logf("scratch dir		%s", dir)

	f, err := os.Create(pubPath)
	if err != nil {
		t.Fatalf("create %s: %v", pubPath, err)
	}
	defer f.Close()

	aw, err := armor.Encode(f, "PGP PUBLIC KEY BLOCK", nil)
	if err != nil {
		t.Fatalf("armor.Encode: %v", err)
	}
	if err := e.Serialize(aw); err != nil {
		t.Fatalf("Serialize: %v", err)
	}
	if err := aw.Close(); err != nil {
		t.Fatalf("closing armor writer: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("closing %s: %v", pubPath, err)
	}

	keyPath := filepath.Join(dir, "vault-unprotected.key")
	kf, err := os.Create(keyPath)
	if err != nil {
		t.Fatalf("create %s: %v", keyPath, err)
	}
	defer kf.Close()

	aw2, err := armor.Encode(kf, "PGP PRIVATE KEY BLOCK", nil)
	if err != nil {
		t.Fatalf("armor.Encode: %v", err)
	}
	if err := e.SerializePrivateWithoutSigning(aw2, nil); err != nil {
		t.Fatalf("SerializePrivateWithoutSigning: %v", err)
	}
	if err := aw2.Close(); err != nil {
		t.Fatalf("closing armor writer: %v", err)
	}
	if err := kf.Close(); err != nil {
		t.Fatalf("closing %s: %v", keyPath, err)
	}
}
