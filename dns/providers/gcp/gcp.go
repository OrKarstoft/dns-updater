package gcp

import (
	"context"
	"fmt"
	"log"

	domain "github.com/orkarstoft/dns-updater"
	"github.com/orkarstoft/dns-updater/dns"
	googledns "google.golang.org/api/dns/v1"
	"google.golang.org/api/option"
)

type Service struct {
	client    *googledns.Service
	projectId string
}

func NewService(projectId, credentialsFile string) dns.DNSImpl {
	ctx := context.TODO()
	client, err := googledns.NewService(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		log.Fatalf("Failed to create DNS client: %v", err)
	}
	return &Service{
		client:    client,
		projectId: projectId,
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
		return fmt.Errorf("record %s of type %s not found in zone %s", req.GetRecordName(), req.GetRecordType(), req.GetZone())
	}

	if recordToUpdate.Rrdatas[0] == req.GetIP() {
		fmt.Println("[DEBUG] Record is already up to date")
		return nil
	}

	if err := s.updateDNSRecord(s.projectId, req.GetZone(), recordToUpdate, req, fullRecordName); err != nil {
		return err
	}

	fmt.Println("[DEBUG] Change applied successfully")
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
			Provider:  "DigitalOcean",
			Operation: "record update",
			Err:       err,
		}
	}
	return nil
}
