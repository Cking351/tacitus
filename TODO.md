# tacitus TODO

Flat, ordered checklist to work through top to bottom. See `PLAN.md` for the
reasoning behind each milestone and the decisions already made.

Since this is your first crypto project: the golden rule is **never invent
your own crypto primitives** — everything here should go through
`go-crypto`'s implementations. Your job is wiring, not designing algorithms.
Also test the *negative* paths (wrong passphrase, corrupted file) as
seriously as the happy path — that's what actually proves the code is
trustworthy, not just that round-trip works once.

## Milestone 0 — Scaffolding ✅ done (`c527b3c`)

- [x] `go mod init`, add `go-crypto` + `cobra` deps
- [x] Directory layout + cobra command skeleton (stubs return "not yet implemented")

## Milestone 1 — Core crypto (symmetric)

- [x] **Implement `EncryptSymmetric`.**
  Use `openpgp.SymmetricallyEncrypt(w io.Writer, passphrase []byte, hints *openpgp.FileHints, config *packet.Config) (io.WriteCloser, error)`.
  It hands back a `WriteCloser` — you write your **plaintext** into that, and
  ciphertext streams out to `w` as you go. Run `go doc
  github.com/ProtonMail/go-crypto/openpgp SymmetricallyEncrypt` to confirm
  the exact signature for the version you pulled (v1.4.1).
  ⚠️ Gotcha: you must `Close()` the returned writer when you're done writing
  plaintext, or the final block never gets flushed and the output is
  truncated/corrupt. `defer` it carefully — check the `Close()` error too,
  don't just defer-and-ignore.
- [x] **Implement `DecryptSymmetric`.**
  Use `openpgp.ReadMessage(r io.Reader, keyring openpgp.KeyRing, prompt openpgp.PromptFunction, config *packet.Config) (*openpgp.MessageDetails, error)`.
  For symmetric-only decryption, `keyring` can be `nil`. `prompt` is a
  callback go-crypto invokes to *ask you* for the passphrase — its signature
  is roughly `func(keys []openpgp.Key, symmetric bool) ([]byte, error)`; when
  `symmetric` is `true`, just return the passphrase. Read plaintext back out
  of `md.UnverifiedBody` (an `io.Reader`) — despite the name, this is normal
  for symmetric messages; the "unverified" part refers to signature
  verification, which doesn't apply here.
- [x] **Round-trip test.** Encrypt a small buffer, decrypt it, assert the
  output equals the input byte-for-byte. Use `bytes.Buffer` for in-memory
  streams — no need to touch the filesystem for this test.
- [x] **Wrong-passphrase test.** Encrypt with one passphrase, try to decrypt
  with a different one, assert you get an error. This is the test that
  actually proves the passphrase is doing something — skipping it is the
  most common way to ship crypto code that silently doesn't protect anything.

## Milestone 2 — `lock` / `unlock` commands

- [x] **`lock`**: prompt for a passphrase (plain `fmt.Scanln` is fine for
  now — no-echo terminal input is milestone 4, don't block on it here), open
  the source file, create `<file>.tct`, call `EncryptSymmetric`, close both.
- [x] **`unlock`**: same shape in reverse. Restore the original filename by
  stripping `.tct`, or honor `--output` if given.
- [ ] **`--delete` flag**: only remove the original *after* the encrypted
  file has been fully written and closed without error. Sequence matters —
  if you delete first and encryption then fails, you've lost the file.
- [ ] **Manual smoke test**: `tacitus lock foo.txt && tacitus unlock
  foo.txt.tct` and diff the result against the original by hand.

## Milestone 3 — Keypairs (asymmetric)

- [ ] **Generate a keypair.**
  `openpgp.NewEntity(name, comment, email string, config *packet.Config) (*openpgp.Entity, error)`
  builds a primary key plus subkeys and self-signs the identity. This is the
  most API-heavy milestone — read `go doc` output for `openpgp.Entity`
  before diving in, there's more structure here (subkeys, identities,
  signatures) than the symmetric path.
