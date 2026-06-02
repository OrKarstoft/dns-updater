package cache

import (
	"context"
	"net/netip"
	"os"
	"strings"
)

// FileCache implements the ports.IPCache interface by storing the IP in a local file.
type FileCache struct {
	filepath string
}

// NewFileCache creates a new FileCache instance.
func NewFileCache(filepath string) *FileCache {
	return &FileCache{
		filepath: filepath,
	}
}

// GetLastIP reads the IP address from the file.
func (c *FileCache) GetLastIP(ctx context.Context) (*netip.Addr, error) {
	data, err := os.ReadFile(c.filepath)
	if err != nil {
		// This happens naturally on the first run before the file is created.
		// Treat it as an empty cache.
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	ipStr := strings.TrimSpace(string(data))
	ip, err := netip.ParseAddr(ipStr)
	if err != nil {
		return nil, err
	}

	return &ip, nil
}

// SetLastIP writes the IP address to the file.
func (c *FileCache) SetLastIP(ctx context.Context, ip *netip.Addr) error {
	// 0644 gives read/write permissions to the owner, and read to everyone else.
	return os.WriteFile(c.filepath, []byte(ip.String()), 0o644)
}
