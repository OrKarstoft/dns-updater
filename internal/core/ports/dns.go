package ports

import (
	"context"
	"net/netip"
)

type DNSProvider interface {
	UpdateRecord(ctx context.Context, zone, domain, recordName, recordType string, ip *netip.Addr) error
}
