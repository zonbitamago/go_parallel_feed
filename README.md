# go_parallel_feed

[https://goparallelfeed.now.sh](https://goparallelfeed.now.sh)

Web API for parallel feed parse by Golang.

Support parsing both [RSS](https://en.wikipedia.org/wiki/RSS) and [Atom](https://en.wikipedia.org/wiki/Atom_(Web_standard)) feeds.

## Features

***Supported feed types:***

- RSS 0.90
- Netscape RSS 0.91
- Userland RSS 0.91
- RSS 0.92
- RSS 0.93
- RSS 0.94
- RSS 1.0
- RSS 2.0
- Atom 0.3
- Atom 1.0

## Usage

### Sample Request

```sh
curl https://goparallelfeed.now.sh -X POST -H 'Content-Type:application/json' -d '{"urls":[{"url":"https://feedforall.com/sample-feed.xml"},{"url":"https://feedforall.com/blog-feed.xml"}]}'
```

### Sample Response

```json
{
    "results":[
        {
            "result": true,
            "url": "https://feedforall.com/sample-feed.xml",
            "feed":{...}
        },
        {
            "result": true,
            "url": "https://feedforall.com/blog-feed.xml",
            "feed":{...}
        }

    ]
}
```