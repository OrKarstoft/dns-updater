package service

import (
	"context"

	"github.com/orkarstoft/dns-updater/internal/config"
	"github.com/orkarstoft/dns-updater/internal/core/ports"
	"github.com/rs/zerolog"
)

type DNSService struct {
	dns    ports.DNSProvider
	ip     ports.IPFetcher
	cache  ports.IPCache
	logger *zerolog.Logger
}

func NewDDNSService(dns ports.DNSProvider, ip ports.IPFetcher, cache ports.IPCache, logger *zerolog.Logger) *DNSService {
	return &DNSService{
		dns:    dns,
		ip:     ip,
		cache:  cache,
		logger: logger,
	}
}

func (s *DNSService) Run(ctx context.Context, cfg []config.Update) error {
	s.logger.Info().Msg("starting DNS update check")

	// 1. Fetch current WAN IP (Guaranteed to be a valid netip.Addr by the adapter)
	currentIP, err := s.ip.Get(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to fetch WAN IP")
		return err
	}

	// 2. Check against cache
	// We ignore the error here because if the cache is empty or fails to read
	// (e.g., first run, or file deleted), we want to proceed with the update anyway.
	lastIP, _ := s.cache.GetLastIP(ctx)

	// Because netip.Addr is a comparable value type, we can check equality directly.
	if lastIP != nil && currentIP != nil && *currentIP == *lastIP {
		s.logger.Info().Msg("IP has not changed, skipping update")
		return nil
	}

	for _, cfg := range cfg {
		for _, record := range cfg.Records {
			err := s.dns.UpdateRecord(ctx, cfg.Zone, cfg.Domain, record, cfg.Type, currentIP)
			if err != nil {
				s.logger.Error().Err(err).Str("domain", cfg.Domain).Msg("failed to update domain")
				// We continue to the next domain instead of returning immediately,
				// so that one failure doesn't block updates for other domains.
				continue
			}
			s.logger.Info().Str("new_ip", currentIP.String()).Msg("successfully updated DNS record")
		}
	}

	// 4. Update the cache for the next run
	err = s.cache.SetLastIP(ctx, currentIP)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to update IP cache")
		// We don't return an error here because the main function of this service is to update DNS records,
		// and a cache failure shouldn't prevent that. The cache is just an optimization to avoid unnecessary updates.
	}

	return nil
}
