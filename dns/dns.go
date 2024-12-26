package dns

import domain "github.com/orkarstoft/dns-updater"

type DNSImpl interface {
	UpdateRecord(*domain.DNSRequest)
}
