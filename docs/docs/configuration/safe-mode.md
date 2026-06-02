---
title: Safe Mode
sidebar_position: 5
---

# Safe Mode

Safe mode prevents accidental modification or overwrite of DNS records that are not managed by `dns-updater`.

When safe mode is enabled, each managed A record has a companion TXT record. `dns-updater` uses this TXT record to verify ownership before it creates or updates the A record.

## Configuration

Configure safe mode under `provider.safemode`:

```yaml
provider:
  safemode:
    enabled: true
    txtOwnerId: dns-updater
    txtPrefix: dns-updater-safemode
```

## Fields

- `enabled` (boolean, default: `true`): enables or disables safe mode.
- `txtOwnerId` (string, default: `dns-updater`): ownership identifier written to the safemode TXT record.
- `txtPrefix` (string, default: `dns-updater-safemode`): prefix used to build the safemode TXT record name.

## Safemode TXT records

For each managed record, `dns-updater` creates or requires a companion TXT record.

- TXT record name: `<txt_prefix><record_name>`
- TXT record data: `managed-by:dns-updater/<txt_owner_id>`

With the default configuration, a record named `home` uses this companion TXT record:

```text
Name: dns-updater-safemodehome
Data: managed-by:dns-updater/dns-updater
```

## Create and update behavior

When safe mode is enabled, `dns-updater` uses the companion TXT record to decide whether it may create or update an A record.

- For each managed record, `dns-updater` creates or requires a companion TXT record with the expected ownership data.
- If the TXT record exists and its data matches the expected ownership value, `dns-updater` may create or update the A record.
- If an existing A record has no matching TXT record, or the TXT record data does not match, `dns-updater` refuses to touch the A record.

This prevents `dns-updater` from modifying records that it does not own.

## Example configuration

```yaml
provider:
  name: digitalocean
  config:
    token: your-digitalocean-api-token
  safemode:
    enabled: true
    txtOwnerId: dns-updater
    txtPrefix: dns-updater-safemode

updates:
  - domain: example.com
    zone: example-com
    type: A
    records:
      - home
      - vpn
```
