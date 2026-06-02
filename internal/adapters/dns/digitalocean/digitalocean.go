package digitalocean

import (
	"context"
	"fmt"
	"strconv"

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
		return nil, fmt.Errorf("digitalocean token is required")
	}

	client := godo.NewFromToken(token)

	return &Provider{
		apiToken: token,
		client:   client,
	}, nil
}

func (p *Provider) GetRecords(ctx context.Context, zone, domain string) ([]ports.DNSRecord, error) {
	records, _, err := p.client.Domains.Records(ctx, domain, &godo.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list domain records: %w", err)
	}
	return toDNSRecords(records), nil
}

func (p *Provider) CreateRecord(ctx context.Context, zone, domain string, record ports.DNSRecord) (ports.DNSRecord, error) {
	drr := &godo.DomainRecordEditRequest{
		Type: record.Type,
		Name: record.Name,
		Data: record.Data,
		TTL:  record.TTL,
	}

	createdRecord, _, err := p.client.Domains.CreateRecord(ctx, domain, drr)
	if err != nil {
		return ports.DNSRecord{}, err
	}

	log.Debug().
		Str("name", drr.Name).
		Str("type", drr.Type).
		Str("data", drr.Data).
		Str("zone", zone).
		Msg("Record created")
	return toDNSRecord(*createdRecord), nil
}

func (p *Provider) UpdateRecord(ctx context.Context, zone, domain, recordID string, record ports.DNSRecord) error {
	id, err := strconv.Atoi(recordID)
	if err != nil {
		return fmt.Errorf("invalid record ID: %s", recordID)
	}
	drer := &godo.DomainRecordEditRequest{
		Data: record.Data,
	}

	_, _, err = p.client.Domains.EditRecord(ctx, domain, id, drer)
	return err
}

func (p *Provider) DeleteRecord(ctx context.Context, zone, domain, recordID string) error {
	id, err := strconv.Atoi(recordID)
	if err != nil {
		return fmt.Errorf("invalid record ID: %s", recordID)
	}
	_, err = p.client.Domains.DeleteRecord(ctx, domain, id)
	return err
}

func toDNSRecord(r godo.DomainRecord) ports.DNSRecord {
	return ports.DNSRecord{
		ID:   strconv.Itoa(r.ID),
		Name: r.Name,
		Type: r.Type,
		Data: r.Data,
		TTL:  r.TTL,
	}
}

func toDNSRecords(rs []godo.DomainRecord) []ports.DNSRecord {
	records := make([]ports.DNSRecord, len(rs))
	for i, r := range rs {
		records[i] = toDNSRecord(r)
	}
	return records
}
