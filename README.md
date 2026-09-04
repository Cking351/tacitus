# tacitus

Encrypt and decrypt personal files from the command line. By default tacitus
uses a passphrase; an optional personal key mode lets you encrypt files for
cloud storage without entering a file passphrase. Built on
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

- **Simple default**: `lock` and `unlock` use a passphrase with no key setup.
- **Standard format**: output is a normal OpenPGP encrypted message,
  optionally ASCII-armored, so it's readable by other OpenPGP tools.
- **Personal key mode**: generate one local identity and use it to encrypt
  personal files without a file passphrase.

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
| `-f, --force` | Overwrite the output file if it already exists. |
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
| `-f, --force` | Overwrite the output file if it already exists. |
| `-p, --password <pw>` | Supply the passphrase directly. Same caveat as above — prefer the interactive prompt. |

Armored (`-a`/`--armor`) input from `lock` is detected automatically; you
don't need to pass any flag for it on `unlock`.

### Personal key mode

For personal files such as cloud backups, create a managed OpenPGP identity:

```
tacitus keygen
```

This creates a private key and matching public key in tacitus's per-user
configuration directory. Use `--key` (or `-k`) to select this mode:

```
tacitus lock photos.zip --key
tacitus unlock photos.zip.tct --key
```

Key mode never prompts for a file passphrase. `--key` cannot be combined with
`--password`. The private key is intentionally not passphrase-protected, so
the operating system account and full-disk encryption protect it.

Back up the private key with:

```
tacitus export /path/to/backup-location
```

Pass `-f`/`--force` to overwrite an existing file at the destination.

Store the export somewhere other than the same cloud storage it's meant to
protect — otherwise losing access to that storage loses the key too. Anyone
with the private key can decrypt key-encrypted files, and losing it
permanently loses access to them, so keep the export secure as well.

## Security notes

- Passphrase mode uses OpenPGP symmetric (password-based) encryption; personal
  key mode uses OpenPGP public-key encryption. Both are produced by
  ProtonMail's `go-crypto` and yield standard OpenPGP messages, with or without
  ASCII armor.
- Passphrase strength is entirely up to you — tacitus does not enforce a
  minimum length or complexity. Anyone who can guess or brute-force your
  passphrase can decrypt the file.
- Prefer the interactive prompt over `--password`/`-p`. A password passed on
  the command line is visible to anyone who can list processes on the
  machine, and is likely to end up in your shell's history file.
- `--delete` removes the original file with a normal filesystem delete, not
  a secure wipe. The data may still be recoverable from disk until it's
  overwritten.
- In personal key mode, the private-key file is the decryption secret. Keep it
  out of cloud storage and back it up somewhere safe.

## License

[BSD 3-Clause](LICENSE)
