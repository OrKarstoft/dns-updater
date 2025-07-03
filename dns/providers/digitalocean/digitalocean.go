package digitalocean

import (
	"context"
	"fmt"
	"log"

	"github.com/digitalocean/godo"
	domain "github.com/orkarstoft/dns-updater"
	"github.com/orkarstoft/dns-updater/dns"
	"github.com/rs/zerolog"
)

type Service struct {
	client *godo.Client
	logger *zerolog.Logger
}

func NewService(logger *zerolog.Logger, apiToken string) dns.DNSImpl {
	if apiToken == "" {
		log.Fatal("API token is required")
	}

	client := godo.NewFromToken(apiToken)

	loggerSvc := logger.With().Str("module", "provider.digitalocean").Logger()

	return &Service{
		client: client,
		logger: &loggerSvc,
	}
}

func (s *Service) UpdateRecord(ctx context.Context, req *domain.DNSRequest) error {
	records, err := s.getRecords(ctx, req)
	if err != nil {
		return err
	}

	record := findMatchingRecord(records, req)
	if record == nil {
		return fmt.Errorf("record %s not found in domain %s", req.GetRecordName(), req.GetDomain())
	}

	if record.Data == req.GetIP() {
		s.logger.Debug().Msg("Record already up to date")
		return nil
	}

	if err := s.updateDNSRecord(ctx, req, record.ID); err != nil {
		return err
	}

	s.logger.Debug().Msg("Record updated")
	return nil
}

func (s *Service) getRecords(ctx context.Context, req *domain.DNSRequest) ([]godo.DomainRecord, error) {
	records, _, err := s.client.Domains.Records(ctx, req.GetDomain(), &godo.ListOptions{WithProjects: true})
	if err != nil {
		return nil, &dns.DNSProviderError{
			Provider:  "DigitalOcean",
			Operation: "list records",
			Err:       err,
		}
	}
	return records, nil
}

func findMatchingRecord(records []godo.DomainRecord, req *domain.DNSRequest) *godo.DomainRecord {
	for _, record := range records {
		if record.Name == req.GetRecordName() && record.Type == req.GetRecordType() {
			return &record
		}
	}
	return nil
}

func (s *Service) updateDNSRecord(ctx context.Context, req *domain.DNSRequest, recordID int) error {
	drer := &godo.DomainRecordEditRequest{
		Data: req.GetIP(),
	}

	_, _, err := s.client.Domains.EditRecord(ctx, req.GetDomain(), recordID, drer)
	if err != nil {
		return &dns.DNSProviderError{
			Provider:  "DigitalOcean",
			Operation: "edit record",
			Err:       err,
		}
	}
	return nil
}
