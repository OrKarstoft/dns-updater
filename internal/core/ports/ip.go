package ports

import (
	"context"
	"net/netip"
)

type IPFetcher interface {
	Get(ctx context.Context) (*netip.Addr, error)
}
