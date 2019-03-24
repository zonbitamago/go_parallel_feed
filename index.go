package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/mmcdole/gofeed"
	parse "github.com/zonbitamago/go_parallel_feed/parse"
)

// FeedResults フィード結果json
type FeedResults struct {
	Results []FeedResult `json:"results"`
}

// FeedResult 個別Feed結果
type FeedResult struct {
	Result bool         `json:"result"`
	URL    string       `json:"url"`
	Feed   *gofeed.Feed `json:"feed"`
}

// URL リクエスト/レスポンス用構造体
type URL struct {
	URL string `json:"url"`
}

// ReqestJSON リクエスト受信用構造体
type ReqestJSON struct {
	Urls []URL `json:"urls"`
}

// Handler nowデプロイ時のエントリーポイント
func Handler(w http.ResponseWriter, req *http.Request) {

	body, _ := ioutil.ReadAll(req.Body)

	var requestJSON ReqestJSON

	json.Unmarshal(body, &requestJSON)

	var feedResult []FeedResult

	// for _, url := range requestJSON.Urls {
	// 	feed, result := FeedParse(url.URL)

	// 	feedResult = append(feedResult, FeedResult{
	// 		Result: result,
	// 		URL:    URL{url.URL},
	// 		Feed:   feed,
	// 	})
	// }

	var wg sync.WaitGroup
	for _, url := range requestJSON.Urls {
		wg.Add(1)
		go func(url URL) {
			feed, result := parse.FeedParse(url.URL)

			feedResult = append(feedResult, FeedResult{
				Result: result,
				URL:    url.URL,
				Feed:   feed,
			})

			wg.Done()
		}(url)
	}
	wg.Wait()

	feedResults := FeedResults{feedResult}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(feedResults); err != nil {
		panic(err)
	}

}

// // FeedParse feed情報取得処理
// func FeedParse(feedURL string) (*gofeed.Feed, bool) {
// 	fp := gofeed.NewParser()
// 	feed, err := fp.ParseURL(feedURL)
// 	result := true
// 	if err != nil {
// 		result = false
// 	}
// 	return feed, result
// }
