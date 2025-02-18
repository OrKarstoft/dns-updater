package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/orkarstoft/dns-updater/application"
	"github.com/orkarstoft/dns-updater/config"
	"github.com/orkarstoft/dns-updater/dns/factory"
	"github.com/orkarstoft/dns-updater/dns/tracing"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	config.LoadConfig()

	dnsProvider, err := factory.CreateProvider(factory.ProviderConfig{
		Type:   factory.ProviderType(config.Conf.Provider.Name),
		Config: config.Conf.Provider.Config,
	})
	if err != nil {
		log.Fatalf("Failed to create DNS provider: %v", err)
	}

	options := application.Options{
		Ctx:            ctx,
		ProviderClient: dnsProvider,
	}

	if config.Conf.Tracing.Enabled {
		tracingService, shutdownTracer := tracing.NewService(ctx, dnsProvider)
		defer shutdownTracer(ctx)

		traceCtx, span := tracingService.Tracer().Start(ctx, "root")
		defer span.End()

		options.Ctx = traceCtx
		options.Tracer = tracingService.Tracer()
		options.ProviderClient = tracingService
	}

	service := application.New(options)

	service.Run()
}
