package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/orkarstoft/dns-updater/application"
	"github.com/orkarstoft/dns-updater/config"
	"github.com/orkarstoft/dns-updater/dns"
	"github.com/orkarstoft/dns-updater/dns/providers/digitalocean"
	"github.com/orkarstoft/dns-updater/dns/providers/gcp"
	"github.com/orkarstoft/dns-updater/dns/tracing"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	config.LoadConfig()

	dnsProvider := getDNSProvider()

	options := application.Options{
		Ctx:            ctx,
		ProviderClient: dnsProvider,
	}

	if config.Conf.TracingEnabled {
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

func getDNSProvider() dns.DNSImpl {
	var dnsProvider dns.DNSImpl
	switch config.Conf.Provider.Name {
	case "googlecloudplatform":
		dnsProvider = gcp.NewService()
	case "digitalocean":
		dnsProvider = digitalocean.NewService(config.Conf.GetProviderString("token"))
	default:
		log.Fatal("No vaild DNS provider specified")
	}
	return dnsProvider
}
