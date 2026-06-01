---
title: General
sidebar_position: 1
---

# Configuration: General

`dns-updater` reads configuration from a `config.yaml` file in the current working directory.

## Top-level keys

- `provider` (required): selects which DNS provider to use and supplies provider-specific config
- `updates` (required): list of DNS updates to apply
- `log` (optional): logging configuration (defaults: `info` + `pretty`)
- `cache` (optional): enable/disable simple caching and configure cache file path
- `schedule` (optional): a cron expression to run the DNS check on a schedule (e.g. `"*/5 * * * *"` for every 5 minutes).
- `provider.safemode` (optional): enable/disable safe mode (defaults to `true`). When enabled, `dns-updater` uses TXT records with data in the format `managed-by:dns-updater/<txt_owner_id>` to mark records it owns. See [Safe Mode](./safe-mode.md).

## Updates

`updates` is a list. Each item supports:

- `domain`: the base domain (e.g. `example.com`)
- `zone`: provider zone identifier/name (provider-specific)
- `type`: DNS record type (commonly `A`)
- `records`: list of record names to update (e.g. `"@"`, `www`, `home`)

Example:

```yaml
updates:
  - domain: example.com
    zone: example-com
    type: A
    records:
      - "@"
      - "www"
```
