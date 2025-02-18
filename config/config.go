package config

import (
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	Provider Provider
	Updates  []Update
	Tracing  Tracing
}

type Provider struct {
	Name   string
	Config map[string]interface{}
}

type Update struct {
	Domain  string
	Zone    string
	Records []string
	Type    string
}

type Tracing struct {
	Enabled       bool
	Host          string
	Port          int
	AllowInsecure bool
}

var Conf Config

func LoadConfig() {
	conf := Config{}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal("Config file not found: %w", err)
		} else {
			log.Fatal("Config file found but error occured: %w", err)
		}
	}

	err := viper.Unmarshal(&conf)
	if err != nil {
		log.Fatal("Can't unmarshal config file:", err)
	}

	Conf = conf

	log.Printf("Config loaded: %+v", conf)
}

func (c Config) GetProviderString(s string) string {
	return viper.GetString(fmt.Sprintf("provider.config.%s", s))
}

func (c Config) GetProviderInt(i int) int {
	return viper.GetInt(fmt.Sprintf("provider.config.%s", strconv.Itoa(i)))
}
