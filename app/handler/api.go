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

type registerResult struct {
	Result string `json:"result"`
}

type registerSubscriptionResult struct {
	Result string `json:"result"`
}

type commonResult struct {
	Result string `json:"r"`
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
