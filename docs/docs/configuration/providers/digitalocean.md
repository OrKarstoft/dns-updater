---
title: DigitalOcean
---

# Provider: DigitalOcean DNS

To use DigitalOcean DNS, set:

- `provider.name: digitalocean`
- `provider.config.token`: a DigitalOcean API token (required)

If `token` is missing/empty, `dns-updater` terminates with a fatal error.

## Example

```yaml
provider:
  name: digitalocean
  config:
    token: "<DO_TOKEN>"

updates:
  - domain: example.com
    zone: example-com
    type: A
    records:
      - "@"
      - "home"
```

## Notes

- Keep tokens out of version control.
- Use your platform’s secret management (CI secrets, Docker secrets, etc.) to inject the token securely.
