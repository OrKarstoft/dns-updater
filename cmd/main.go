package main

import (
	"log"

	domain "github.com/orkarstoft/dns-updater"
	"github.com/orkarstoft/dns-updater/config"
	"github.com/orkarstoft/dns-updater/dns"
	"github.com/orkarstoft/dns-updater/dns/providers/digitalocean"
	"github.com/orkarstoft/dns-updater/dns/providers/gcp"
	"github.com/orkarstoft/dns-updater/ip"
)

func main() {
	config.LoadConfig()

	var dnsService dns.DNSImpl
	// TODO: Handle multiple DNS providers better than this
	if config.Conf.DOToken != "" {
		dnsService = digitalocean.NewService(config.Conf.DOToken)
	} else if config.Conf.GCP != (config.GCP{}) {
		dnsService = gcp.NewService()
	} else {
		log.Fatal("No valid DNS provider found")
	}
	ip := ip.GetIP()

	for _, update := range config.Conf.Updates {
		for _, record := range update.Records {
			dnsReq := domain.NewDNSRequest(record, update.Domain, update.Zone, ip, update.Type)
			if dnsReq == nil {
				log.Fatalf("Invalid DNS request: %+v", dnsReq)
			}

			dnsService.UpdateRecord(dnsReq)
		}
	}
}
