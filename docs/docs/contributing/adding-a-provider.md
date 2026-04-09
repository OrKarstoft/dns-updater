---
sidebar_position: 2
---

# Adding a new DNS Provider

This guide will walk you through the process of adding a new DNS provider to this project.

## Introduction

A DNS provider is a service that manages your DNS records. This project uses a common interface to interact with different DNS providers, allowing users to choose the provider that best suits their needs. To add a new provider, you'll need to implement this interface and register your provider with the application.

## Provider Interface

The core of the provider integration is the `ports.DNSProvider` interface. Any new provider must implement this interface.

```go
// internal/core/ports/dns.go

// DNSRecord represents a generic DNS record.
type DNSRecord struct {
	ID   string
	Name string // Name of the record, e.g., 'www'
	Type string // Type of the record, e.g., 'A', 'TXT'
	Data string // Data for the record, e.g., an IP address or text
	TTL  int
}

// DNSProvider is the interface for DNS providers.
type DNSProvider interface {
	// GetRecords retrieves all records for a given domain.
	GetRecords(ctx context.Context, zone, domain string) ([]DNSRecord, error)
	// CreateRecord creates a new DNS record.
	CreateRecord(ctx context.Context, zone, domain string, record DNSRecord) (DNSRecord, error)
	// UpdateRecord updates an existing DNS record.
	UpdateRecord(ctx context.Context, zone, domain, recordID string, record DNSRecord) error
	// DeleteRecord deletes an existing DNS record.
	DeleteRecord(ctx context.Context, zone, domain, recordID string) error
}
```

## Implementation Steps

### 1. Create a new Provider Package

Create a new directory for your provider under `internal/adapters/dns/`. For example, if you're adding a provider named "cloudprovider", you would create the directory `internal/adapters/dns/cloudprovider/`.

Inside this directory, create a `.go` file for your implementation, e.g., `cloudprovider.go`.

### 2. Implement the `DNSProvider` Interface

In your `cloudprovider.go` file, define a struct for your provider that holds any necessary data, such as an API client or configuration values. Then, implement the methods of the `DNSProvider` interface.

```go
package cloudprovider

import (
	"context"
	"github.com/orkarstoft/dns-updater/internal/config"
	"github.com/orkarstoft/dns-updater/internal/core/ports"
)

type Provider struct {
	// Add provider-specific fields here, e.g., an API client.
}

func (p *Provider) GetRecords(ctx context.Context, zone, domain string) ([]ports.DNSRecord, error) {
	// Implementation to get DNS records
}

func (p *Provider) CreateRecord(ctx context.Context, zone, domain string, record ports.DNSRecord) (ports.DNSRecord, error) {
	// Implementation to create a DNS record
}

func (p *Provider) UpdateRecord(ctx context.Context, zone, domain, recordID string, record ports.DNSRecord) error {
	// Implementation to update a DNS record
}

func (p *Provider) DeleteRecord(ctx context.Context, zone, domain, recordID string) error {
	// Implementation to delete a DNS record
}

```

### 3. Create the `NewFromConfig` Function

This function is responsible for creating an instance of your provider and initializing it with the necessary configuration. It reads the provider-specific configuration from the `config.Provider` object.

```go
func NewFromConfig(cfg config.Provider) (ports.DNSProvider, error) {
	apiToken := cfg.GetString("api_token")
	if apiToken == "" {
		return nil, fmt.Errorf("CloudProvider API token is required")
	}

	// Initialize your provider's client here
	// client := cloudproviderclient.New(apiToken)

	return &Provider{
		// client: client,
	}, nil
}
```

### 4. Register the Provider

Use an `init()` function to register your new provider with the application's provider registry. This makes it available for use when specified in the configuration.

```go
// internal/adapters/dns/cloudprovider/cloudprovider.go

import (
    "github.com/orkarstoft/dns-updater/internal/registry"
)

func init() {
	registry.RegisterDNSProvider("cloudprovider", NewFromConfig)
}

```

The `registry.RegisterDNSProvider` function is defined in `internal/registry/dns.go`:

```go
// internal/registry/dns.go
var dnsProviders = make(map[string]DNSProviderFactory)

// DNSProviderFactory is a function that creates a DNSProvider.
type DNSProviderFactory func(cfg config.Provider) (ports.DNSProvider, error)

// RegisterDNSProvider registers a new DNS provider.
func RegisterDNSProvider(name string, factory DNSProviderFactory) {
	dnsProviders[name] = factory
}
```

### 5. Add Provider to Main Package

To ensure your provider's `init` function is called, you need to import it in the main application. Open `cmd/main.go` and add a blank import for your new provider package.

```go
// cmd/main.go

import (
    // ... other imports
    _ "github.com/orkarstoft/dns-updater/internal/adapters/dns/cloudprovider"
)
```

## Configuration

The application uses a `config.yaml` file for configuration. To use your new provider, the user needs to specify it in the `provider` section of this file. The `config` block within the `provider` section is where your provider-specific settings go.

Here's an example of how the configuration for your new "cloudprovider" would look:

```yaml
provider:
  name: cloudprovider
  safemode: true
  config:
    api_token: "your-api-token-here"
    # Add other provider-specific configuration keys here
```

Your `NewFromConfig` function will have access to these values through the `config.Provider` object passed to it.

## Data Models

Your provider will likely have its own data structures for representing DNS records. You'll need to map these to the project's common `ports.DNSRecord` struct. It's a good practice to create helper functions for this conversion.

Example from the DigitalOcean provider:

```go
// internal/adapters/dns/digitalocean/digitalocean.go

import (
    "github.com/digitalocean/godo"
    "github.com/orkarstoft/dns-updater/internal/core/ports"
    "strconv"
)

func toDNSRecord(r godo.DomainRecord) ports.DNSRecord {
	return ports.DNSRecord{
		ID:   strconv.Itoa(r.ID),
		Name: r.Name,
		Type: r.Type,
		Data: r.Data,
		TTL:  r.TTL,
	}
}

func toDNSRecords(rs []godo.DomainRecord) []ports.DNSRecord {
	records := make([]ports.DNSRecord, len(rs))
	for i, r := range rs {
		records[i] = toDNSRecord(r)
	}
	return records
}
```

By following these steps, you can successfully add a new DNS provider to the project, extending its capabilities and providing more options for users.
