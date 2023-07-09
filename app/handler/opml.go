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

func OpmlImport(c echo.Context) error {
	o, err := opml.NewOPML([]byte(c.FormValue("xml")))
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	if err = walkOutlines(c, "Default Category", o.Outlines()); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	return c.JSON(http.StatusOK, opmlImportResult{true})
}
