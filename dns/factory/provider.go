package factory

// ProviderType represents the type of DNS provider as a value object
type ProviderType string

const (
	ProviderGCP          ProviderType = "googlecloudplatform"
	ProviderDigitalOcean ProviderType = "digitalocean"
)

// ProviderConfig represents the configuration for a DNS provider
type ProviderConfig struct {
	Type   ProviderType
	Config map[string]interface{}
}
