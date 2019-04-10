package feeds

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/laurentlbm/comics-rss/feed"
)

func main() {
	var port = flag.Int64("port", 20480, "port to run the server on")
	flag.Parse()

	for comic, transform := range feed.Feeds {
		transformFeed(comic, transform)
	}

	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}

func transformFeed(comicURL string, transform feed.Transform) {
	a := make(chan struct{})
	go transform.Run(a)
	http.HandleFunc(fmt.Sprintf("/%s/", comicURL), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "text/xml")
		w.Write([]byte(transform.Generate()))
	})
	http.HandleFunc(fmt.Sprintf("/%s/favicon.ico", comicURL), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "image/png")
		w.Write(transform.GetFavicon(nil))
	})
}
