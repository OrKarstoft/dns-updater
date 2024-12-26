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

	var svc dns.DNSImpl
	// TODO: Handle multiple DNS providers better than this
	if config.Conf.DOToken != "" {
		svc = digitalocean.NewService(config.Conf.DOToken)
	} else if config.Conf.GCP != (config.GCP{}) {
		svc = gcp.NewService()
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

			svc.SetRecord(dnsReq)
		}
	}
}
