package application

import (
	"context"
	"errors"
	"fmt"
	"log"

	domain "github.com/orkarstoft/dns-updater"
	"github.com/orkarstoft/dns-updater/config"
	"github.com/orkarstoft/dns-updater/dns"
	"github.com/orkarstoft/dns-updater/ip"
	"go.opentelemetry.io/otel/trace"
)

type IPResolver interface {
	Get() (string, error)
}

type Service struct {
	ctx         context.Context
	dnsProvider dns.DNSImpl
	ipResolver  IPResolver
	tracer      trace.Tracer
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
		ctx:         opts.Ctx,
		dnsProvider: opts.ProviderClient,
		tracer:      opts.Tracer,
	}
}

func (s *Service) Run() {
	actualIP, err := ip.Get()
	if err != nil {
		log.Println(fmt.Errorf("application exited because of error from ip.Get(): %w", err))
	}

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
			var dnsErr *dns.DNSProviderError
			if errors.As(err, &dnsErr) {
				log.Printf("DNS provider error: %v", err)
			} else {
				log.Printf("General error: %v", err)
			}
		}
	} else {
		log.Println("All records updated successfully")
	}
}

// processRecord handles creating the DNS request and updating the record.
func (s *Service) processRecord(record string, update config.Update, actualIP string) error {
	_, span := s.startSpan("Creating DNS Request")

	dnsReq, err := domain.NewDNSRequest(record, update.Domain, update.Zone, actualIP, update.Type)
	if dnsReq == nil {
		return fmt.Errorf("failed to create DNS request for record %s: %w", record, err)
	}

	s.endSpan(span)

	return s.dnsProvider.UpdateRecord(s.ctx, dnsReq)
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
