package cache

import (
	"context"
	"net/netip"
)

// NoOpCache implements the ports.IPCache interface but does nothing.
// It is used when the user disables caching in their configuration.
type NoOpCache struct{}

// NewNoOpCache creates a new NoOpCache instance.
func NewNoOpCache() *NoOpCache {
	return &NoOpCache{}
}

// GetLastIP always returns a nil IP address.
// This guarantees that in updater.go, `currentIP == lastIP` will always evaluate to false,
// forcing a DNS update every time the schedule runs.
func (c *NoOpCache) GetLastIP(ctx context.Context) (*netip.Addr, error) {
	return nil, nil
}

// SetLastIP immediately returns success without saving anything.
func (c *NoOpCache) SetLastIP(ctx context.Context, ip *netip.Addr) error {
	return nil
}
