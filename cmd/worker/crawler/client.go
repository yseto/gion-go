package crawler

import (
	"io"
	"net/http"
	"net/url"

	"github.com/hakobe/paranoidhttp"

	"github.com/yseto/gion-go/db/db"
)

type userAgent struct {
	StatusCode    int
	Location      string
	Cache         db.CacheJson
	Content       string
	RedirectCount int
}

func (u *userAgent) Get(xmlUrl string, cache db.CacheJson) error {
	client := paranoidhttp.DefaultClient
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	req, err := http.NewRequest("GET", xmlUrl, nil)
	if err != nil {
		return err
	}
	if cache.Etag != "" {
		req.Header.Add("If-None-Match", cache.Etag)
	}
	if cache.Modified != "" {
		req.Header.Add("If-Modified-Since", cache.Modified)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	u.StatusCode = resp.StatusCode

	if etag := resp.Header.Get("Etag"); etag != "" {
		u.Cache.Etag = etag
	}
	if lm := resp.Header.Get("Last-Modified"); lm != "" {
		u.Cache.Modified = lm
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	u.Content = string(body)

	code := u.StatusCode
	if code == 200 || code == 304 || code == 400 {
		return nil
	}

	if code == 301 || code == 302 || code == 303 || code == 307 {
		redirectLoaction, err := url.Parse(resp.Header.Get("Location"))
		if err != nil {
			return err
		}
		parsedUrl, err := url.Parse(xmlUrl)
		if err != nil {
			return err
		}
		nextLocation := parsedUrl.ResolveReference(redirectLoaction).String()
		// 301 は URL更新が必要
		if code == 301 {
			// リダイレクト先が取得できるか評価する
			redResp, err := client.Get(nextLocation)
			// リダイレクト先のパスがおかしい場合に拾う
			if err != nil {
				var c db.CacheJson
				u.Cache = c
				u.StatusCode = 400
				return nil
			}
			defer func() {
				io.Copy(io.Discard, redResp.Body)
				redResp.Body.Close()
			}()
			if redResp.StatusCode >= 200 && redResp.StatusCode < 300 {
				u.Location = nextLocation
			}
		}

		// リダイレクトループを検出する
		u.RedirectCount++
		if u.RedirectCount > 5 {
			u.StatusCode = 400
			return nil
		}

		return u.Get(nextLocation, cache)
	}

	return nil
}
