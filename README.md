# tacitus

Encrypt and decrypt files from the command line, with a passphrase. Built on
[OpenPGP](https://www.openpgp.org/) via ProtonMail's `go-crypto`, so the file
format is standard and your encrypted files are never tied to this tool.

```
$ tacitus lock secrets.txt
Passphrase:
Confirm Passphrase:
Wrote secrets.txt.tct

$ tacitus unlock secrets.txt.tct
Passphrase:
Wrote secrets.txt
```

## Features

- **Simple**: two commands, `lock` and `unlock`. No key management to set up.
- **Standard format**: output is a normal OpenPGP symmetrically-encrypted
  message, optionally ASCII-armored, so it's readable by other OpenPGP tools.

## Installation

### Prerequisites

- [Go](https://go.dev/) 1.25 or later

### Build from source

```
git clone https://github.com/Cking351/tacitus.git
cd tacitus
go build -o tacitus ./cmd/tacitus
```

This puts a `tacitus` binary in the current directory. Move it onto your
`$PATH` (e.g. `mv tacitus /usr/local/bin/`) to run it from anywhere.

## Usage

### `lock` — encrypt a file

```
tacitus lock <file> [flags]
```

| Flag | Effect |
|---|---|
| `-o, --output <path>` | Output path. Defaults to `<file>.tct`. |
| `-a, --armor` | ASCII-armor the output (base64 text, ~33% larger). Use this when the file needs to survive being pasted into an email, chat, or other text-only channel. |
| `-d, --delete` | Delete the original file once encryption succeeds. |
| `-p, --password <pw>` | Supply the passphrase directly. **Avoid this** — it can be recorded in your shell history and process list. Omit it and you'll be prompted instead. |

Without `-p`, `lock` prompts for the passphrase twice and refuses to proceed
if the two don't match.

### `unlock` — decrypt a file

```
tacitus unlock <file> [flags]
```

| Flag | Effect |
|---|---|
| `-o, --output <path>` | Output path. Defaults to `<file>` with the `.tct` suffix stripped, or `<file>.decrypted` if there was no `.tct` suffix. |
| `-d, --delete` | Delete the encrypted file once decryption succeeds. |
| `-p, --password <pw>` | Supply the passphrase directly. Same caveat as above — prefer the interactive prompt. |

Armored (`-a`/`--armor`) input from `lock` is detected automatically; you
don't need to pass any flag for it on `unlock`.

## Security notes

- Encryption is OpenPGP symmetric (password-based) encryption, produced by
  ProtonMail's `go-crypto`. The output is a standard OpenPGP message, with or
  without ASCII armor, decryptable by any compliant OpenPGP implementation
  that knows the passphrase.
- Passphrase strength is entirely up to you — tacitus does not enforce a
  minimum length or complexity. Anyone who can guess or brute-force your
  passphrase can decrypt the file.
- Prefer the interactive prompt over `--password`/`-p`. A password passed on
  the command line is visible to anyone who can list processes on the
  machine, and is likely to end up in your shell's history file.
- `--delete` removes the original file with a normal filesystem delete, not
  a secure wipe. The data may still be recoverable from disk until it's
  overwritten.

## License

[BSD 3-Clause](LICENSE)
