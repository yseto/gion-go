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

func walkOutlines(c echo.Context, categoryName string, o []opml.Outline) error {
	for i := range o {
		if o[i].Type == "rss" {
			if err := saveOutline(c, categoryName, o[i]); err != nil {
				return err
			}
		}
		if err := walkOutlines(c, o[i].Title, o[i].Outlines); err != nil {
			return err
		}
	}
	return nil
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
