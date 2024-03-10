package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/digitalocean/godo"
)

func main() {
	doToken := os.Getenv("DO_TOKEN")
	if doToken == "" {
		log.Fatal("DO_TOKEN is not set")
	}

	client := godo.NewFromToken(doToken)

	ctx := context.TODO()

	// Get IP address
	req, err := http.NewRequest("GET", "https://myip.dk", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "curl/8.4.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	ip := string(body)

	records, _, err := client.Domains.Records(ctx, "karstoft.pro", &godo.ListOptions{WithProjects: true})
	if err != nil {
		log.Fatal(err)
	}

	for _, record := range records {
		if record.Name == "ha" {
			if record.Data == ip {
				fmt.Println("Record is up to date")
				break
			}
			_, _, err := client.Domains.EditRecord(ctx, "karstoft.pro", record.ID, &godo.DomainRecordEditRequest{
				Data: ip,
			})
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Record updated")
			break
		}
	}
}
