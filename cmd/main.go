package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"

	"github.com/orkarstoft/dns-updater/application"
	"github.com/orkarstoft/dns-updater/config"
	"github.com/orkarstoft/dns-updater/dns"
	"github.com/orkarstoft/dns-updater/dns/providers/digitalocean"
	"github.com/orkarstoft/dns-updater/dns/providers/gcp"
	"github.com/orkarstoft/dns-updater/dns/tracing"
	"github.com/orkarstoft/dns-updater/logger"
	"github.com/rs/zerolog"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	config.LoadConfig()

	loggerSvc, err := logger.New(config.Conf.Log.Type, config.Conf.Log.Level)
	if err != nil {
		// You can unwrap the error if needed
		var logErr *logger.LoggerError
		if errors.As(err, &logErr) {
			log.Fatal("Logger operation '%s' failed: %v\n", logErr.Operation, logErr.Err)
		}
		// Handle error appropriately
		log.Fatal(err)
	}

	dnsProvider := getDNSProvider(loggerSvc)

	options := application.Options{
		Ctx:            ctx,
		ProviderClient: dnsProvider,
		Logger:         loggerSvc,
	}

	if config.Conf.Tracing.GetBool("enabled") {
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

func getDNSProvider(logger *zerolog.Logger) dns.DNSImpl {
	var dnsProvider dns.DNSImpl
	switch config.Conf.Provider.GetString("name") {
	case "googlecloudplatform":
		dnsProvider = gcp.NewService(logger, config.Conf.Provider.GetString("projectId"), config.Conf.Provider.GetString("credentialsFile"))
	case "digitalocean":
		dnsProvider = digitalocean.NewService(logger, config.Conf.Provider.GetString("token"))
	default:
		logger.Fatal().Msg("No valid DNS provider specified in config file")
	}
	return dnsProvider
}
