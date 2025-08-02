package client

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

type Response struct {
	StatusCode int
	Location   string
	Cache      Cache
	Content    string
}

type Cache struct {
	Etag          string `json:"If-None-Match,omitempty"`
	Modified      string `json:"If-Modified-Since,omitempty"`
	RedirectCount int    `json:"-"`
	Location      string
}

var client *http.Client
var transport *http.Transport

func init() {
	transport = http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		Control:   control,
	}).DialContext

	client = &http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func Get(url string) (*http.Response, error) {
	return (&http.Client{
		Transport: transport,
	}).Get(url)
}

var ErrRedirectURL = errors.New("cant get redirect URL")

// リダイレクト先URLを取得する
func GetRedirectURL(url string) (string, error) {
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode == 301 || resp.StatusCode == 308 {
		url, err := resp.Location()
		if err == nil {
			return GetRedirectURL(url.String())
		}
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return url, nil
	}
	return "", ErrRedirectURL
}

func GetWithCache(ctx context.Context, xmlUrl string, cache Cache) (*Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", xmlUrl, nil)
	if err != nil {
		return nil, err
	}
	if cache.Etag != "" {
		req.Header.Add("If-None-Match", cache.Etag)
	}
	if cache.Modified != "" {
		req.Header.Add("If-Modified-Since", cache.Modified)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	u := &Response{}
	u.StatusCode = resp.StatusCode
	// リダイレクト先をリストア
	u.Location = cache.Location

	if etag := resp.Header.Get("Etag"); etag != "" {
		u.Cache.Etag = etag
	}
	if lm := resp.Header.Get("Last-Modified"); lm != "" {
		u.Cache.Modified = lm
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	u.Content = string(body)

	if code := u.StatusCode; code == 200 || code == 304 || code == 400 {
		return u, nil
	}

	if code := u.StatusCode; code == 301 || code == 302 || code == 303 || code == 307 || code == 308 {
		redirectLoaction, err := url.Parse(resp.Header.Get("Location"))
		if err != nil {
			return nil, err
		}
		parsedUrl, err := url.Parse(xmlUrl)
		if err != nil {
			return nil, err
		}
		nextLocation := parsedUrl.ResolveReference(redirectLoaction).String()
		// 301 は URL更新が必要
		if code == 301 {
			// リダイレクト先が取得できるか評価する
			redResp, err := client.Get(nextLocation)

			// リダイレクト先のパスがおかしい場合に拾う
			if err != nil {
				u.Cache = Cache{}
				u.StatusCode = 400
				return u, nil
			}
			defer func() {
				io.Copy(io.Discard, redResp.Body)
				redResp.Body.Close()
			}()

			if redResp.StatusCode >= 200 && redResp.StatusCode < 300 {
				cache.Location = nextLocation
			}
		}

		// リダイレクトループを検出する
		cache.RedirectCount++
		if cache.RedirectCount > 5 {
			u.Cache = Cache{}
			u.StatusCode = 400
			return u, nil
		}

		return GetWithCache(ctx, nextLocation, cache)
	}

	return u, nil
}
