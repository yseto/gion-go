package handler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gilliek/go-opml/opml"
	"github.com/labstack/echo/v4"

	"github.com/yseto/gion-go/db/db"
)

type opmlResult struct {
	Xml string `json:"xml"`
}

type opmlImportResult struct {
	Done bool `json:"done"`
}
