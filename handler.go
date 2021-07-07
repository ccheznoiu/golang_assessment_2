package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"sync"
	"time"
)

func releasesAPI(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	} else if req.URL.Path != releasesURL {
		http.NotFound(w, req)
		return
	}

	q := req.URL.Query()
	artist := q.Get("artist")

	var (
		from, until time.Time
		err         error
	)
	if from, err = time.Parse(iso8601, q.Get("from")); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if until, err = time.Parse(iso8601, q.Get("until")); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if from.After(until) {
		http.Error(w, "\"from\" after \"until\"", http.StatusBadRequest)
		return
	}

	var resp []*songsByDate
	var wg sync.WaitGroup

	songC := make(chan *songsByDate)

	go func(wg *sync.WaitGroup, resp *[]*songsByDate) {
		defer wg.Done()

		for s := range songC {
			*resp = append(*resp, s)
		}
	}(&wg, &resp)

	wg.Add(1)
	getSongs(from, until, artist, false, false, songC, &wg)

	wg.Wait()
	wg.Add(1)
	close(songC)
	wg.Wait()

	sort.Slice(resp, func(i, j int) bool { return resp[i].Date < resp[j].Date })

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	enc.Encode(resp)
}
