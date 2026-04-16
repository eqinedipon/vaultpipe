# vaultpipe

> Stream secrets from Vault into process environments without writing to disk.

---

## Installation

```bash
go install github.com/youruser/vaultpipe@latest
```

Or download a pre-built binary from the [releases page](https://github.com/youruser/vaultpipe/releases).

---

## Usage

`vaultpipe` fetches secrets from HashiCorp Vault and injects them as environment variables into a child process — no temp files, no disk writes.

```bash
vaultpipe run --path secret/data/myapp -- ./myapp server
```

The secrets stored at `secret/data/myapp` are injected directly into `myapp`'s environment at runtime.

### Options

| Flag | Description |
|------|-------------|
| `--path` | Vault secret path to read from |
| `--addr` | Vault server address (default: `$VAULT_ADDR`) |
| `--token` | Vault token (default: `$VAULT_TOKEN`) |
| `--mount` | Secret engine mount point (default: `secret`) |

### Example

```bash
export VAULT_ADDR=https://vault.example.com
export VAULT_TOKEN=s.xxxxxxxx

vaultpipe run --path secret/data/db-creds -- python manage.py runserver
```

The child process receives `DB_USER`, `DB_PASSWORD`, etc. as normal environment variables.

---

## Requirements

- Go 1.21+
- A running HashiCorp Vault instance

---

## License

MIT © [youruser](https://github.com/youruser)