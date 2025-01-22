package tracing

import (
	"context"
	"log"
	"os"
	"os/signal"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	ctx            context.Context
	tracer         trace.Tracer
	shutdownTracer func(context.Context) error
}

func NewService(ctx context.Context) (*Service, func(context.Context) error) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

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

	name := "dns-updater"
	tracer := otel.Tracer(name)

	// Attributes represent additional key-value descriptors that can be bound
	// to a metric observer or recorder.
	// commonAttrs := []attribute.KeyValue{
	// 	attribute.String("attrA", "chocolate"),
	// 	attribute.String("attrB", "raspberry"),
	// 	attribute.String("attrC", "vanilla"),
	// }

	return &Service{
		ctx:    ctx,
		tracer: tracer,
	}, shutdownTracerProvider
}

// func (s *ServiceWithTracing) UpdateRecord(req *domain.DNSRequest) {
// 	_, span := s.tracer.Start(s.ctx, "UpdateRecord")
// 	defer span.End()
//
// 	s.service.UpdateRecord(req)
// }
