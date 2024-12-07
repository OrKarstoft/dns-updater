package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DOToken string
	GCP     GCP
	Updates []Update `mapstructure:"updates"`
}
type GCP struct {
	CredentialsFilePath string `mapstructure:"credentialsFile"`
	ProjectId           string `mapstructure:"projectId"`
}
type Update struct {
	Domain  string
	Zone    string
	Records []string
	Type    string
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
}
