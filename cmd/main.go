package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/orkarstoft/dns-updater/internal/adapters/cache"
	_ "github.com/orkarstoft/dns-updater/internal/adapters/dns/digitalocean"
	_ "github.com/orkarstoft/dns-updater/internal/adapters/dns/gcp"
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
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to initialize file cache")
		}
	} else {
		cacheSvc = cache.NewNoOpCache()
	}

	updaterSvc := service.NewDDNSService(dnsAdapter, ipsvc, cacheSvc, &log.Logger)

	err = updaterSvc.Run(ctx, cfg.Updates)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to run DNS updater service")
	}
}
