package simply

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/netip"
	"strings"

	"github.com/orkarstoft/dns-updater/internal/config"
	"github.com/orkarstoft/dns-updater/internal/core/ports"
	"github.com/orkarstoft/dns-updater/internal/registry"
	"github.com/rs/zerolog/log"
)

const (
	providerName = "simply"
	baseURL      = "https://api.simply.com/2"
)

func init() {
	registry.RegisterDNSProvider(providerName, NewFromConfig)
}

type Provider struct {
	accountName string
	apiKey      string
	httpClient  *http.Client
}

type SimplyRecord struct {
	ID   int    `json:"record_id,omitempty"`
	Name string `json:"name"`
	Type string `json:"type"`
	Data string `json:"data"`
	TTL  int    `json:"ttl,omitempty"`
}

type recordsResponse struct {
	Records []SimplyRecord `json:"records"`
}

func NewFromConfig(cfg config.Provider) (ports.DNSProvider, error) {
	accountName := cfg.GetString("accountName")
	apiKey := cfg.GetString("apiKey")

	if accountName == "" || apiKey == "" {
		log.Fatal().Msg("Simply account name and API key are required")
	}

	return &Provider{
		accountName: accountName,
		apiKey:      apiKey,
		httpClient:  &http.Client{},
	}, nil
}

func (p *Provider) doRequest(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	url := fmt.Sprintf("%s/%s", baseURL, strings.TrimPrefix(path, "/"))
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(p.accountName, p.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("simply API error: status=%d body=%s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (p *Provider) UpdateRecord(ctx context.Context, zone, domain, recordName, recordType string, ip *netip.Addr) error {
	records, err := p.getRecords(ctx, domain)
	if err != nil {
		return err
	}
	log.Debug().Msgf("Found %d records in domain %s", len(records), domain)

	record := findMatchingRecord(records, recordType, recordName)
	if record == nil {
		log.Debug().Msgf("Record %s not found in zone %s, creating new record", fmt.Sprintf("%s.%s", recordName, domain), zone)
		if err := p.createDNSRecord(ctx, recordType, recordName, domain, zone, ip); err != nil {
			return err
		}
		return nil
	}

	log.Debug().Msgf("Record %s found in zone %s, updating record", fmt.Sprintf("%s.%s", recordName, domain), zone)
	if record.Data == ip.String() {
		log.Debug().Msg("Record already up to date")
		return nil
	}

	if err := p.updateDNSRecord(ctx, domain, ip, record.ID); err != nil {
		return err
	}

	log.Debug().Msg("Record updated")
	return nil
}

func (p *Provider) getRecords(ctx context.Context, domain string) ([]SimplyRecord, error) {
	path := fmt.Sprintf("/my/products/%s/dns/records", domain)
	respBody, err := p.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list simply domain records: %w", err)
	}

	var parsedResp recordsResponse
	if err := json.Unmarshal(respBody, &parsedResp); err != nil {
		// Fallback for flat array response in case the API spec differs slightly
		var flatRecords []SimplyRecord
		if err2 := json.Unmarshal(respBody, &flatRecords); err2 == nil {
			return flatRecords, nil
		}
		return nil, fmt.Errorf("failed to parse simply API response: %w", err)
	}

	return parsedResp.Records, nil
}

func findMatchingRecord(records []SimplyRecord, recordType, recordName string) *SimplyRecord {
	for i := range records {
		if records[i].Name == recordName && records[i].Type == recordType {
			return &records[i]
		}
	}
	return nil
}

func (p *Provider) updateDNSRecord(ctx context.Context, domain string, ip *netip.Addr, recordID int) error {
	path := fmt.Sprintf("/my/products/%s/dns/records/%d", domain, recordID)
	updateReq := SimplyRecord{
		Data: ip.String(),
	}

	_, err := p.doRequest(ctx, http.MethodPut, path, updateReq)
	if err != nil {
		return fmt.Errorf("failed to update simply DNS record: %w", err)
	}
	return nil
}

func (p *Provider) createDNSRecord(ctx context.Context, recordType, recordName, domain, zone string, ip *netip.Addr) error {
	path := fmt.Sprintf("/my/products/%s/dns/records", domain)
	createReq := SimplyRecord{
		Name: recordName,
		Type: recordType,
		Data: ip.String(),
		TTL:  3600,
	}

	_, err := p.doRequest(ctx, http.MethodPost, path, createReq)
	if err != nil {
		return fmt.Errorf("failed to create simply DNS record: %w", err)
	}

	log.Debug().
		Str("name", recordName).
		Str("type", recordType).
		Str("ip", ip.String()).
		Str("zone", zone).
		Msg("Record created")
	return nil
}
