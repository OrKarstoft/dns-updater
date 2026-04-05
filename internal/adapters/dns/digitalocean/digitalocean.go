package digitalocean

import (
	"context"
	"fmt"
	"net/netip"

	"github.com/digitalocean/godo"
	"github.com/orkarstoft/dns-updater/internal/config"
	"github.com/orkarstoft/dns-updater/internal/core/ports"
	"github.com/orkarstoft/dns-updater/internal/registry"
	"github.com/rs/zerolog/log"
)

func init() {
	registry.RegisterDNSProvider("digitalocean", NewFromConfig)
}

type Provider struct {
	apiToken string
	client   *godo.Client
}

func NewFromConfig(cfg config.Provider) (ports.DNSProvider, error) {
	token := cfg.GetString("token")
	if token == "" {
		log.Fatal().Msg("DigitalOcean API token is required")
	}

	client := godo.NewFromToken(token)

	return &Provider{
		apiToken: token,
		client:   client,
	}, nil
}

func (s *Provider) UpdateRecord(ctx context.Context, zone, domain, recordName, recordType string, ip *netip.Addr) error {
	records, err := s.getRecords(ctx, domain)
	if err != nil {
		return err
	}
	log.Debug().Msgf("Found %d records in domain %s", len(records), domain)

	record := findMatchingRecord(records, recordType, recordName)
	if record == nil {
		log.Debug().Msgf("Record %s not found in zone %s, creating new record", fmt.Sprintf("%s.%s", recordName, domain), zone)
		if err := s.createDNSRecord(ctx, recordType, recordName, domain, zone, ip); err != nil {
			return err
		}
		return nil
	}

	log.Debug().Msgf("Record %s found in zone %s, updating record", fmt.Sprintf("%s.%s", recordName, domain), zone)
	if record.Data == ip.String() {
		log.Debug().Msg("Record already up to date")
		return nil
	}

	if err := s.updateDNSRecord(ctx, domain, ip, record.ID); err != nil {
		return err
	}

	log.Debug().Msg("Record updated")

	return nil
}

func (s *Provider) getRecords(ctx context.Context, domain string) ([]godo.DomainRecord, error) {
	records, _, err := s.client.Domains.Records(ctx, domain, &godo.ListOptions{WithProjects: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list domain records: %w", err)
	}
	return records, nil
}

func findMatchingRecord(records []godo.DomainRecord, recordType, recordName string) *godo.DomainRecord {
	for _, record := range records {
		if record.Name == recordName && record.Type == recordType {
			return &record
		}
	}
	return nil
}

func (p *Provider) updateDNSRecord(ctx context.Context, domain string, ip *netip.Addr, recordID int) error {
	drer := &godo.DomainRecordEditRequest{
		Data: ip.String(),
	}

	_, _, err := p.client.Domains.EditRecord(ctx, domain, recordID, drer)
	if err != nil {
		return err
	}
	return nil
}

func (p *Provider) createDNSRecord(ctx context.Context, recordType, recordName, domain, zone string, ip *netip.Addr) error {
	drr := &godo.DomainRecordEditRequest{
		Type: recordType,
		Name: recordName,
		Data: ip.String(),
	}

	_, _, err := p.client.Domains.CreateRecord(ctx, domain, drr)
	if err != nil {
		return err
	}

	log.Debug().
		Str("name", drr.Name).
		Str("type", drr.Type).
		Str("ip", ip.String()).
		Str("zone", zone).
		Msg("Record created")
	return nil
}
