package main

import (
	"context"
	"log"

	"github.com/orkarstoft/dns-updater/application"
	"github.com/orkarstoft/dns-updater/config"
	"github.com/orkarstoft/dns-updater/dns"
	"github.com/orkarstoft/dns-updater/dns/providers/digitalocean"
	"github.com/orkarstoft/dns-updater/dns/providers/gcp"
)

func main() {
	config.LoadConfig()

	dnsProvider := getDNSProvider()

	options := application.Options{
		Ctx:            context.Background(),
		ProviderClient: dnsProvider,
	}

	if config.Conf.TracingEnabled {
		options.Tracer = nil
	}

	service := application.New(options)

	service.Run()
}

func getDNSProvider() dns.DNSImpl {
	var dnsProvider dns.DNSImpl
	switch config.Conf.Provider.Name {
	case "googlecloudplatform":
		dnsProvider = gcp.NewService()
	case "digitalocean":
		dnsProvider = digitalocean.NewService(config.Conf.GetProviderString("token"))
	default:
		log.Fatal("No vaild DNS provider specified")
	}
	return dnsProvider
}
