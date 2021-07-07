package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	client = http.Client{Timeout: 30 * time.Second}
	apiKey = os.Getenv("APIKEY")
)

type releaseServiceCache struct {
	mu    sync.RWMutex
	songs map[string]songsByDateMeta
}

func (c *releaseServiceCache) callRS(date, url string) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	q := req.URL.Query()
	q.Set("released_at", date)
	q.Set("api_key", apiKey)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("%v from %s", err, req.URL)
		return err
	} else if resp.StatusCode != http.StatusOK {
		err := errors.New(fmt.Sprintf("%d from %s", resp.StatusCode, req.URL))
		log.Print(err)
		return err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	if _, err := dec.Token(); err != nil {
		log.Fatal(err)
	}

	exp := time.Now().AddDate(0, 0, 30)

	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.songs[date]; !ok && url == dailyURL {
		c.songs[date] = songsByDateMeta{exp: exp}
	}

Songs:
	for dec.More() {
		var s songMeta
		if err := dec.Decode(&s); err != nil {
			log.Fatal(err)
		}

		if cached := c.songs[s.Date].Songs; cached != nil {
			for _, v := range *cached {
				if v.ID == s.ID {
					continue Songs
				}
			}

			*cached = append(*cached, s)
			continue
		}

		c.songs[s.Date] = songsByDateMeta{
			exp:         exp,
			songsByDate: songsByDate{Date: s.Date, Songs: &[]songMeta{s}},
		}
	}

	return nil
}
