package config

import (
	"fmt"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

type Config struct {
	Provider Provider
	Updates  []Update
	Log      Log
	Cache    Cache
	Schedule string `mapstructure:"schedule"`
}

type Provider struct {
	Name     string `mapstructure:"name"`
	SafeMode SafeMode
	Config   map[string]any `mapstructure:"config"`
}

type SafeMode struct {
	// If enabled, dns-updater will only update DNS records where a matching TXT record is found with the correct ownership identifier. This provides an extra layer of safety to prevent accidental updates.
	// Default: true
	Enabled bool `mapstructure:"enabled"`
	// Identifier for the TXT record used to verify ownership of the domain before making changes.
	// Useful if you have multiple instances of dns-updater running for the same domain.
	// Default: "dns-updater"
	TxtOwnerId string `mapstructure:"txtOwnerId"`
	// Prefix for the TXT record used to verify ownership. The full TXT record name will be constructed as "<txt_prefix>.<record_name>". This allows you to use a consistent naming convention for your ownership verification records.
	// Default: "dns-updater-safemode"
	TxtPrefix string `mapstructure:"txtPrefix"`
}

type Update struct {
	Domain  string
	Zone    string
	Records []string
	Type    string
}

type LogType string

const (
	LOGTYPE_JSON   LogType = "json"
	LOGTYPE_PRETTY LogType = "pretty"
	LOGTYPE_FILE   LogType = "file"
)

type Log struct {
	Level string
	Type  LogType
}

type Cache struct {
	Enabled  bool
	FilePath string
}

func LoadConfig() (*Config, error) {
	var conf Config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Default config
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.type", "pretty")
	viper.SetDefault("provider.safemode.enabled", true)
	viper.SetDefault("provider.safemode.txt_owner_id", "dns-updater")
	viper.SetDefault("provider.safemode.txt_prefix", "dns-updater-safemode")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file not found: %w", err)
		}
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Note: Assuming DecodeLogLevelHookFunc and DecodeLogTypeHookFunc are defined in this package
	err := viper.Unmarshal(&conf, viper.DecodeHook(
		mapstructure.ComposeDecodeHookFunc(
		// DecodeLogLevelHookFunc(),
		// DecodeLogTypeHookFunc(),
		),
	))
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal config file: %w", err)
	}

	return &conf, nil
}

// Retaining your getter methods for Viper fallback/utility
func (p Provider) GetString(s string) string {
	if s == "name" {
		return viper.GetString("provider.name")
	}
	return viper.GetString(fmt.Sprintf("provider.config.%s", s))
}

func (p Provider) GetInt(s string) int {
	return viper.GetInt(fmt.Sprintf("provider.config.%s", s))
}

func (p Provider) GetBool(s string) bool {
	return viper.GetBool(fmt.Sprintf("provider.config.%s", s))
}
