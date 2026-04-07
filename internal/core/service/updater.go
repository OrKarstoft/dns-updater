package service

import (
	"context"
	"fmt"
	"net/netip"

	"github.com/orkarstoft/dns-updater/internal/config"
	"github.com/orkarstoft/dns-updater/internal/core/ports"
	"github.com/rs/zerolog"
)

type DNSService struct {
	dns      ports.DNSProvider
	ip       ports.IPFetcher
	cache    ports.IPCache
	logger   *zerolog.Logger
	safeMode bool
}

func NewDDNSService(dns ports.DNSProvider, ip ports.IPFetcher, cache ports.IPCache, logger *zerolog.Logger, safeMode bool) *DNSService {
	return &DNSService{
		dns:      dns,
		ip:       ip,
		cache:    cache,
		logger:   logger,
		safeMode: safeMode,
	}
}

func (s *DNSService) Run(ctx context.Context, cfg []config.Update) error {
	s.logger.Info().Msg("starting DNS update check")

	// 1. Fetch current WAN IP
	currentIP, err := s.ip.Get(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to fetch WAN IP")
		return err
	}

	// 2. Check against cache
	lastIP, _ := s.cache.GetLastIP(ctx)
	if lastIP != nil && currentIP != nil && *currentIP == *lastIP {
		s.logger.Info().Msg("IP has not changed, skipping update")
		return nil
	}

	for _, updateCfg := range cfg {
		for _, recordName := range updateCfg.Records {
			if s.safeMode {
				err := s.updateRecordSafe(ctx, updateCfg.Zone, updateCfg.Domain, recordName, updateCfg.Type, currentIP)
				if err != nil {
					s.logger.Error().Err(err).Str("domain", updateCfg.Domain).Msg("failed to update domain")
					continue
				}
			} else {
				err := s.updateRecord(ctx, updateCfg.Zone, updateCfg.Domain, recordName, updateCfg.Type, currentIP)
				if err != nil {
					s.logger.Error().Err(err).Str("domain", updateCfg.Domain).Msg("failed to update domain")
					continue
				}
			}
			s.logger.Info().Str("new_ip", currentIP.String()).Msg("successfully updated DNS record")
		}
	}

	// 4. Update the cache for the next run
	if err := s.cache.SetLastIP(ctx, currentIP); err != nil {
		s.logger.Error().Err(err).Msg("failed to update IP cache")
	}

	return nil
}

func (s *DNSService) updateRecord(ctx context.Context, zone, domain, recordName, recordType string, ip *netip.Addr) error {
	records, err := s.dns.GetRecords(ctx, zone, domain)
	if err != nil {
		return err
	}

	record := findMatchingRecord(records, recordType, recordName)
	if record == nil {
		s.logger.Debug().Msgf("Record %s not found in zone %s, creating new record", fmt.Sprintf("%s.%s", recordName, domain), zone)
		_, err := s.dns.CreateRecord(ctx, zone, domain, ports.DNSRecord{
			Name: recordName,
			Type: recordType,
			Data: ip.String(),
			TTL:  3600,
		})
		return err
	}

	s.logger.Debug().Msgf("Record %s found in zone %s, updating record", fmt.Sprintf("%s.%s", recordName, domain), zone)
	if record.Data == ip.String() {
		s.logger.Debug().Msg("Record already up to date")
		return nil
	}

	return s.dns.UpdateRecord(ctx, zone, domain, record.ID, ports.DNSRecord{Data: ip.String(), TTL: record.TTL})
}

func (s *DNSService) updateRecordSafe(ctx context.Context, zone, domain, recordName, recordType string, ip *netip.Addr) error {
	records, err := s.dns.GetRecords(ctx, zone, domain)
	if err != nil {
		return err
	}

	safemodeRecordName := "dns-updater-safemode." + recordName
	safemodeRecord := findMatchingRecord(records, "TXT", safemodeRecordName)
	record := findMatchingRecord(records, recordType, recordName)

	fullRecordName := fmt.Sprintf("%s.%s", recordName, domain)

	if record != nil && safemodeRecord == nil {
		return fmt.Errorf("record %s exists, but no safemode TXT record found. refusing to touch this record", fullRecordName)
	}

	if record == nil {
		s.logger.Debug().Msgf("Record %s not found in zone %s, creating new record", fullRecordName, zone)
		if _, err := s.dns.CreateRecord(ctx, zone, domain, ports.DNSRecord{Name: recordName, Type: recordType, Data: ip.String(), TTL: 3600}); err != nil {
			return err
		}
		s.logger.Debug().Msgf("Creating safemode TXT record for %s", fullRecordName)
		if _, err := s.dns.CreateRecord(ctx, zone, domain, ports.DNSRecord{Name: safemodeRecordName, Type: "TXT", Data: "managed-by-dns-updater", TTL: 3600}); err != nil {
			return err
		}
		return nil
	}

	s.logger.Debug().Msgf("Record %s found in zone %s, updating record", fullRecordName, zone)
	if record.Data == ip.String() {
		s.logger.Debug().Msg("Record already up to date")
		return nil
	}

	return s.dns.UpdateRecord(ctx, zone, domain, record.ID, ports.DNSRecord{Data: ip.String(), TTL: record.TTL})
}

func findMatchingRecord(records []ports.DNSRecord, recordType, recordName string) *ports.DNSRecord {
	for i := range records {
		if records[i].Name == recordName && records[i].Type == recordType {
			return &records[i]
		}
	}
	return nil
}
