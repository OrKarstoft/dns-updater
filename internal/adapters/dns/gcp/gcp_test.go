package gcp

import (
	"testing"

	"github.com/orkarstoft/dns-updater/internal/core/ports"
	googledns "google.golang.org/api/dns/v1"
)

func TestToDNSRecordNormalizesRecordNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		record   googledns.ResourceRecordSet
		domain   string
		expected ports.DNSRecord
	}{
		{
			name:   "zone apex becomes at",
			domain: "example.com",
			record: googledns.ResourceRecordSet{
				Name:    "example.com.",
				Type:    "A",
				Rrdatas: []string{"1.2.3.4"},
				Ttl:     300,
			},
			expected: ports.DNSRecord{
				ID:   "example.com.:A",
				Name: "@",
				Type: "A",
				Data: "1.2.3.4",
				TTL:  300,
			},
		},
		{
			name:   "subdomain strips domain suffix",
			domain: "example.com",
			record: googledns.ResourceRecordSet{
				Name:    "home.example.com.",
				Type:    "A",
				Rrdatas: []string{"1.2.3.4"},
				Ttl:     600,
			},
			expected: ports.DNSRecord{
				ID:   "home.example.com.:A",
				Name: "home",
				Type: "A",
				Data: "1.2.3.4",
				TTL:  600,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := toDNSRecord(tt.record, tt.domain); got != tt.expected {
				t.Fatalf("toDNSRecord() = %#v, want %#v", got, tt.expected)
			}
		})
	}
}

func TestExpandRecordName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		record   string
		domain   string
		expected string
	}{
		{name: "zone apex", record: "@", domain: "example.com", expected: "example.com."},
		{name: "subdomain", record: "home", domain: "example.com", expected: "home.example.com."},
		{name: "already fqdn", record: "home.example.com.", domain: "example.com", expected: "home.example.com."},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := expandRecordName(tt.record, tt.domain); got != tt.expected {
				t.Fatalf("expandRecordName() = %q, want %q", got, tt.expected)
			}
		})
	}
}
