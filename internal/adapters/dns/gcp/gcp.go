package gcp

import (
	"context"
	"fmt"
	"net/netip"

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

func (p *Provider) UpdateRecord(ctx context.Context, zone, domain, recordName, recordType string, ip *netip.Addr) error {
	fullRecordName := fmt.Sprintf("%s.%s.", domain, zone)

	// If the record name is @, it means the root domain
	if domain == "@" {
		fullRecordName = fmt.Sprintf("%s.", domain)
	}

	records, err := p.listRecords(p.projectId, zone)
	if err != nil {
		return err
	}

	recordToUpdate := findMatchingRecord(records, fullRecordName, recordType)
	if recordToUpdate == nil {
		log.Debug().Msgf("Record %s not found in zone %s, creating new record", fullRecordName, zone)
		if err := p.createDNSRecord(p.projectId, zone, recordType, ip, fullRecordName); err != nil {
			return err
		}
		return nil
	}

	if recordToUpdate.Rrdatas[0] == ip.String() {
		log.Debug().Msg("Record is already up to date")
		return nil
	}

	if err := p.updateDNSRecord(p.projectId, zone, recordToUpdate, ip, fullRecordName); err != nil {
		return err
	}

	log.Debug().Msg("Record updated")
	return nil
}

func (p *Provider) listRecords(projectID, zone string) ([]*googledns.ResourceRecordSet, error) {
	resp, err := p.client.ResourceRecordSets.List(projectID, zone).Do()
	if err != nil {
		return nil, err
	}
	return resp.Rrsets, nil
}

func findMatchingRecord(records []*googledns.ResourceRecordSet, name, recordType string) *googledns.ResourceRecordSet {
	for _, record := range records {
		if record.Name == name && record.Type == recordType {
			return record
		}
	}
	return nil
}

func (p *Provider) updateDNSRecord(projectID, zone string, oldRecord *googledns.ResourceRecordSet, ip *netip.Addr, fullRecordName string) error {
	newRecord := &googledns.ResourceRecordSet{
		Name:    fullRecordName,
		Type:    oldRecord.Type,
		Ttl:     oldRecord.Ttl, // Preserve TTL
		Rrdatas: []string{ip.String()},
	}

	change := &googledns.Change{
		Additions: []*googledns.ResourceRecordSet{newRecord},
		Deletions: []*googledns.ResourceRecordSet{oldRecord},
	}

	_, err := p.client.Changes.Create(projectID, zone, change).Do()
	if err != nil {
		return err
	}
	return nil
}

func (p *Provider) createDNSRecord(projectID, zone, recordType string, ip *netip.Addr, fullRecordName string) error {
	newRecord := &googledns.ResourceRecordSet{
		Name:    fullRecordName,
		Type:    recordType,
		Ttl:     300, // Default TTL
		Rrdatas: []string{ip.String()},
	}
	_, err := p.client.ResourceRecordSets.Create(projectID, zone, newRecord).Do()
	if err != nil {
		return err
	}
	log.Debug().
		Str("name", newRecord.Name).
		Str("type", newRecord.Type).
		Str("ip", ip.String()).
		Str("zone", zone).
		Msg("Record created")
	return nil
}
