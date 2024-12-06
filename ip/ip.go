package ip

import (
	"io"
	"log"
	"net/http"
)

func GetIP() string {
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

	return string(body)
}
