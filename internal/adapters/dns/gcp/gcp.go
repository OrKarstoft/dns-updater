package gcp

import (
	"context"
	"fmt"
	"net/netip"
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
		log.Fatal().Msgf("Failed to create DNS client: %v", err)
	}

	return &Provider{
		client:    client,
		projectId: cfg.GetString("projectId"),
	}, nil
}

func (p *Provider) GetRecords(ctx context.Context, zone, domain string) ([]ports.DNSRecord, error) {
	resp, err := p.client.ResourceRecordSets.List(p.projectId, zone).Do()
	if err != nil {
		return nil, err
	}
	return toDNSRecords(resp.Rrsets), nil
}

func (p *Provider) CreateRecord(ctx context.Context, zone, domain string, record ports.DNSRecord) (ports.DNSRecord, error) {
	newRecord := &googledns.ResourceRecordSet{
		Name:    record.Name,
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
	return toDNSRecord(*newRecord), nil
}

func (p *Provider) UpdateRecord(ctx context.Context, zone, domain, recordID string, record ports.DNSRecord) error {
	records, err := p.GetRecords(ctx, zone, domain)
	if err != nil {
		return err
	}

	oldRecord := findMatchingRecord(records, recordID, "")
	if oldRecord == nil {
		return fmt.Errorf("could not find record with id %s to update", recordID)
	}

	godoOldRecord := toResourceRecordSet(*oldRecord)

	newRecord := &googledns.ResourceRecordSet{
		Name:    oldRecord.Name,
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

	recordToDelete := findMatchingRecord(records, recordID, "")
	if recordToDelete == nil {
		return fmt.Errorf("could not find record with id %s to delete", recordID)
	}
	godoRecordToDelete := toResourceRecordSet(*recordToDelete)

	change := &googledns.Change{
		Deletions: []*googledns.ResourceRecordSet{&godoRecordToDelete},
	}

	_, err = p.client.Changes.Create(p.projectId, zone, change).Do()
	return err
}





func toDNSRecord(r googledns.ResourceRecordSet) ports.DNSRecord {
	return ports.DNSRecord{
		ID:   fmt.Sprintf("%s:%s", r.Name, r.Type),
		Name: r.Name,
		Type: r.Type,
		Data: r.Rrdatas[0],
		TTL:  int(r.Ttl),
	}
}

func toDNSRecords(rs []*googledns.ResourceRecordSet) []ports.DNSRecord {
	records := make([]ports.DNSRecord, 0, len(rs))
	for _, r := range rs {
		if len(r.Rrdatas) > 0 {
			records = append(records, toDNSRecord(*r))
		}
	}
	return records
}

func toResourceRecordSet(r ports.DNSRecord) googledns.ResourceRecordSet {
	return googledns.ResourceRecordSet{
		Name:    r.Name,
		Type:    r.Type,
		Rrdatas: []string{r.Data},
		Ttl:     int64(r.TTL),
	}
}
