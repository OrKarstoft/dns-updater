package dns

import (
	"context"

	domain "github.com/orkarstoft/dns-updater"
)

type DNSImpl interface {
	UpdateRecord(context.Context, *domain.DNSRequest) error
}
