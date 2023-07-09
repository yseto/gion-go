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
