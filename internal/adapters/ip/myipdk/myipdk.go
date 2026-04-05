package myipdk

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/netip"
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

	req.Header.Set("User-Agent", "curl/8.4.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ip.Get returned an error at http.DefaultClient.Do: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ip.Get returned an error at io.ReadAll: %w", err)
	}

	ip, err := netip.ParseAddr(string(body))
	if err != nil {
		return nil, fmt.Errorf("ip.Get returned an error at netip.ParseAddr: %w", err)
	}

	return &ip, nil
}
