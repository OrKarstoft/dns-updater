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

	if config.Conf.Tracing.Enabled && opts.Tracer == nil {
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
			if err := s.processRecord(record, update, actualIP); err != nil {
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

// processRecord handles creating the DNS request and updating the record.
func (s *Service) processRecord(record string, update config.Update, actualIP string) error {
	_, span := s.startSpan("Creating DNS Request")

	dnsReq := domain.NewDNSRequest(record, update.Domain, update.Zone, actualIP, update.Type)
	if dnsReq == nil {
		return fmt.Errorf("invalid DNS request: %+v", dnsReq)
	}

	s.endSpan(span)

	return s.providerClient.UpdateRecord(s.ctx, dnsReq)
}

// startSpan safely starts a tracing span if the tracer is available.
func (s *Service) startSpan(name string) (context.Context, trace.Span) {
	if s.tracer != nil {
		return s.tracer.Start(s.ctx, name)
	}
	return s.ctx, nil
}

// endSpan safely ends a span if it exists.
func (s *Service) endSpan(span trace.Span) {
	if span != nil {
		span.End()
	}
}
