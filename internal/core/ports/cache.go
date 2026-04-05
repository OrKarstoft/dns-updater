package ports

import (
	"context"
	"net/netip"
)

type IPCache interface {
	// GetLastIP retrieves the last known IP address from the cache.
	// It returns an error if the cache is empty or if there was an issue reading it.
	GetLastIP(ctx context.Context) (*netip.Addr, error)

	// SetLastIP updates the cache with the new IP address.
	// It returns an error if there was an issue writing to the cache.
	SetLastIP(ctx context.Context, ip *netip.Addr) error
}
