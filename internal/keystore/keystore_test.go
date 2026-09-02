package keystore

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/ProtonMail/go-crypto/openpgp"
)

func TestGenerateAndLoadIdentity(t *testing.T) {
	store := Store{Dir: t.TempDir() + "/identity"}
	fingerprint, err := store.Generate()
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	public, err := store.LoadPublic()
	if err != nil {
		t.Fatalf("LoadPublic: %v", err)
	}
	private, err := store.LoadPrivate()
	if err != nil {
		t.Fatalf("LoadPrivate: %v", err)
	}
	if got := stringFingerprint(public); got != fingerprint {
		t.Fatalf("public fingerprint = %s, want %s", got, fingerprint)
	}
	if got := stringFingerprint(private); got != fingerprint {
		t.Fatalf("private fingerprint = %s, want %s", got, fingerprint)
	}
	if private.PrivateKey == nil {
		t.Fatal("private identity has no private key")
	}

	info, err := os.Stat(store.PrivatePath())
	if err != nil {
		t.Fatalf("stat private key: %v", err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("private key mode = %o, want 600", got)
	}
}

func TestGenerateRefusesExistingIdentity(t *testing.T) {
	store := Store{Dir: t.TempDir() + "/identity"}
	if _, err := store.Generate(); err != nil {
		t.Fatalf("first Generate: %v", err)
	}
	if _, err := store.Generate(); !errors.Is(err, ErrIdentityExists) {
		t.Fatalf("second Generate error = %v, want ErrIdentityExists", err)
	}
}

func TestGenerateRebuildsMissingPublicKey(t *testing.T) {
	store := Store{Dir: t.TempDir() + "/identity"}
	fingerprint, err := store.Generate()
	if err != nil {
		t.Fatalf("first Generate: %v", err)
	}
	if err := os.Remove(store.PublicPath()); err != nil {
		t.Fatalf("removing public key: %v", err)
	}

	got, err := store.Generate()
	if err != nil {
		t.Fatalf("repair Generate: %v", err)
	}
	if got != fingerprint {
		t.Fatalf("repaired fingerprint = %s, want %s", got, fingerprint)
	}

	if _, err := store.LoadPublic(); err != nil {
		t.Fatalf("LoadPublic after repair: %v", err)
	}
	private, err := store.LoadPrivate()
	if err != nil {
		t.Fatalf("LoadPrivate after repair: %v", err)
	}
	if stringFingerprint(private) != fingerprint {
		t.Fatalf("private key changed during repair: got %s, want %s", stringFingerprint(private), fingerprint)
	}
}

func TestGenerateReplacesOrphanedPublicKey(t *testing.T) {
	store := Store{Dir: t.TempDir() + "/identity"}
	oldFingerprint, err := store.Generate()
	if err != nil {
		t.Fatalf("first Generate: %v", err)
	}
	if err := os.Remove(store.PrivatePath()); err != nil {
		t.Fatalf("removing private key: %v", err)
	}

	newFingerprint, err := store.Generate()
	if err != nil {
		t.Fatalf("second Generate: %v", err)
	}
	if newFingerprint == oldFingerprint {
		t.Fatal("expected a new identity, got the orphaned one's fingerprint")
	}

	private, err := store.LoadPrivate()
	if err != nil {
		t.Fatalf("LoadPrivate: %v", err)
	}
	if stringFingerprint(private) != newFingerprint {
		t.Fatalf("private fingerprint = %s, want %s", stringFingerprint(private), newFingerprint)
	}
}

func TestLoadPrivateMissingIdentity(t *testing.T) {
	store := Store{Dir: t.TempDir() + "/identity"}
	if _, err := store.LoadPrivate(); err == nil {
		t.Fatal("LoadPrivate succeeded without an identity")
	}
}

func stringFingerprint(entity *openpgp.Entity) string {
	return fmt.Sprintf("%X", entity.PrimaryKey.Fingerprint)
}
