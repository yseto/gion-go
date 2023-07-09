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
