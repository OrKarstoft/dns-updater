package registry

import (
	"fmt"

	"github.com/orkarstoft/dns-updater/internal/config"
	"github.com/orkarstoft/dns-updater/internal/core/ports"
)

type DNSFactory func(cfg config.Provider) (ports.DNSProvider, error)

var dnsProviders = make(map[string]DNSFactory)

func RegisterDNSProvider(name string, factory DNSFactory) {
	if _, exists := dnsProviders[name]; exists {
		panic("DNS provider already registered: " + name)
	}
	dnsProviders[name] = factory
}

func GetDNSProvider(cfg config.Provider) (ports.DNSProvider, error) {
	factory, exists := dnsProviders[cfg.GetString("name")]
	if !exists {
		return nil, fmt.Errorf("DNS provider not found: %s", cfg.GetString("name"))
	}
	return factory(cfg)
}
