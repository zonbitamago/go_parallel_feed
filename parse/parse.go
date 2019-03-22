package parse

import "github.com/mmcdole/gofeed"

// FeedParse feed情報取得処理
func FeedParse(feedURL string) (*gofeed.Feed, bool) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(feedURL)
	result := true
	if err != nil {
		result = false
	}
	return feed, result
}
