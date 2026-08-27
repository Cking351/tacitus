# tacitus (dev notes)

Internal notes only — not a public-facing README yet.

## What this is

A Go CLI for encrypting/decrypting files, built on ProtonMail's OpenPGP fork
(`github.com/ProtonMail/go-crypto`, the engine behind `gopenpgp`).

## Naming

- Landed on `tacitus` after ruling out `ion` (collides with `lfaoro/ion` and
  SST's "Ion" deploy engine), `ironvail`/`ironward` (iron-* felt crowded —
  IronKey, IronClad, IronCore Labs all taken; `ironward.com` is an active
  Croatian game studio), and anything Proton-branded (avoid implying an
  official Proton product).
- `tacitus` checked clean: no competing security/crypto tool, no npm
  collision, only an unrelated bioinformatics GitHub repo (`alaimos/tacitus`).
  Domain availability (`.com`/`.dev`/`.sh`) still unverified — check a
  registrar directly before announcing publicly.
- Roman historian Tacitus was known for terse, guarded prose — fits an
  encryption tool. Root of "tacit."

## Planned CLI shape

```
tacitus lock <file>                    # encrypt -> file.tct
tacitus unlock <file.tct>              # decrypt back to original
tacitus keygen                         # generate a new keypair
tacitus keys list
tacitus keys import <path>
tacitus lock <file> --to <recipient>   # encrypt for a recipient's pubkey
```

- Verbs: `lock` / `unlock` (not `encrypt`/`decrypt`)
- Encrypted extension: `.tct`
- Go module path: TBD — `github.com/<you>/tacitus`

## Dev loop

```
go build ./...                        # compile everything
go run ./cmd/tacitus <command>         # e.g. go run ./cmd/tacitus lock foo.txt
go test ./...                         # run tests
```

See `PLAN.md` for rationale/decisions and `TODO.md` for the current
milestone checklist.

## Open questions / next steps

- [ ] Confirm domain availability if we want one
- [ ] `go mod init`, pick CLI framework (cobra vs urfave/cli)
- [ ] Wire up ProtonMail/go-crypto for actual lock/unlock
- [ ] Decide key storage location/format (keyring dir, passphrase-protected?)
