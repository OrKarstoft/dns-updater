package gcp

import (
	"context"
	"fmt"
	"log"

	domain "github.com/orkarstoft/dns-updater"
	"github.com/orkarstoft/dns-updater/dns"
	"github.com/rs/zerolog"
	googledns "google.golang.org/api/dns/v1"
	"google.golang.org/api/option"
)

const (
	providerName = "GoogleCloudPlatform"
)

type Service struct {
	client    *googledns.Service
	projectId string
	logger    *zerolog.Logger
}

func NewService(logger *zerolog.Logger, projectId, credentialsFile string) dns.DNSImpl {
	ctx := context.TODO()
	client, err := googledns.NewService(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		log.Fatalf("Failed to create DNS client: %v", err)
	}

	loggerSvc := logger.With().Str("module", "provider.gcp").Logger()
	return &Service{
		client:    client,
		projectId: projectId,
		logger:    &loggerSvc,
	}
}

func (s *Service) UpdateRecord(ctx context.Context, req *domain.DNSRequest) error {
	fullRecordName := fmt.Sprintf("%s.%s.", req.GetRecordName(), req.GetDomain())

	// If the record name is @, it means the root domain
	if req.GetRecordName() == "@" {
		fullRecordName = fmt.Sprintf("%s.", req.GetDomain())
	}

	records, err := s.listRecords(s.projectId, req.GetZone())
	if err != nil {
		return err
	}

	recordToUpdate := findMatchingRecord(records, fullRecordName, req.GetRecordType())
	if recordToUpdate == nil {
		s.logger.Debug().Msgf("Record %s not found in zone %s, creating new record", fullRecordName, req.GetZone())
		if err := s.createDNSRecord(s.projectId, req.GetZone(), req, fullRecordName); err != nil {
			return fmt.Errorf("failed to create record: %w", err)
		}
		return nil
	}

	if recordToUpdate.Rrdatas[0] == req.GetIP() {
		s.logger.Debug().Msg("Record is already up to date")
		return nil
	}

	if err := s.updateDNSRecord(s.projectId, req.GetZone(), recordToUpdate, req, fullRecordName); err != nil {
		return err
	}

	s.logger.Debug().Msg("Record updated")
	return nil
}

func (s *Service) listRecords(projectID, zone string) ([]*googledns.ResourceRecordSet, error) {
	resp, err := s.client.ResourceRecordSets.List(projectID, zone).Do()
	if err != nil {
		return nil, &dns.DNSProviderError{
			Provider:  "GoogleCloudPlatform",
			Operation: "list records",
			Err:       err,
		}
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

func (s *Service) updateDNSRecord(projectID, zone string, oldRecord *googledns.ResourceRecordSet, req *domain.DNSRequest, fullRecordName string) error {
	newRecord := &googledns.ResourceRecordSet{
		Name:    fullRecordName,
		Type:    oldRecord.Type,
		Ttl:     oldRecord.Ttl, // Preserve TTL
		Rrdatas: []string{req.GetIP()},
	}

	change := &googledns.Change{
		Additions: []*googledns.ResourceRecordSet{newRecord},
		Deletions: []*googledns.ResourceRecordSet{oldRecord},
	}

	_, err := s.client.Changes.Create(projectID, zone, change).Do()
	if err != nil {
		return &dns.DNSProviderError{
			Provider:  providerName,
			Operation: dns.OperationUpdateRecord,
			Err:       err,
		}
	}
	return nil
}

func (s *Service) createDNSRecord(projectID, zone string, req *domain.DNSRequest, fullRecordName string) error {
	newRecord := &googledns.ResourceRecordSet{
		Name:    fullRecordName,
		Type:    req.GetRecordType(),
		Ttl:     300, // Default TTL
		Rrdatas: []string{req.GetIP()},
	}
	_, err := s.client.ResourceRecordSets.Create(projectID, zone, newRecord).Do()
	if err != nil {
		return &dns.DNSProviderError{
			Provider:  providerName,
			Operation: dns.OperationCreateRecord,
			Err:       err,
		}
	}
	s.logger.Debug().Msg("Record created")
	return nil
}
