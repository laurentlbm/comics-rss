package feed

import (
	"flag"
	"regexp"
	"time"
)

var (
	RefreshInterval = flag.Duration("refresh_interval", time.Hour*3, "interval to check for new data")
	Feeds           = map[string]Transform{
		"buttersafe": Transform{
			FeedURL:         "https://feeds.feedburner.com/buttersafe",
			FaviconURL:      "https://www.google.com/s2/favicons?domain=buttersafe.com",
			ImageRegexp:     regexp.MustCompile(`(?s)<div id="comic">(.*?)<\/div>`),
			ExtraRegexp:     regexp.MustCompile(`(?s)<div class="entry">(.*?)<\/div>`),
			RefreshInterval: *RefreshInterval,
		},
		"fowllanguage": Transform{
			FeedURL:         "http://www.fowllanguagecomics.com/feed/",
			FaviconURL:      "https://www.google.com/s2/favicons?domain=fowllanguagecomics.com",
			ImageRegexp:     regexp.MustCompile(`(?s)data-lazy-src="(.*?)\?fit=`),
			ExtraRegexp:     nil,
			RefreshInterval: *RefreshInterval,
		},
		"explosm": Transform{
			FeedURL: "http://feeds.feedburner.com/Explosm",
			FaviconURL:  "https://www.google.com/s2/favicons?domain=explosm.com",
			ImageRegexp:     regexp.MustCompile(`(?s)<div id="comic-wrap">(.*?)</div>`),
			ExtraRegexp:     nil,
			RefreshInterval: *RefreshInterval,
		},
	}
)
