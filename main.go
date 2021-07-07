package main

import (
	"log"
	"net/http"
)

const (
	releasesURL string = "/releases"
	dailyURL    string = "/daily"
	monthlyURL  string = "/monthly"
	iso8601     string = "2006-01-02"
	yyyymm      string = "2006-01"
)

var cache *releaseServiceCache

func main() {
	if apiKey == "" {
		log.Fatal("Set APIKEY variable?")
	}
	cache = &releaseServiceCache{songs: make(map[string]songsByDateMeta)}

	http.HandleFunc(releasesURL, releasesAPI)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
