package main

import (
	"encoding/json"
	"sync"
	"time"
)

type songsByDateMeta struct {
	exp time.Time
	songsByDate
}

type songsByDate struct {
	Date  string      `json:"released_at"`
	Songs *[]songMeta `json:"songs"`
}

func (r songsByDateMeta) isExpired() bool {
	return r.exp.Before(time.Now())
}

func (r songsByDateMeta) hasReleases() bool {
	return r.Songs != nil
}

func (r *songsByDate) filterByArtist(artist string) bool {
	var byArtist []songMeta
	for _, v := range *r.Songs {
		if v.Artist == artist {
			byArtist = append(byArtist, v)
		}
	}

	if len(byArtist) == 0 {
		return false
	}

	r.Songs = &byArtist
	return true
}

type songMeta struct {
	Date string `json:"released_at"`
	ID   string `json:"song_id"`
	song
}

type song struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
}

func (s *songMeta) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.song)
}

func getSongs(from, until time.Time, artist string, callDaily, useExpired bool, songC chan *songsByDate, wg *sync.WaitGroup) {
	defer wg.Done()

	exp := make([]time.Time, 0, 31)

	for ; from.Before(until) || from.Equal(until); from = from.AddDate(0, 0, 1) {
		if callDaily && cache.callRS(from.Format(iso8601), dailyURL) != nil {
			useExpired = true
		}

		date := from.Format(iso8601)
		cache.mu.RLock()
		cached, ok := cache.songs[date]
		cache.mu.RUnlock()

		switch {
		case !ok && useExpired: // new query, downstream error on recursive call: give up
			continue
		case !ok || (cached.isExpired() && !useExpired): // new or expired query, first pass: to call downstream
			exp = append(exp, from)
		case cached.hasReleases() && (artist == "" || cached.filterByArtist(artist)): // is usable
			songC <- &songsByDate{Date: date, Songs: cached.Songs}
		}

		if from.Equal(until) || from.Equal(from.AddDate(0, 1, -from.Day())) {
			if l := len(exp); l > 0 {
				from, until := exp[0], exp[0]
				callDaily := true
				useExpired := false
				if l > 24 {
					callDaily = false
					useExpired = cache.callRS(from.Format(yyyymm), monthlyURL) != nil
				}

				if l > 1 {
					for _, v := range exp[1:] {
						if v.Day()-until.Day() == 1 {
							until = v
							continue
						}
						wg.Add(1)
						go getSongs(from, until, artist, callDaily, useExpired, songC, wg)
						from, until = v, v
					}
				}

				wg.Add(1)
				go getSongs(from, until, artist, callDaily, useExpired, songC, wg)

			}
		}
	}
}
