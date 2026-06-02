---
title: Google Cloud DNS (GCP)
---

> [!WARNING]
> Google Cloud DNS provider is untested since release 3.0.0. Use with caution and report any issues you encounter.

# Provider: Google Cloud DNS (GCP)

To use Google Cloud DNS, set:

- `provider.name: googlecloudplatform`
- `provider.config.credentialsFile`: path to a Google service account JSON file
- `provider.config.projectId`: your GCP project ID

## Example

```yaml
provider:
  name: googlecloudplatform
  config:
    credentialsFile: "/path/to/credentials.json"
    projectId: "your-gcp-project-id"

updates:
  - domain: example.com
    zone: example-com
    type: A
    records:
      - "@"
      - "home"
```

## Notes

- The credentials file must be readable by the process/container.
- Record creation and updates are governed by the `provider.safemode` setting. See [Safe Mode](../safe-mode.md) for more details.
