package application

import (
	"context"
	"fmt"
	"log"

	domain "github.com/orkarstoft/dns-updater"
	"github.com/orkarstoft/dns-updater/config"
	"github.com/orkarstoft/dns-updater/dns"
	"github.com/orkarstoft/dns-updater/ip"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	ctx            context.Context
	providerClient dns.DNSImpl
	tracer         trace.Tracer
}

type Options struct {
	Ctx            context.Context
	ProviderClient dns.DNSImpl
	Tracer         trace.Tracer
}

func New(opts Options) *Service {
	if opts.Ctx == nil {
		fmt.Println("No context provided, creating a blank")
		opts.Ctx = context.Background()
	}

	if opts.ProviderClient == nil {
		log.Fatal("No valid DNS provider specified")
	}

	if config.Conf.TracingEnabled && opts.Tracer == nil {
		log.Fatal("No valid tracer specified")
	}

	return &Service{
		ctx:            opts.Ctx,
		providerClient: opts.ProviderClient,
		tracer:         opts.Tracer,
	}
}

func (s *Service) Run() {
	actualIP := ip.Get()

	var errs []error
	for _, update := range config.Conf.Updates {
		for _, record := range update.Records {
			_, span := s.tracer.Start(s.ctx, "Creating DNS Request")
			dnsReq := domain.NewDNSRequest(record, update.Domain, update.Zone, actualIP, update.Type)
			if dnsReq == nil {
				log.Fatalf("Invalid DNS request: %+v", dnsReq)
			}
			span.End()
			err := s.providerClient.UpdateRecord(s.ctx, dnsReq)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		for _, err := range errs {
			log.Println(err)
		}
	} else {
		log.Println("All records updated successfully")
	}
}
