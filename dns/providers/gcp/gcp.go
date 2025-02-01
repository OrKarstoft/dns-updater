package gcp

import (
	"context"
	"fmt"
	"log"

	domain "github.com/orkarstoft/dns-updater"
	"github.com/orkarstoft/dns-updater/config"
	"github.com/orkarstoft/dns-updater/dns"
	googledns "google.golang.org/api/dns/v1"
	"google.golang.org/api/option"
)

type Service struct {
	client *googledns.Service
}

func NewService() dns.DNSImpl {
	ctx := context.TODO()
	client, err := googledns.NewService(ctx, option.WithCredentialsFile(config.Conf.GetProviderString("CredentialsFilePath")))
	if err != nil {
		log.Fatalf("Failed to create DNS client: %v", err)
	}
	return &Service{
		client: client,
	}
}

func (s *Service) UpdateRecord(ctx context.Context, req *domain.DNSRequest) error {
	projectID := config.Conf.GetProviderString("ProjectId")
	fullRecordName := fmt.Sprintf("%s.%s.", req.GetRecordName(), req.GetDomain())

	records, err := s.listRecords(projectID, req.GetZone())
	if err != nil {
		return fmt.Errorf("failed to list resource record sets: %w", err)
	}

	recordToUpdate := findMatchingRecord(records, fullRecordName, req.GetRecordType())
	if recordToUpdate == nil {
		return fmt.Errorf("record %s of type %s not found in zone %s", req.GetRecordName(), req.GetRecordType(), req.GetZone())
	}

	if recordToUpdate.Rrdatas[0] == req.GetIP() {
		fmt.Println("[DEBUG] Record is already up to date")
		return nil
	}

	if err := s.updateDNSRecord(projectID, req.GetZone(), recordToUpdate, req); err != nil {
		return fmt.Errorf("failed to create change: %w", err)
	}

	fmt.Println("[DEBUG] Change applied successfully")
	return nil
}

func (s *Service) listRecords(projectID, zone string) ([]*googledns.ResourceRecordSet, error) {
	resp, err := s.client.ResourceRecordSets.List(projectID, zone).Do()
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

func (s *Service) updateDNSRecord(projectID, zone string, oldRecord *googledns.ResourceRecordSet, req *domain.DNSRequest) error {
	newRecord := &googledns.ResourceRecordSet{
		Name:    oldRecord.Name,
		Type:    oldRecord.Type,
		Ttl:     oldRecord.Ttl, // Preserve TTL
		Rrdatas: []string{req.GetIP()},
	}

	change := &googledns.Change{
		Additions: []*googledns.ResourceRecordSet{newRecord},
		Deletions: []*googledns.ResourceRecordSet{oldRecord},
	}

	_, err := s.client.Changes.Create(projectID, zone, change).Do()
	return err
}
