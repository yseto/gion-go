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

func OpmlExport(c echo.Context) error {
	db := c.(*CustomContext).DBUser()
	cats, err := db.Category()
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	o := opml.Body{}
	for i := range cats {
		feeds, err := db.FeedsByCategoryID(cats[i].ID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, nil)
		}

		b := opml.Outline{Text: cats[i].Name, Title: cats[i].Name}
		for j := range feeds {
			b.Outlines = append(b.Outlines, opml.Outline{
				Type:    "rss",
				Text:    feeds[j].Title,
				Title:   feeds[j].Title,
				HTMLURL: feeds[j].SiteURL,
				XMLURL:  feeds[j].URL})
		}
		o.Outlines = append(o.Outlines, b)
	}

	xml, err := opml.OPML{Version: "1.0", Head: opml.Head{Title: "export data"}, Body: o}.XML()
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	return c.JSON(http.StatusOK, opmlResult{xml})
}

func categoryByName(c echo.Context, categoryName string) (*db.Category, error) {
	db := c.(*CustomContext).DBUser()
	tx := db.MustBegin()
	cat, err := tx.CategoryByName(categoryName)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return nil, err
	}
	if cat != nil {
		tx.Commit()
		return cat, nil
	}

	err = tx.InsertCategory(categoryName)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	cat, err = tx.CategoryByName(categoryName)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return cat, nil
}

func saveOutline(c echo.Context, categoryName string, o opml.Outline) error {
	if o.XMLURL == "" || o.HTMLURL == "" || o.Title == "" {
		fmt.Printf("missing required parameter : %+v\n", o)
		return nil
	}

	category, err := categoryByName(c, categoryName)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	feed, err := insertFeed(c, o.XMLURL, o.HTMLURL, o.Title)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	tx := c.(*CustomContext).DBUser().MustBegin()

	sub, err := tx.SubscriptionByFeedID(feed.ID)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return c.JSON(http.StatusBadRequest, nil)
	}
	if sub != nil {
		fmt.Printf("already registered : %s\n", o.Title)
		tx.Commit()
		return nil
	}

	if err = tx.InsertSubscription(feed.ID, category.ID); err != nil {
		tx.Rollback()
		return err
	}
	fmt.Printf("registered : %s\n", o.Title)
	tx.Commit()

	return nil
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
