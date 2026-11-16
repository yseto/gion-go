package client

import (
	"net/http"
	"os"
)

const envCrawlerUA = "CRAWLER_USER_AGENT"

type transportForUserAgent struct {
	Transport http.RoundTripper
}

func (t *transportForUserAgent) RoundTrip(req *http.Request) (*http.Response, error) {
	if ua, ok := os.LookupEnv(envCrawlerUA); ok {
		req.Header.Set("User-Agent", ua)
	}
	return t.Transport.RoundTrip(req)
}
