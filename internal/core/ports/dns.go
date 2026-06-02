package ports

import (
	"context"
)

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
