package main

import (
	"log"
	"net/http"
	"time"
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

	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	cache = newCache(ticker.C)

	http.HandleFunc(releasesURL, releasesAPI)

	server := &http.Server{
		ReadTimeout:  time.Second,
		WriteTimeout: time.Minute,
		Addr:         "8000",
	}

	log.Fatal(server.ListenAndServe())
}
