package feed

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

type Transform struct {
	FeedURL         string
	FaviconURL      string
	ImageRegexp     *regexp.Regexp
	ExtraRegexp     *regexp.Regexp
	RefreshInterval time.Duration

	data    []byte
	rssData channel
}

func (t *Transform) Run(abort chan struct{}) {
	// Execute it the first time
	t.Do(nil)
	for {
		select {
		case <-abort:
			return
		case <-time.After(t.RefreshInterval):
			t.Do(nil)
		}
	}
}

func (t *Transform) Do(r *http.Request) error {
	fp := gofeed.NewParser()
	if r != nil {
		fp.Client = urlfetch.Client(appengine.NewContext(r))
	}
	feed, err := fp.ParseURL(t.FeedURL)
	if err != nil {
		return err
	}
	var is []Item
	for _, i := range feed.Items {
		data, lastURLQuery, err := getDataFromNet(i.Link, r)
		if err != nil {
			safeLog(r, "%s", err)
		}
		image := findMatch(data, t.ImageRegexp)
		if image == "" {
			safeLog(r, "cannot find image. Data is: %s", t.data)
		}
		if isValidURL(image) {
			if strings.HasPrefix(image, "//") {
				image = "https:" + image
			}
			image = fmt.Sprintf("<img alt=\"\" src=\"%s\" />", image)
		}
		extraElements := findMatch(data, t.ExtraRegexp)
		is = append(is, Item{
			Title:       i.Title,
			Link:        lastURLQuery,
			Description: CData{fmt.Sprintf("%s %s", image, extraElements)},
			Category:    i.Categories,
			Guid:        i.GUID,
			PubDate:     i.Published,
		})
	}
	t.rssData = channel{
		Title:       feed.Title,
		Description: feed.Description,
		Link:        feed.Link,
		Item:        is,
		Image: Image{
			URL:   t.FaviconURL,
			Link:  feed.Link,
			Title: feed.Title,
		},
	}
	return nil
}

func getDataFromNet(url string, r *http.Request) ([]byte, string, error) {
	httpClient := http.DefaultClient
	if r != nil {
		httpClient = urlfetch.Client(appengine.NewContext(r))
	}

	lastURLQuery := url
	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {

		if len(via) > 10 {
			return errors.New("too many redirects")
		}
		lastURLQuery = req.URL.String()
		return nil
	}

	res, err := httpClient.Get(url)
	if err != nil {
		return nil, lastURLQuery, err
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, lastURLQuery, err
	}
	return data, lastURLQuery, nil
}

func findMatch(data []byte, matchRegexp *regexp.Regexp) string {
	if matchRegexp == nil {
		return ""
	}
	matches := matchRegexp.FindSubmatch(data)
	if len(matches) > 1 {
		return string(matches[1])
	}
	return ""
}

func isValidURL(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	return true
}

func (t *Transform) GetFavicon(r *http.Request) []byte {
	f, _, err := getDataFromNet(t.FaviconURL, r)
	if err != nil {
		return nil
	}
	return f
}

func (t *Transform) Generate() string {
	if len(t.rssData.Item) == 0 {
		return "Please try again. There was an error retrieving the feeds: no feeds."
	}
	return generate(t.rssData)
}

// safeLog only logs if request is not nil
func safeLog(r *http.Request, format string, args ...interface{}) {
	if r == nil {
		return
	}
	log.Errorf(appengine.NewContext(r), format, args)
}
