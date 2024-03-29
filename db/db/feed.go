package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Feed struct {
	ID         uint64    `db:"id"`
	URL        string    `db:"url"`
	SiteURL    string    `db:"siteurl"`
	Title      string    `db:"title"`
	HTTPStatus string    `db:"http_status"`
	Pubdate    time.Time `db:"pubdate"`
	Term       string    `db:"term"`
	Cache      CacheJson `db:"cache"`
	NextSerial uint64    `db:"next_serial"`
}

type CacheJson struct {
	Etag     string `json:"If-None-Match,omitempty"`
	Modified string `json:"If-Modified-Since,omitempty"`
}

func (pc *CacheJson) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		json.Unmarshal(v, &pc)
		return nil
	case string:
		json.Unmarshal([]byte(v), &pc)
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

func (pc CacheJson) Value() (driver.Value, error) {
	return json.Marshal(pc)
}

func (f *Feed) UpdateTerm() {
	seconds := time.Now().Sub(f.Pubdate).Seconds()
	switch {
	case seconds > 86400*14:
		f.Term = "5"
	case seconds > 86400*7:
		f.Term = "4"
	case seconds > 86400*4:
		f.Term = "3"
	case seconds > 3600*12:
		f.Term = "2"
	default:
		f.Term = "1"
	}
}

func (c *ClientTxn) FeedByID(id uint64) (*Feed, error) {
	f := Feed{}
	if err := c.Get(&f, c.sql("SELECT * FROM feed WHERE id = ?"), id); err != nil {
		return nil, err
	}
	return &f, nil
}

func (c *Client) FeedsByID(ids []uint64) ([]*Feed, error) {
	sql, params, err := sqlx.In(`SELECT * FROM feed WHERE id IN (?)`, ids)
	if err != nil {
		return nil, err
	}

	feeds := []*Feed{}
	if err := c.Select(&feeds, c.sql(sql), params...); err != nil {
		return nil, err
	}
	return feeds, nil
}

func (c *Client) FeedsByTerm(term uint64) ([]*Feed, error) {
	f := []*Feed{}
	err := c.Select(&f, c.sql("SELECT * FROM feed WHERE term = ?"), term)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (c *Client) Feeds() ([]*Feed, error) {
	f := []*Feed{}
	err := c.Select(&f, "SELECT * FROM feed")
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (c *UserClientTxn) FeedByUrl(feedUrl, siteUrl string) (*Feed, error) {
	f := Feed{}
	if err := c.Get(&f, c.sql("SELECT * FROM feed WHERE url = ? AND siteurl = ?"), feedUrl, siteUrl); err != nil {
		return nil, err
	}
	return &f, nil
}

func (c *UserClientTxn) InsertFeed(feedUrl, siteUrl, title string) error {
	_, err := c.Exec(c.sql(`
INSERT INTO feed 
    (url, siteurl, title, http_status, pubdate, cache)
VALUES
    (?, ?, ?, 0, CURRENT_TIMESTAMP, '{}')
    `), feedUrl, siteUrl, title)
	return err
}

func (c *ClientTxn) GetNextSerial(feedID uint64) (*uint64, error) {
	var f uint64
	if err := c.Get(&f, c.sql("SELECT next_serial FROM feed WHERE id = ?"), feedID); err != nil {
		return nil, err
	}
	return &f, nil
}

func (c *ClientTxn) UpdateNextSerial(feedID uint64) error {
	_, err := c.Exec(c.sql("UPDATE feed SET next_serial = next_serial + 1 WHERE id = ?"), feedID)
	return err
}

func (c *ClientTxn) UpdateFeed(item Feed) error {
	_, err := c.Exec(c.sql("UPDATE feed SET http_status = ?, term = ?, cache = ? WHERE id = ?"), item.HTTPStatus, item.Term, item.Cache, item.ID)
	return err
}

func (c *ClientTxn) UpdateFeedRSSUrl(item Feed) error {
	_, err := c.Exec(c.sql("UPDATE feed SET url = ? WHERE id = ?"), item.URL, item.ID)
	return err
}

func (c *ClientTxn) UpdateFeedWithPubdate(item Feed) error {
	_, err := c.Exec(c.sql(`
UPDATE feed
SET
    http_status = ?, pubdate = ?, term = ?, cache = ?
WHERE
    id = ?
`), item.HTTPStatus, item.Pubdate, item.Term, item.Cache, item.ID)

	return err
}

func (c *ClientTxn) DeleteFeed(feedID uint64) error {
	_, err := c.Exec(c.sql("DELETE FROM feed WHERE id = ?"), feedID)
	return err
}
