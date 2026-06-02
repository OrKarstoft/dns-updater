package myipdk

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/netip"
	"strings"
	"time"
)

type MyIPDK struct{}

func New() *MyIPDK {
	return &MyIPDK{}
}

func (s *MyIPDK) Get(ctx context.Context) (*netip.Addr, error) {
	ctx, cancelCtx := context.WithTimeout(ctx, 5*time.Second)
	defer cancelCtx()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://myip.dk", nil)
	if err != nil {
		return nil, fmt.Errorf("ip.Get returned an error at http.NewRequest: %w", err)
	}

	// Fake curl, so we don't get an HTML page back.
	req.Header.Set("User-Agent", "curl/8.20.0")
	req.Header.Set("Content-Type", "text/plain")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ip.Get returned an error at http.DefaultClient.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ip.Get unexpected status %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ip.Get returned an error at io.ReadAll: %w", err)
	}

	// The response from myip.dk may contain trailing whitespace, so we trim it before parsing.
	cleanBody := strings.TrimSpace(string(body))

	ip, err := netip.ParseAddr(string(cleanBody))
	if err != nil {
		return nil, fmt.Errorf("ip.Get returned an error at netip.ParseAddr: %w", err)
	}

	return &ip, nil
}
