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

func PinnedItems(c echo.Context) error {
	pin, err := c.(*CustomContext).DBUser().PinnedItems()
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	return c.JSON(http.StatusOK, pin)
}

func Profile(c echo.Context) error {
	pin, err := c.(*CustomContext).DBUser().Profile()
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	return c.JSON(http.StatusOK, pin)
}

func Categories(c echo.Context) error {
	cat, err := c.(*CustomContext).DBUser().Category()
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	return c.JSON(http.StatusOK, cat)
}

func CategoryAndUnreadEntryCount(c echo.Context) error {
	cat, err := c.(*CustomContext).DBUser().CategoryAndUnreadEntryCount()
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	return c.JSON(http.StatusOK, cat)
}

func UnreadEntry(c echo.Context) error {
	category, err := strconv.ParseUint(c.FormValue("category"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	db := c.(*CustomContext).DBUser()
	u, err := db.Profile()
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	cat, err := db.UnreadEntryByCategory(category)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	if u.EntryCount > 0 && len(cat) > int(u.EntryCount) {
		cat = cat[:u.EntryCount]
	}

	p := bluemonday.NewPolicy()
	for i := range cat {
		d := p.Sanitize(cat[i].Description)
		if u.SubstringLength > 0 && uint64(utf8.RuneCountInString(d)) > u.SubstringLength {
			d = string([]rune(d)[:u.SubstringLength])
		}
		cat[i].Description = d
	}

	return c.JSON(http.StatusOK, cat)
}

type categoryAndSubscription struct {
	ID   uint64                    `json:"id"`
	Name string                    `json:"name"`
	Subs []*db.SubscriptionForUser `json:"subscription"`
}

func Subscriptions(c echo.Context) error {
	dbClient := c.(*CustomContext).DBUser()
	subs, err := dbClient.Subscriptions()
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	cat, err := dbClient.Category()
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	var resp []categoryAndSubscription
	for i := range cat {
		var subsOnCategory []*db.SubscriptionForUser
		for j := range subs {
			if cat[i].ID == subs[j].CategoryID {
				subsOnCategory = append(subsOnCategory, subs[j])
			}
		}

		resp = append(resp, categoryAndSubscription{
			ID:   cat[i].ID,
			Name: cat[i].Name,
			Subs: subsOnCategory,
		})
	}

	return c.JSON(http.StatusOK, resp)
}

type asread struct {
	FeedID uint64 `json:"feed_id"`
	Serial uint64 `json:"serial"`
}
type asreadResult struct {
	Result bool `json:"result"`
}

func SetAsread(c echo.Context) error {
	var reads []asread
	err := json.NewDecoder(c.Request().Body).Decode(&reads)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	defer c.Request().Body.Close()

	//	return c.JSON(http.StatusOK, asreadResult{true}) // FOR DEBUG

	db := c.(*CustomContext).DBUser()
	for i := range reads {
		err := db.UpdateEntrySeen(reads[i].FeedID, reads[i].Serial)
		if err != nil {
			return c.JSON(http.StatusBadRequest, nil)
		}
	}
	return c.JSON(http.StatusOK, asreadResult{true})
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
