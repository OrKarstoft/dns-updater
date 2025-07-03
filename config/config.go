package config

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-viper/mapstructure/v2"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type Config struct {
	Provider Provider
	Updates  []Update
	Tracing  Tracing
	Log      Log
}

type Provider struct {
	name   string
	config map[string]interface{}
}

type Update struct {
	Domain  string
	Zone    string
	Records []string
	Type    string
}

type Tracing struct {
	enabled       bool
	host          string
	port          int
	allowInsecure bool
}

type LogType string

const (
	LOGTYPE_JSON   LogType = "json"
	LOGTYPE_PRETTY LogType = "pretty"
	LOGTYPE_FILE   LogType = "file"
)

type Log struct {
	Level zerolog.Level
	Type  LogType
}

var Conf Config

func LoadConfig() {
	conf := Config{}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Default config
	//// Tracing
	viper.SetDefault("tracing.enabled", false)
	//// Log
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.type", "pretty")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal("Config file not found: %w", err)
		} else {
			log.Fatal("Config file found but error occured: %w", err)
		}
	}

	err := viper.Unmarshal(&conf, viper.DecodeHook(
		mapstructure.ComposeDecodeHookFunc(
			DecodeLogLevelHookFunc(),
			DecodeLogTypeHookFunc(),
		),
	))
	if err != nil {
		log.Fatal("Can't unmarshal config file:", err)
	}

	Conf = conf
}

func (p Provider) GetString(s string) string {
	if s == "name" {
		return viper.GetString("provider.name")
	}
	return viper.GetString(fmt.Sprintf("provider.config.%s", s))
}

func (p Provider) GetInt(i int) int {
	return viper.GetInt(fmt.Sprintf("provider.config.%s", strconv.Itoa(i)))
}

func (p Provider) GetBool(s string) bool {
	return viper.GetBool(fmt.Sprintf("provider.config.%s", s))
}

func (t Tracing) GetString(s string) string {
	return viper.GetString(fmt.Sprintf("tracing.%s", s))
}

func (t Tracing) GetInt(s string) int {
	return viper.GetInt(fmt.Sprintf("tracing.%s", s))
}

func (t Tracing) GetBool(s string) bool {
	return viper.GetBool(fmt.Sprintf("tracing.%s", s))
}
