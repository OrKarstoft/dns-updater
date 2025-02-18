package factory

import (
	"fmt"

	"github.com/orkarstoft/dns-updater/dns"
	"github.com/orkarstoft/dns-updater/dns/providers/digitalocean"
	"github.com/orkarstoft/dns-updater/dns/providers/gcp"
)

// CreateProvider creates a new DNS provider based on the configuration
func CreateProvider(config ProviderConfig) (dns.DNSImpl, error) {
	switch config.Type {
	case ProviderGCP:
		return gcp.NewService(), nil
	case ProviderDigitalOcean:
		token, ok := config.Config["token"].(string)
		if !ok {
			return nil, fmt.Errorf("digitalocean provider requires a token")
		}
		return digitalocean.NewService(token), nil
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", config.Type)
	}
}