- [ ] **Storage.** Serialize the public half with `entity.Serialize(w)` and
  the private half with `entity.SerializePrivate(w, config)`. Pick a
  location (e.g. `~/.config/tacitus/keys/`) and set the private key file to
  `0600` — that's not optional, a world-readable private key defeats the
  point.
- [ ] **Passphrase-protect the private key at rest.** Before serializing,
  call `entity.PrivateKey.Encrypt(passphrase)` — and do the same for each
  entry in `entity.Subkeys[i].PrivateKey`. Miss the subkeys and you've left
  part of the private key material unprotected on disk.
- [ ] **`keygen` command**: prompt for name/email (used as the OpenPGP
  identity string) and a passphrase, wire up the above.
- [ ] **`keys list`**: read every key file in the store, parse with
  `openpgp.ReadKeyRing`/`ReadArmoredKeyRing`, print
  `entity.PrimaryKey.KeyIdString()` plus the identity name/email.
- [ ] **`keys import <path>`**: read someone else's *public* key file the
  same way, copy it into your public keyring store.
- [ ] **`lock --to <recipient>`**: different function from milestone 1 —
  `openpgp.Encrypt(w io.Writer, to []*openpgp.Entity, signed *openpgp.Entity, hints *openpgp.FileHints, config *packet.Config) (io.WriteCloser, error)`.
  Pass `nil` for `signed` for now (signing is a separate concern from
  encrypting — don't conflate them).
- [ ] **`unlock` auto-detection**: the same `prompt` callback from milestone
  1 handles both cases — go-crypto calls it with `symmetric=true` for
  passphrase messages and `symmetric=false` (plus a list of candidate keys)
  when it needs you to unlock a *private key* instead. One `unlock` command,
  one prompt function, branching on that bool.

## Milestone 4 — UX polish

- [ ] **No-echo passphrase input.** `go get golang.org/x/term`, then
  `term.ReadPassword(int(os.Stdin.Fd()))`. Swap this in for the `Scanln`
  placeholders from milestone 2.
- [ ] **Confirm-twice on `keygen`.** Prompt for the passphrase twice, compare
  with `bytes.Equal` (not `==` — you're comparing byte slices, not strings),
  reject on mismatch before generating anything.
- [ ] **Translate go-crypto errors.** Check for the library's specific
  sentinel/typed errors (e.g. incorrect-passphrase cases) rather than
  printing raw internal error strings — a user typing the wrong passphrase
  should see "wrong passphrase," not a stack of wrapped library errors.
- [ ] **`--armor` flag.** Wrap the output writer with
  `armor.Encode(w, "PGP MESSAGE", nil)` *before* passing it to
  `SymmetricallyEncrypt`/`Encrypt`. ⚠️ Gotcha: you now have two nested
  `WriteCloser`s — close the encryption writer first, *then* the armor
  writer, then the underlying file. Wrong order corrupts the armored output.
- [ ] **Progress indicator** for large files — only worth doing if you
  notice real files taking long enough to need it. Don't build this
  speculatively.

## Milestone 5 — Testing

- [ ] **Unit tests for `internal/crypto`**: round-trip (done in milestone 1),
  wrong-passphrase failure (done in milestone 1), plus a **tampered
  ciphertext test** — encrypt something, flip one byte in the output, assert
  decrypt fails. OpenPGP's integrity protection (MDC/AEAD) exists precisely
  to catch this; a test that doesn't check it isn't testing the part of
  "encryption" that actually matters most (confidentiality *and*
  tamper-evidence, not just confidentiality).
- [ ] **CLI integration tests**: build the real binary, exec it against
  fixture files, assert on stdout/exit codes/output files.
  `github.com/rogpeppe/go-internal/testscript` is a good fit for
  script-style CLI tests.
- [ ] **CI**: GitHub Actions workflow running `go test ./...` on every push.

## Milestone 6 — Packaging / release

- [ ] `goreleaser` config (darwin/linux/windows, amd64/arm64)
- [ ] Tag `v0.1.0` once `lock`/`unlock` + keygen are solid
- [ ] Public-facing README, then flip the repo public

## Stretch (post-v1)

- [ ] Recursive directory encryption
- [ ] Stdin/stdout piping
- [ ] Shell completion
- [ ] Homebrew tap
