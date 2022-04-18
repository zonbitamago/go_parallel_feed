package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/mmcdole/gofeed"
	// parse "github.com/zonbitamago/go_parallel_feed/parse"
	"go_parallel_feed/parse"
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

	var wg sync.WaitGroup
	for _, url := range requestJSON.Urls {
		wg.Add(1)
		go func(url URL) {
			feed, result,err := parse.FeedParse(url.URL)
			if err != nil {
				fmt.Println(err)
			}

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
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set( "Access-Control-Allow-Methods","GET, POST, PUT, DELETE, OPTIONS" )
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(feedResults); err != nil {
		panic(err)
	}

}
