package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"unicode/utf8"

	"github.com/labstack/echo/v4"
	"github.com/microcosm-cc/bluemonday"

	"github.com/yseto/gion-go/db/db"
)

type categoryAndSubscription struct {
	ID   uint64                    `json:"id"`
	Name string                    `json:"name"`
	Subs []*db.SubscriptionForUser `json:"subscription"`
}

type asread struct {
	FeedID uint64 `json:"feed_id"`
	Serial uint64 `json:"serial"`
}
type asreadResult struct {
	Result bool `json:"result"`
}

type setPinResult struct {
	Readflag uint64 `json:"readflag"`
}

func SetPin(c echo.Context) error {
	rawReadflag, rErr := strconv.ParseUint(c.FormValue("readflag"), 10, 64)
	serial, sErr := strconv.ParseUint(c.FormValue("serial"), 10, 64)
	feedID, fErr := strconv.ParseUint(c.FormValue("feed_id"), 10, 64)
	if rErr != nil || sErr != nil || fErr != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	var readflag db.ReadFlag
	if rawReadflag == 2 {
		readflag = db.Seen
		rawReadflag = 1
	} else {
		readflag = db.SetPin
		rawReadflag = 2
	}

	fmt.Printf("PIN feed_id:%d\tserial:%d\treadflag:%d\n", feedID, serial, readflag)

	tx := c.(*CustomContext).DBUser().MustBegin()
	if tx.UpdateEntry(feedID, serial, readflag) != nil {
		tx.Rollback()
		return c.JSON(http.StatusBadRequest, nil)
	}
	tx.Commit()

	return c.JSON(http.StatusOK, setPinResult{rawReadflag})
}

type registerResult struct {
	Result string `json:"result"`
}

func RegisterCategory(c echo.Context) error {
	categoryName := c.FormValue("name")
	if categoryName == "" {
		return c.JSON(http.StatusBadRequest, nil)
	}

	tx := c.(*CustomContext).DBUser().MustBegin()

	cat, err := tx.CategoryByName(categoryName)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return c.JSON(http.StatusBadRequest, nil)
	}
	if cat != nil {
		tx.Commit()
		return c.JSON(http.StatusOK, registerResult{"ERROR_ALREADY_REGISTER"})
	}

	if err = tx.InsertCategory(categoryName); err != nil {
		tx.Rollback()
		return c.JSON(http.StatusBadRequest, nil)
	}
	tx.Commit()

	return c.JSON(http.StatusOK, registerResult{"OK"})
}

type registerSubscriptionResult struct {
	Result string `json:"result"`
}

func insertFeed(c echo.Context, rssUrl, siteUrl, title string) (*db.Feed, error) {
	tx := c.(*CustomContext).DBUser().MustBegin()
	feed, err := tx.FeedByUrl(rssUrl, siteUrl)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return nil, err
	}
	if feed != nil {
		tx.Commit()
		return feed, nil
	}

	err = tx.InsertFeed(rssUrl, siteUrl, title)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	feed, err = tx.FeedByUrl(rssUrl, siteUrl)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return feed, nil
}

func RegisterSubscription(c echo.Context) error {
	rssUrl, rErr := url.Parse(c.FormValue("rss"))
	siteUrl, sErr := url.Parse(c.FormValue("url"))
	title := c.FormValue("title")
	category, cErr := strconv.ParseUint(c.FormValue("category"), 10, 64)
	if rErr != nil || sErr != nil || title == "" || cErr != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	feed, err := insertFeed(c, rssUrl.String(), siteUrl.String(), title)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	db := c.(*CustomContext).DBUser()
	tx := db.MustBegin()
	sub, err := tx.SubscriptionByFeedID(feed.ID)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return c.JSON(http.StatusBadRequest, nil)
	}
	if sub != nil {
		tx.Rollback()
		return c.JSON(http.StatusOK, registerResult{"ERROR_ALREADY_REGISTER"})
	}

	cat, err := db.CategoryByID(category)
	if err != nil {
		tx.Rollback()
		return c.JSON(http.StatusBadRequest, nil)
	}
	if cat == nil {
		tx.Rollback()
		return c.JSON(http.StatusBadRequest, nil)
	}
	if tx.InsertSubscription(feed.ID, cat.ID) != nil {
		tx.Rollback()
		return c.JSON(http.StatusBadRequest, nil)
	}
	tx.Commit()

	return c.JSON(http.StatusOK, registerResult{"OK"})
}

type commonResult struct {
	Result string `json:"r"`
}

func DeleteSubscription(c echo.Context) error {
	deleteType := c.FormValue("subscription")
	id, err := strconv.ParseUint(c.FormValue("id"), 10, 64)
	if err != nil || deleteType == "" {
		return c.JSON(http.StatusBadRequest, nil)
	}

	db := c.(*CustomContext).DBUser()
	switch deleteType {
	case "category":
		err = db.DeleteCategory(id)
	case "entry":
		err = db.DeleteSubscription(id)
	default:
		err = fmt.Errorf("invalid type")
	}
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	return c.JSON(http.StatusOK, commonResult{"OK"})
}

func ChangeSubscription(c echo.Context) error {
	categoryID, cErr := strconv.ParseUint(c.FormValue("category"), 10, 64)
	feedID, iErr := strconv.ParseUint(c.FormValue("id"), 10, 64)
	if cErr != nil || iErr != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	if c.(*CustomContext).DBUser().UpdateSubscription(feedID, categoryID) != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	return c.JSON(http.StatusOK, commonResult{"OK"})
}

func UpdateProfile(c echo.Context) error {
	var profile db.UserProfile
	if err := json.NewDecoder(c.Request().Body).Decode(&profile); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	defer c.Request().Body.Close()

	if c.(*CustomContext).DBUser().UpdateProfile(profile) != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	return c.JSON(http.StatusOK, commonResult{"OK"})
}

func RemoveAllPin(c echo.Context) error {
	if c.(*CustomContext).DBUser().RemovePinnedItem() != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	return c.JSON(http.StatusOK, asreadResult{true})
}
