package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/go-co-op/gocron/v2"
	"github.com/orkarstoft/dns-updater/internal/adapters/cache"
	_ "github.com/orkarstoft/dns-updater/internal/adapters/dns/digitalocean"
	_ "github.com/orkarstoft/dns-updater/internal/adapters/dns/gcp"
	_ "github.com/orkarstoft/dns-updater/internal/adapters/dns/simply"
	"github.com/orkarstoft/dns-updater/internal/adapters/ip/myipdk"
	"github.com/orkarstoft/dns-updater/internal/config"
	"github.com/orkarstoft/dns-updater/internal/core/ports"
	"github.com/orkarstoft/dns-updater/internal/core/service"
	"github.com/orkarstoft/dns-updater/internal/logger"
	"github.com/orkarstoft/dns-updater/internal/registry"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	logger.New(cfg.Log)

	dnsAdapter, err := registry.GetDNSProvider(cfg.Provider)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize DNS provider")
	}

	ipsvc := myipdk.New()

	// Initialize the cache.
	var cacheSvc ports.IPCache
	if cfg.Cache.Enabled {
		cacheSvc = cache.NewFileCache(cfg.Cache.FilePath)
	} else {
		cacheSvc = cache.NewNoOpCache()
	}

	updaterSvc := service.NewDDNSService(dnsAdapter, ipsvc, cacheSvc, &log.Logger, cfg.Provider.SafeMode)

	if cfg.Schedule != "" {
		runScheduled(ctx, updaterSvc, cfg)
	} else {
		runOnce(ctx, updaterSvc, cfg)
	}
}

func runOnce(ctx context.Context, updaterSvc *service.DNSService, cfg *config.Config) {
	if err := updaterSvc.Run(ctx, cfg.Updates); err != nil {
		log.Fatal().Err(err).Msg("Failed to run DNS updater service")
	}
	log.Info().Msg("DNS update complete")
}

func runScheduled(ctx context.Context, updaterSvc *service.DNSService, cfg *config.Config) {
	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create scheduler")
	}
	_, err = s.NewJob(
		gocron.CronJob(cfg.Schedule, false),
		gocron.NewTask(func() {
			if err := updaterSvc.Run(ctx, cfg.Updates); err != nil {
				log.Error().Err(err).Msg("Failed to run DNS updater service")
			}
		}),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create job")
	}

	s.Start()
	log.Info().Msgf("Scheduler started with schedule: %s", cfg.Schedule)

	// Block until the context is cancelled.
	<-ctx.Done()

	log.Info().Msg("Shutting down scheduler")
	if err := s.Shutdown(); err != nil {
		log.Error().Err(err).Msg("Failed to shutdown scheduler")
	}
}
