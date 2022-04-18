package handler

import (
	"bytes"
	"encoding/json"
	"go_parallel_feed/testutils"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"testing"

	"github.com/mmcdole/gofeed"
)

var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	Handler(w, r)
})

func callIndexHandler(t *testing.T, jsonStr string) FeedResults {
	ts := httptest.NewServer(handler)
	defer ts.Close()

	resp, err := http.Post(ts.URL, "application/json", bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		t.Fatalf("Error by http.Get(). %v", err)
	}

	var got FeedResults
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("Error by json.NewDecoder(). %v", err)
	}
	return got
}

func TestIndexHandler(t *testing.T) {
	t.Run("request_http_success", func(t *testing.T) {

		data := `{"urls":[{"url":"https://feedforall.com/sample-feed.xml"}]}`

		got := callIndexHandler(t, data)

		var correctURLFeed gofeed.Feed
		json.Unmarshal([]byte(sampleFeedJSON), &correctURLFeed)

		want := FeedResults{
			Results: []FeedResult{
				FeedResult{
					Result: true,
					URL:    "https://feedforall.com/sample-feed.xml",
					Feed:   &correctURLFeed,
				},
			},
		}

		if !reflect.DeepEqual(want, got) {
			testutils.ErrorfHandler(t, want, got)
		}
	})

	t.Run("request_single", func(t *testing.T) {

		data := `{"urls":[{"url":"a"}]}`

		got := callIndexHandler(t, data)

		want := FeedResults{
			Results: []FeedResult{
				FeedResult{
					Result: false,
					URL:    "a",
				},
			},
		}

		if !reflect.DeepEqual(want, got) {
			testutils.ErrorfHandler(t, want, got)
		}
	})

	t.Run("request_multi", func(t *testing.T) {

		data := `{"urls":[{"url":"a"},{"url":"b"}]}`

		got := callIndexHandler(t, data)

		want := FeedResults{
			Results: []FeedResult{
				FeedResult{
					Result: false,
					URL:    "a",
				},
				{
					Result: false,
					URL:    "b",
				},
			},
		}

		// goroutineで順序が保証されないので、比較のために一旦ソートを行う。
		sort.SliceStable((got.Results), func(i, j int) bool { return got.Results[i].URL < got.Results[j].URL })

		if !reflect.DeepEqual(want, got) {
			testutils.ErrorfHandler(t, want, got)
		}

	})
}

var sampleFeedJSON = `
{
	"title": "Sample Feed - Favorite RSS Related Software & Resources",
	"description": "Take a look at some of FeedForAll's favorite software and resources for learning more about RSS.",
	"link": "http://www.feedforall.com",
	"updated": "Mon, 1 Nov 2004 13:17:17 -0500",
	"updatedParsed": "2004-11-01T18:17:17Z",
	"published": "Tue, 26 Oct 2004 14:06:44 -0500",
	"publishedParsed": "2004-10-26T19:06:44Z",
	"author": {
		"email": "marketing@feedforall.com"
	},
	"language": "en-us",
	"image": {
		"url": "http://www.feedforall.com/feedforall-temp.gif",
		"title": "FeedForAll Sample Feed"
	},
	"copyright": "Copyright 2004 NotePage, Inc.",
	"generator": "FeedForAll Beta1 (0.0.1.8)",
	"categories": [
		"Computers/Software/Internet/Site Management/Content Management"
	],
	"items": [
		{
			"title": "RSS Resources",
			"description": "Be sure to take a look at some of our favorite RSS Resources<br>\r\n<a href=\"http://www.rss-specifications.com\">RSS Specifications</a><br>\r\n<a href=\"http://www.blog-connection.com\">Blog Connection</a><br>\r\n<br>",
			"link": "http://www.feedforall.com",
			"published": "Tue, 26 Oct 2004 14:01:01 -0500",
			"publishedParsed": "2004-10-26T19:01:01Z"
		},
		{
			"title": "Recommended Desktop Feed Reader Software",
			"description": "<b>FeedDemon</b> enables you to quickly read and gather information from hundreds of web sites - without having to visit them. Don't waste any more time checking your favorite web sites for updates. Instead, use FeedDemon and make them come to you. <br>\r\nMore <a href=\"http://store.esellerate.net/a.asp?c=1_SKU5139890208_AFL403073819\">FeedDemon Information</a>",
			"link": "http://www.feedforall.com/feedforall-partners.htm",
			"published": "Tue, 26 Oct 2004 14:03:25 -0500",
			"publishedParsed": "2004-10-26T19:03:25Z"
		},
		{
			"title": "Recommended Web Based Feed Reader Software",
			"description": "<b>FeedScout</b> enables you to view RSS/ATOM/RDF feeds from different sites directly in Internet Explorer. You can even set your Home Page to show favorite feeds. Feed Scout is a plug-in for Internet Explorer, so you won't have to learn anything except for how to press 2 new buttons on Internet Explorer toolbar. <br>\r\nMore <a href=\"http://www.bytescout.com/feedscout.html\">Information on FeedScout</a><br>\r\n<br>\r\n<br>\r\n<b>SurfPack</b> can feature search tools, horoscopes, current weather conditions, LiveJournal diaries, humor, web modules and other dynamically updated content. <br>\r\nMore <a href=\"http://www.surfpack.com/\">Information on SurfPack</a><br>",
			"link": "http://www.feedforall.com/feedforall-partners.htm",
			"published": "Tue, 26 Oct 2004 14:06:44 -0500",
			"publishedParsed": "2004-10-26T19:06:44Z"
		}
	],
	"feedType": "rss",
	"feedVersion": "2.0"
}
`
