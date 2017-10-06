package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/laurentlbm/feeds/feed"
)

var timeoutDuration = time.Hour

func init() {
	for comic, transform := range feed.Feeds {
		transformFeed(comic, transform)
	}
}

func transformFeed(comicURL string, transform feed.Transform) {
	var lastUpdated time.Time
	http.HandleFunc(fmt.Sprintf("/%s/", comicURL), func(w http.ResponseWriter, r *http.Request) {
		// This is a hack because it seems like I can't do polling in the background without a user generated request
		if time.Now().Sub(lastUpdated) > timeoutDuration {
			transform.Do(r)
			lastUpdated = time.Now()
		}
		w.Header().Add("content-type", "text/xml")
		w.Write([]byte(transform.Generate()))
	})
	http.HandleFunc(fmt.Sprintf("/%s/favicon.ico", comicURL), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "image/png")
		w.Write(transform.GetFavicon(r))
	})
}
