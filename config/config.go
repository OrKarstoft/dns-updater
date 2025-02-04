package config

import (
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	Provider Provider `mapstructure:"provider"`
	Updates  []Update `mapstructure:"updates"`
	Tracing  Tracing  `mapstructure:"tracing"`
}

type Provider struct {
	Name   string                 `mapstructure:"name"`
	Config map[string]interface{} `mapstructure:"config"`
}

type Update struct {
	Domain  string
	Zone    string
	Records []string
	Type    string
}

type Tracing struct {
	Enabled       bool   `mapstructure:"enabled"`
	Host          string `mapstructure:"host"`
	Port          int    `mapstructure:"port"`
	AllowInsecure bool   `mapstructure:"AllowInsecure"`
}

var Conf Config

func LoadConfig() {
	conf := Config{}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't read config file:", err)
	}

	err = viper.Unmarshal(&conf)
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
