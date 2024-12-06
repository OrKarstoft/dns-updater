package dns

import domain "github.com/orkarstoft/dns-updater"

type DNSImpl interface {
	SetRecord(*domain.DNSRequest)
}
