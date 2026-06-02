package gcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/orkarstoft/dns-updater/internal/config"
	"github.com/orkarstoft/dns-updater/internal/core/ports"
	"github.com/orkarstoft/dns-updater/internal/registry"
	"github.com/rs/zerolog/log"
	googledns "google.golang.org/api/dns/v1"
	"google.golang.org/api/option"
)

func init() {
	registry.RegisterDNSProvider("googlecloudplatform", NewFromConfig)
}

type Provider struct {
	client    *googledns.Service
	projectId string
}

func NewFromConfig(cfg config.Provider) (ports.DNSProvider, error) {
	ctx := context.Background()
	client, err := googledns.NewService(ctx, option.WithCredentialsFile(cfg.GetString("credentialsFile")))
	if err != nil {
		return nil, fmt.Errorf("failed to create GCP DNS client: %w", err)
	}

	projectID := cfg.GetString("projectId")
	if projectID == "" {
		return nil, fmt.Errorf("gcp projectId is required")
	}

	return &Provider{
		client:    client,
		projectId: projectID,
	}, nil
}

func (p *Provider) GetRecords(ctx context.Context, zone, domain string) ([]ports.DNSRecord, error) {
	resp, err := p.client.ResourceRecordSets.List(p.projectId, zone).Do()
	if err != nil {
		return nil, err
	}
	return toDNSRecords(resp.Rrsets, domain), nil
}

func (p *Provider) CreateRecord(ctx context.Context, zone, domain string, record ports.DNSRecord) (ports.DNSRecord, error) {
	newRecord := &googledns.ResourceRecordSet{
		Name:    expandRecordName(record.Name, domain),
		Type:    record.Type,
		Ttl:     int64(record.TTL),
		Rrdatas: []string{record.Data},
	}
	_, err := p.client.ResourceRecordSets.Create(p.projectId, zone, newRecord).Do()
	if err != nil {
		return ports.DNSRecord{}, err
	}
	log.Debug().
		Str("name", newRecord.Name).
		Str("type", newRecord.Type).
		Str("data", record.Data).
		Str("zone", zone).
		Msg("Record created")
	return toDNSRecord(*newRecord, domain), nil
}

func (p *Provider) UpdateRecord(ctx context.Context, zone, domain, recordID string, record ports.DNSRecord) error {
	records, err := p.GetRecords(ctx, zone, domain)
	if err != nil {
		return err
	}

	oldRecord := findRecordByID(records, recordID)
	if oldRecord == nil {
		return fmt.Errorf("could not find record with id %s to update", recordID)
	}

	godoOldRecord := toResourceRecordSet(*oldRecord, domain)

	newRecord := &googledns.ResourceRecordSet{
		Name:    expandRecordName(oldRecord.Name, domain),
		Type:    oldRecord.Type,
		Ttl:     int64(record.TTL),
		Rrdatas: []string{record.Data},
	}

	change := &googledns.Change{
		Additions: []*googledns.ResourceRecordSet{newRecord},
		Deletions: []*googledns.ResourceRecordSet{&godoOldRecord},
	}

	_, err = p.client.Changes.Create(p.projectId, zone, change).Do()
	return err
}

func (p *Provider) DeleteRecord(ctx context.Context, zone, domain, recordID string) error {
	records, err := p.GetRecords(ctx, zone, domain)
	if err != nil {
		return err
	}

	recordToDelete := findRecordByID(records, recordID)
	if recordToDelete == nil {
		return fmt.Errorf("could not find record with id %s to delete", recordID)
	}
	godoRecordToDelete := toResourceRecordSet(*recordToDelete, domain)

	change := &googledns.Change{
		Deletions: []*googledns.ResourceRecordSet{&godoRecordToDelete},
	}

	_, err = p.client.Changes.Create(p.projectId, zone, change).Do()
	return err
}

func toDNSRecord(r googledns.ResourceRecordSet, domain string) ports.DNSRecord {
	return ports.DNSRecord{
		ID:   fmt.Sprintf("%s:%s", r.Name, r.Type),
		Name: normalizeRecordName(r.Name, domain),
		Type: r.Type,
		Data: r.Rrdatas[0],
		TTL:  int(r.Ttl),
	}
}

func toDNSRecords(rs []*googledns.ResourceRecordSet, domain string) []ports.DNSRecord {
	records := make([]ports.DNSRecord, 0, len(rs))
	for _, r := range rs {
		if len(r.Rrdatas) > 0 {
			records = append(records, toDNSRecord(*r, domain))
		}
	}
	return records
}

func toResourceRecordSet(r ports.DNSRecord, domain string) googledns.ResourceRecordSet {
	return googledns.ResourceRecordSet{
		Name:    expandRecordName(r.Name, domain),
		Type:    r.Type,
		Rrdatas: []string{r.Data},
		Ttl:     int64(r.TTL),
	}
}

func normalizeRecordName(name, domain string) string {
	trimmedName := strings.TrimSuffix(name, ".")
	trimmedDomain := strings.TrimSuffix(domain, ".")

	if trimmedDomain == "" {
		return trimmedName
	}

	if trimmedName == trimmedDomain {
		return "@"
	}

	suffix := "." + trimmedDomain
	if strings.HasSuffix(trimmedName, suffix) {
		return strings.TrimSuffix(trimmedName, suffix)
	}

	return trimmedName
}

func expandRecordName(name, domain string) string {
	trimmedName := strings.TrimSuffix(name, ".")
	trimmedDomain := strings.TrimSuffix(domain, ".")

	if trimmedDomain == "" {
		if trimmedName == "" {
			return "."
		}
		return trimmedName + "."
	}

	switch trimmedName {
	case "", "@", trimmedDomain:
		return trimmedDomain + "."
	}

	if strings.HasSuffix(trimmedName, "."+trimmedDomain) {
		return trimmedName + "."
	}

	return fmt.Sprintf("%s.%s.", trimmedName, trimmedDomain)
}

func findRecordByID(records []ports.DNSRecord, id string) *ports.DNSRecord {
	for i := range records {
		if records[i].ID == id {
			return &records[i]
		}
	}
	return nil
}
