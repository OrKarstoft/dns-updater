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
- `safeMode` (optional): enable/disable safe mode (defaults to `true`). When enabled, `dns-updater` will only update A records if a corresponding TXT record with the value `managed-by-dns-updater` exists. This prevents accidental updates to records not managed by `dns-updater`.

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
