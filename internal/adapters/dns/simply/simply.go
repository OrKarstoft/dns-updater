package simply

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/netip"
	"strconv"
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

func (p *Provider) doRequest(ctx context.Context, method, path string, body any) ([]byte, error) {
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

func (p *Provider) GetRecords(ctx context.Context, zone, domain string) ([]ports.DNSRecord, error) {
	path := fmt.Sprintf("/my/products/%s/dns/records", domain)
	respBody, err := p.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list simply domain records: %w", err)
	}

	var parsedResp recordsResponse
	if err := json.Unmarshal(respBody, &parsedResp); err != nil {
		var flatRecords []SimplyRecord
		if err2 := json.Unmarshal(respBody, &flatRecords); err2 == nil {
			return toDNSRecords(flatRecords), nil
		}
		return nil, fmt.Errorf("failed to parse simply API response: %w", err)
	}

	return toDNSRecords(parsedResp.Records), nil
}

func (p *Provider) CreateRecord(ctx context.Context, zone, domain string, record ports.DNSRecord) (ports.DNSRecord, error) {
	path := fmt.Sprintf("/my/products/%s/dns/records", domain)
	createReq := SimplyRecord{
		Name: record.Name,
		Type: record.Type,
		Data: record.Data,
		TTL:  record.TTL,
	}

	respBody, err := p.doRequest(ctx, http.MethodPost, path, createReq)
	if err != nil {
		return ports.DNSRecord{}, fmt.Errorf("failed to create simply DNS record: %w", err)
	}

	var createdRecord SimplyRecord
	if err := json.Unmarshal(respBody, &createdRecord); err != nil {
		return ports.DNSRecord{}, fmt.Errorf("failed to parse simply API response: %w", err)
	}

	log.Debug().
		Str("name", record.Name).
		Str("type", record.Type).
		Str("data", record.Data).
		Str("domain", domain).
		Msg("Record created")
	return toDNSRecord(createdRecord), nil
}

func (p *Provider) UpdateRecord(ctx context.Context, zone, domain, recordID string, record ports.DNSRecord) error {
	path := fmt.Sprintf("/my/products/%s/dns/records/%s", domain, recordID)
	updateReq := SimplyRecord{
		Data: record.Data,
	}

	_, err := p.doRequest(ctx, http.MethodPut, path, updateReq)
	if err != nil {
		return fmt.Errorf("failed to update simply DNS record: %w", err)
	}
	return nil
}

func (p *Provider) DeleteRecord(ctx context.Context, zone, domain, recordID string) error {
	path := fmt.Sprintf("/my/products/%s/dns/records/%s", domain, recordID)
	_, err := p.doRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("failed to delete simply DNS record: %w", err)
	}
	return nil
}



func toDNSRecord(r SimplyRecord) ports.DNSRecord {
	return ports.DNSRecord{
		ID:   strconv.Itoa(r.ID),
		Name: r.Name,
		Type: r.Type,
		Data: r.Data,
		TTL:  r.TTL,
	}
}

func toDNSRecords(rs []SimplyRecord) []ports.DNSRecord {
	records := make([]ports.DNSRecord, len(rs))
	for i, r := range rs {
		records[i] = toDNSRecord(r)
	}
	return records
}


