# dns-updater

A lightweight DNS record updater that supports multiple DNS providers (DigitalOcean, Simply.com, and Google Cloud Platform) to automatically update DNS records with your current public IP address.

## Documentation

- Hosted documentation: https://orkarstoft.github.io/dns-updater
- Docs source (Docusaurus): [`/docs`](./docs)

## Features

- Supports multiple DNS providers:
  - DigitalOcean DNS
  - Simply.com DNS
  - Google Cloud DNS (GCP)
- Automatic public IP address detection
- Optional caching to avoid unnecessary DNS updates (file cache)
- YAML-based configuration (`config.yaml`)
- Docker support with a minimal, secure image
- Configurable logging:
  - `pretty`, `json`, or `file`

## Installation

### Using Docker (Recommended)

Pull the latest version from GitHub Container Registry:

```bash
docker pull ghcr.io/orkarstoft/dns-updater:latest
```

### From source

```bash
git clone https://github.com/OrKarstoft/dns-updater.git
cd dns-updater
go build -o dnsupdater cmd/main.go
```

## Configuration

`dns-updater` expects a `config.yaml` in the **current working directory** (it looks for `config.yaml` via Viper using config name `config`).

### Provider configuration

Providers are selected via:

```yaml
provider:
  name: <providerName>
  config:
    # provider-specific keys here
```

Supported providers available:

| Provider               | `provider.name`       | State                                      |
| ---------------------- | --------------------- | ------------------------------------------ |
| DigitalOcean           | `digitalocean`        | ✅ Tested                                  |
| Simply.com             | `simply`              | ✅ Tested                                  |
| Google Cloud DNS (GCP) | `googlecloudplatform` | ⚠️ Untested (untested since release 3.0.0) |

### Update configuration

Updates are configured under `updates` as a list. Each item includes:

- `domain`: domain name (e.g. `example.com`)
- `zone`: provider zone identifier (provider-specific)
- `type`: record type (e.g. `A`)
- `records`: list of record names (e.g. `"@"`, `www`, etc.)

#### DigitalOcean example

```yaml
provider:
  name: digitalocean
  config:
    token: <DO_TOKEN>

updates:
  - domain: example.com
    zone: example-com
    type: A
    records:
      - "@" # Root record => example.com
      - record1 # Subdomain => record1.example.com
```

#### Simply.com example

```yaml
provider:
  name: simply
  config:
    accountName: "YOUR_ACCOUNT_NAME"
    apiKey: "YOUR_API_KEY"

updates:
  - domain: example.com
    zone: example-com
    type: "A"
    records:
      - record1
      - record2
```

#### Google Cloud DNS example

```yaml
provider:
  name: googlecloudplatform
  config:
    credentialsFile: "/path/to/credentials.json"
    projectId: "your-project-id"

updates:
  - domain: example.com
    zone: example-com
    type: "A"
    records:
      - record1
      - record2
```

### Safe Mode

Safe mode is enabled by default and helps prevent accidental modification of records not managed by `dns-updater`. For each managed record, `dns-updater` creates/requires a companion TXT record named `<txt_prefix><record_name>` whose data is exactly `managed-by:dns-updater/<txt_owner_id>`.

When enabled, `dns-updater` only creates or updates an A record if the matching safe mode TXT record exists and matches the expected ownership data. If a record exists without a matching safe mode TXT record, it refuses to touch it.

```yaml
provider:
  name: digitalocean
  safemode:
    enabled: true
    txtOwnerId: "dns-updater"
    txtPrefix: "dns-updater-safemode"
  config:
    token: "<DO_TOKEN>"
```

Full details on safe mode and cleanup are available in the hosted docs: https://orkarstoft.github.io/dns-updater

### Cache configuration (optional)

This branch supports a simple file cache to persist the last observed IP and skip updates when it hasn’t changed:

```yaml
cache:
  enabled: true
  filePath: "dns-updater.cache"
```

If `cache.enabled` is `false`, a no-op cache is used.

### Logging configuration

Logging defaults to:

- `log.level: info`
- `log.type: pretty`

You can override with:

```yaml
log:
  level: debug # info, warning, debug (as supported by the logger)
  type: pretty # pretty, json, file
```

## Usage

### Running with Docker

```bash
docker run -v /path/to/config.yaml:/config.yaml ghcr.io/orkarstoft/dns-updater:latest
```

> Note: the application looks for `config.yaml` in the working directory. Ensure the container’s working directory and mount path align with how the image is built/run.

For GCP authentication, mount your credentials file:

```bash
docker run -v /path/to/config.yaml:/config.yaml \
          -v /path/to/credentials.json:/credentials.json \
          ghcr.io/orkarstoft/dns-updater:latest
```

### Running from binary

Place `config.yaml` next to the binary (or run from the directory containing it):

```bash
./dnsupdater
```

### Cleanup

Use `--clean` to delete the managed records and their associated safe mode TXT records for every record in your config. Cleanup is guarded: it only deletes a record when the matching safe mode TXT record with the correct ownership data is present, so it won’t delete records it doesn’t own.

```bash
./dnsupdater --clean
```

Docker equivalent:

```bash
docker run -v /path/to/config.yaml:/config.yaml ghcr.io/orkarstoft/dns-updater:latest --clean
```

Full details on safe mode and cleanup are available in the hosted docs: https://orkarstoft.github.io/dns-updater

## Development

### Prerequisites

- Go (per `go.mod`)
- Docker (for container builds)
- Node.js >= 20 (only needed for `/docs`)

### Testing

```bash
go test ./...
```

## Contributing

See [`CONTRIBUTING.md`](./CONTRIBUTING.md).

## Security

- The Docker image runs as a non-root user (per project intent)
- Never commit provider credentials (tokens, service account JSON files)

## License

This project is licensed under the MIT License - see [`LICENSE.md`](./LICENSE.md).

## Acknowledgments

- Uses the [godo](https://github.com/digitalocean)
