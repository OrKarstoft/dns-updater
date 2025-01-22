package application

import (
	"context"
	"log"

	domain "github.com/orkarstoft/dns-updater"
	"github.com/orkarstoft/dns-updater/config"
	"github.com/orkarstoft/dns-updater/dns"
	"github.com/orkarstoft/dns-updater/ip"
)

type Service struct {
	ctx            context.Context
	providerClient dns.DNSImpl
}

type Options struct {
	Ctx            context.Context
	ProviderClient dns.DNSImpl
}

func New(opts Options) *Service {
	if opts.Ctx == nil {
		opts.Ctx = context.Background()
	}

	if opts.ProviderClient == nil {
		log.Fatal("No valid DNS provider specified")
	}

	return &Service{
		ctx:            opts.Ctx,
		providerClient: opts.ProviderClient,
		tracer:         opts.Tracer,
	}
}

func (s *Service) Run() {
	actualIP := ip.Get()

	for _, update := range config.Conf.Updates {
		for _, record := range update.Records {
			dnsReq := domain.NewDNSRequest(record, update.Domain, update.Zone, actualIP, update.Type)
			if dnsReq == nil {
				log.Fatalf("Invalid DNS request: %+v", dnsReq)
			}

			s.providerClient.UpdateRecord(dnsReq)
		}
	}
}
