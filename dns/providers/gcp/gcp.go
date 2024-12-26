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
	ctx    context.Context
	client *googledns.Service
}

func NewService() dns.DNSImpl {
	ctx := context.TODO()
	client, err := googledns.NewService(ctx, option.WithCredentialsFile(config.Conf.GCP.CredentialsFilePath))
	if err != nil {
		log.Fatalf("Failed to create DNS client: %v", err)
	}
	return &Service{ctx: ctx, client: client}
}

func (s *Service) UpdateRecord(req *domain.DNSRequest) {
	fullRecordName := fmt.Sprintf("%s.%s.", req.GetRecordName(), req.GetDomain())

	// List existing records in the zone
	recordSets, err := s.client.ResourceRecordSets.List(config.Conf.GCP.ProjectId, req.GetZone()).Do()
	if err != nil {
		log.Fatalf("Failed to list resource record sets: %v", err)
	}

	var recordToUpdate *googledns.ResourceRecordSet
	for _, record := range recordSets.Rrsets {
		if record.Name == fullRecordName && record.Type == req.GetRecordType() {
			if record.Rrdatas[0] == req.GetIP() {
				fmt.Println("Record already up to date")
				return
			}
			recordToUpdate = record
			break
		}
	}

	if recordToUpdate == nil {
		log.Fatalf("Record %s of type %s not found in zone %s", req.GetRecordName(), req.GetRecordType(), req.GetZone())
	}

	// Prepare the new record
	newRecord := &googledns.ResourceRecordSet{
		Name:    fullRecordName,
		Type:    req.GetRecordType(),
		Ttl:     recordToUpdate.Ttl, // Keep the same TTL
		Rrdatas: []string{req.GetIP()},
	}

	// Create a change request to update the record
	change := &googledns.Change{
		Additions: []*googledns.ResourceRecordSet{newRecord},
		Deletions: []*googledns.ResourceRecordSet{recordToUpdate},
	}

	_, err = s.client.Changes.Create(config.Conf.GCP.ProjectId, req.GetZone(), change).Do()
	if err != nil {
		log.Fatalf("Failed to create change: %v", err)
	}

	fmt.Println("Change applied successfully")
}
