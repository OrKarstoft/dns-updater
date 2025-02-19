# dns-updater

A lightweight DNS record updater that supports multiple DNS providers (DigitalOcean and Google Cloud Platform) to automatically update DNS records with your current public IP address.

## Features

- Supports multiple DNS providers:
  - DigitalOcean DNS
  - Google Cloud Platform DNS
- Automatic IP address detection
- Docker support with minimal secure image
- Configurable updates for multiple domains and records
- YAML-based configuration

## Installation

### Using Docker (Recommended)

Pull the latest version from GitHub Container Registry:
```bash
docker pull ghcr.io/orkarstoft/dns-updater:latest
```


### From source

1. Clone the repository:
```bash
git clone https://github.com/orkarstoft/dns-updater.git
cd dns-updater
```

2. Build the binary:
```bash
go build -o dnsupdater cmd/main.go
```

## Configuration

Create a `config.yaml` file with your DNS provider credentials and update configuration:

### DigitalOcean Example:

```bash
provider:
  name: digitalocean
  config:
    token: <DO_TOKEN>

updates:
  - domain: example.com
    zone: example-com
    type: A
    records:
    - "@" # Set the root level record, so example.com
    - record1 # Set the subdomain record, so record1.example.com
```

### Google Cloud Platform Example:

```bash
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

## Usage

### Running with Docker

```bash
docker run -v /path/to/config.yaml:/config.yaml ghcr.io/orkarstoft/dns-updater:latest
```

For GCP authentication, mount your credentials file:
```bash
docker run -v /path/to/config.yaml:/config.yaml \
          -v /path/to/credentials.json:/credentials.json \
          ghcr.io/orkarstoft/dns-updater:latest
```

### Running from binary

```bash
./dns-updater
```

## Development

### Prerequisites

- Go 1.23 or higher
- Docker (for container builds)

### Building

Build the binary:
```bash
go build -o dnsupdater cmd/main.go
```

Build the Docker image:
```bash
docker build -t dns-updater .
```

### Testing

```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create your feature branch (git checkout -b feature/amazing-feature)
3. Commit your changes (git commit -m 'Add some amazing feature')
4. Push to the branch (git push origin feature/amazing-feature)
5. Open a Pull Request

## Security

- The Docker image runs as a non-root user
- Credentials should be properly secured and never committed to version control
- For production use, store sensitive configuration in secure locations

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Uses the [godo](https://github.com/digitalocean/godo) client for DigitalOcean API 
- Uses the Google Cloud DNS API for GCP integration

