# tacitus build plan

You're driving the implementation — this is a roadmap, not code I'm going to
write for you. Ping me for API lookups, design sanity checks, or debugging
along the way.

**Decided so far:** cobra for the CLI framework · symmetric/passphrase
encryption before keypairs · original file kept by default, deleted only via
`--delete`.

## 0. Scaffolding

- [ ] `go mod init github.com/Cking351/tacitus`
- [ ] `go get github.com/ProtonMail/go-crypto/openpgp`
- [ ] Directory layout:
  ```
  cmd/tacitus/main.go   entrypoint, wires cobra root cmd
  internal/crypto/      wraps ProtonMail/go-crypto
  internal/keystore/    key generation, storage, lookup
  internal/cli/         cobra command definitions
  ```

## 1. Core crypto — symmetric (passphrase) encryption

Smallest end-to-end slice; proves the go-crypto wiring before key management
complexity gets added. Stream via `io.Reader`/`io.Writer` from the start —
`openpgp.SymmetricallyEncrypt` is already stream-based, and it keeps large
files off the heap.

- [ ] `internal/crypto.EncryptSymmetric(r io.Reader, w io.Writer, passphrase string) error`
- [ ] `internal/crypto.DecryptSymmetric(r io.Reader, w io.Writer, passphrase string) error`
- [ ] Round-trip test: encrypt a file, decrypt it, diff against the original

## 2. `lock` / `unlock` commands (passphrase-only)

- [ ] `tacitus lock <file>` → prompts for passphrase (never accept via flag —
      shows up in shell history), writes `<file>.tct`
- [ ] `tacitus unlock <file>.tct` → prompts for passphrase, restores original
      name (strip `.tct`) or writes to `--output`
- [ ] `--delete` flag removes the original file after a successful encrypt
- [ ] Manual smoke test end-to-end before moving on

## 3. Keypair support (asymmetric)

Biggest complexity jump in the plan — worth its own pass once 0-2 feel
solid, rather than building alongside passphrase support.

- [ ] `internal/keystore`: generate a keypair via `openpgp.NewEntity`;
      storage location (likely `~/.config/tacitus/keys/`), private key
      permissions `600`
- [ ] `tacitus keygen` — prompts for name/email (OpenPGP identity) and a
      passphrase to protect the private key at rest
- [ ] `tacitus keys list` — fingerprints/identities of stored keys
- [ ] `tacitus keys import <path>` — import someone else's public key
- [ ] `tacitus lock <file> --to <recipient>` — encrypt with a recipient's
      public key instead of/in addition to a passphrase
- [ ] `tacitus unlock` auto-detects passphrase vs. private-key decryption
      from the PGP packet headers

## 4. UX polish

- [ ] Passphrase prompts: no echo (`golang.org/x/term.ReadPassword`),
      confirm-twice on keygen
- [ ] Translate go-crypto errors into plain messages (wrong passphrase,
      corrupt file, missing key)
- [ ] `--armor` flag for ASCII-armored output vs. default binary `.tct`
- [ ] Progress indicator for large files, if encryption isn't near-instant

## 5. Testing

- [ ] Unit tests for `internal/crypto` (round-trip, wrong-passphrase
      failure, tampered-ciphertext detection)
- [ ] CLI integration tests — build the binary, exec against fixture files
      (`testscript` is a good fit)
- [ ] CI: GitHub Actions running `go test ./...` on push

## 6. Packaging / release

- [ ] `goreleaser` config for cross-platform binaries (darwin/linux/windows,
      amd64/arm64)
- [ ] Tag `v0.1.0` once `lock`/`unlock` + keygen are solid
- [ ] Public-facing README (separate from the current dev-notes one) before
      flipping the repo public

## Later / stretch (not blocking v1)

- Recursive directory encryption (`tacitus lock ./some-dir`)
- Stdin/stdout piping (`cat file | tacitus lock - > file.tct`)
- Shell completion (cobra generates this almost for free)
- Homebrew tap
