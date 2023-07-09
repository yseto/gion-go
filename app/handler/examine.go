package handler

import (
	"net/http"
	"net/url"

	"github.com/antchfx/htmlquery"
	"github.com/hakobe/paranoidhttp"
	"github.com/labstack/echo/v4"
	"github.com/mmcdole/gofeed"
	"golang.org/x/net/html/charset"
)

type previewFeed struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Date  string `json:"date"`
}

type examine struct {
	Title       string         `json:"title"`
	URL         string         `json:"url"`
	PreviewFeed []*previewFeed `json:"preview_feed"`
}

func ExamineSubscription(c echo.Context) error {
	v := examineSubscriptionGetContent(c.FormValue("url"))
	if v.URL == "" {
		return c.JSON(http.StatusOK, v)
	}

	resp, err := paranoidhttp.DefaultClient.Get(v.URL)
	if err != nil {
		v.URL = ""
		return c.JSON(http.StatusOK, v)
	}
	defer resp.Body.Close()

	feed, err := gofeed.NewParser().Parse(resp.Body)
	if err != nil {
		v.URL = ""
		return c.JSON(http.StatusOK, v)
	}

	var pf []*previewFeed
	for i := range feed.Items {
		var date string
		if dt := feed.Items[i].PublishedParsed; dt != nil {
			date = dt.Format("01/02 15:04")
		}
		pf = append(pf, &previewFeed{
			Title: feed.Items[i].Title,
			URL:   feed.Items[i].Link,
			Date:  date,
		})

		if len(pf) == 3 {
			break
		}
	}
	v.PreviewFeed = pf

	return c.JSON(http.StatusOK, v)
}
