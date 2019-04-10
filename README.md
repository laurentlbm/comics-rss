# comics-rss

The missing image rss feed for [https://explosm.net](https://www.buttersafe.com), [https://www.buttersafe.com](https://www.buttersafe.com) and [https://www.fowllanguagecomics.com](https://www.fowllanguagecomics.com). This service downloads the rss feed and replaces the description of the content with the image instead of a link.

Forked from [https://github.com/daniellowtw/explosm-rss](https://github.com/daniellowtw/explosm-rss)

## Installation

`go get github.com/laurentlbm/comics-rss`

### On appengine

This software is compatible with [AppEngine](https://cloud.google.com/appengine/docs/go/quickstart).

To use with AppEngine:

* Make a directory with necessary `app.yaml` file (or use the one provided)
* Clone this repo inside that directory
* Download the [Go SDK](https://cloud.google.com/appengine/downloads#Google_App_Engine_SDK_for_Go) and run `goapp serve` to make sure it is working
* Upload it with `appcfg.py`
  * `python appcfg.py update "<path-to-appengine-folder>" -A <app-name> -V <version-number>`

## Usage

```
go get .
go build
./comics-rss
```

Go to [http://localhost:20480/buttersafe](http://localhost:20480/buttersafe)

### Configuration

* `port` - the port that the server is listening on
* `refresh_interval` - how often to poll the actual feeds

## License

[MIT License](http://choosealicense.com/licenses/mit/)