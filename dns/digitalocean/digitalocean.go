package digitalocean

import (
	"context"
	"fmt"
	"log"

	"github.com/digitalocean/godo"
	domain "github.com/orkarstoft/dns-updater"
	"github.com/orkarstoft/dns-updater/dns"
)

type Service struct {
	ctx    context.Context
	client *godo.Client
}

func NewService(apiToken string) dns.DNSImpl {
	if apiToken == "" {
		log.Fatal("API token is required")
	}

	client := godo.NewFromToken(apiToken)

	ctx := context.TODO()
	return &Service{ctx: ctx, client: client}
}

func (s *Service) SetRecord(req *domain.DNSRequest) {
	records, _, err := s.client.Domains.Records(s.ctx, req.GetDomain(), &godo.ListOptions{WithProjects: true})
	if err != nil {
		log.Fatal(err)
	}

	for _, record := range records {
		if record.Name == req.GetRecordName() {
			if record.Data == req.GetIP() {
				fmt.Println("Record is up to date")
				break
			}
			_, _, err := s.client.Domains.EditRecord(s.ctx, req.GetDomain(), record.ID, &godo.DomainRecordEditRequest{
				Data: req.GetIP(),
			})
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Record updated")
			break
		}
	}
}
