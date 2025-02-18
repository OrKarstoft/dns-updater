package tracing

import (
	"context"
	"log"

	domain "github.com/orkarstoft/dns-updater"
	"github.com/orkarstoft/dns-updater/dns"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	tracer         trace.Tracer
	providerClient dns.DNSImpl
}

func NewService(ctx context.Context, providerClient dns.DNSImpl) (*Service, func(context.Context) error) {
	conn, err := initConn()
	if err != nil {
		log.Fatal(err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// The service name used to display traces in backends
			serviceName,
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	shutdownTracerProvider, err := initTracerProvider(ctx, res, conn)
	if err != nil {
		log.Fatal(err)
	}

	tracer := otel.Tracer("dns-updater")

	return &Service{
		tracer:         tracer,
		providerClient: providerClient,
	}, shutdownTracerProvider
}

func (s *Service) Tracer() trace.Tracer {
	return s.tracer
}

func (s *Service) UpdateRecord(ctx context.Context, req *domain.DNSRequest) error {
	newCtx, span := s.tracer.Start(ctx, "Updating Record")
	defer span.End()
	err := s.providerClient.UpdateRecord(newCtx, req)
	return err
}
