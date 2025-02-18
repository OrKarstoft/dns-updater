package ip

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

var cachedIp string

func Get() (string, error) {
	if cachedIp != "" {
		return cachedIp, nil
	}

	ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelCtx()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://myip.dk", nil)
	if err != nil {
		return "", fmt.Errorf("ip.Get returned an error at http.NewRequest: %w", err)
	}

	req.Header.Set("User-Agent", "curl/8.4.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("ip.Get returned an error at http.DefaultClient.Do: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ip.Get returned an error at io.ReadAll: %w", err)
	}

	return string(body), nil
}
