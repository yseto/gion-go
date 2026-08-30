package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Feed struct {
	ID            uint64       `db:"id"`
	URL           string       `db:"url"`
	SiteURL       string       `db:"siteurl"`
	Title         string       `db:"title"`
	HTTPStatus    string       `db:"http_status"`
	Pubdate       time.Time    `db:"pubdate"`
	Term          string       `db:"term"`
	Cache         CacheJson    `db:"cache"`
	NextSerial    uint64       `db:"next_serial"`
	RevisitAt     sql.NullTime `db:"revisit_at"`
	LatestEntryAt sql.NullTime `db:"latest_entry_at"`
}

func (f *Feed) GetRevisitAt() *time.Time {
	if f.RevisitAt.Valid {
		return &f.RevisitAt.Time
	}
	return nil
}

// no implemented for nil
func (f *Feed) SetRevisitAt(ts time.Time) {
	f.RevisitAt.Time, f.RevisitAt.Valid = ts, true
}

func (f *Feed) GetLatestEntryAt() *time.Time {
	if f.LatestEntryAt.Valid {
		return &f.LatestEntryAt.Time
	}
	return nil
}

// no implemented for nil
func (f *Feed) SetLatestEntryAt(ts time.Time) {
	f.LatestEntryAt.Time, f.LatestEntryAt.Valid = ts, true
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

func (c *ClientTxn) FeedByID(ctx context.Context, id uint64) (*Feed, error) {
	f := Feed{}
	if err := c.GetContext(ctx, &f, c.sql("SELECT * FROM feed WHERE id = ?"), id); err != nil {
		return nil, err
	}
	return &f, nil
}

func (c *Client) FeedsByID(ctx context.Context, ids []uint64) ([]*Feed, error) {
	sql, params, err := sqlx.In(`SELECT * FROM feed WHERE id IN (?)`, ids)
	if err != nil {
		return nil, err
	}

	feeds := []*Feed{}
	if err := c.SelectContext(ctx, &feeds, c.sql(sql), params...); err != nil {
		return nil, err
	}
	return feeds, nil
}

func (c *Client) Feeds(ctx context.Context) ([]*Feed, error) {
	f := []*Feed{}
	err := c.SelectContext(ctx, &f, "SELECT * FROM feed")
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (c *ClientTxn) FeedByUrl(ctx context.Context, feedUrl, siteUrl string) (*Feed, error) {
	f := Feed{}
	if err := c.GetContext(ctx, &f, c.sql("SELECT * FROM feed WHERE url = ? AND siteurl = ?"), feedUrl, siteUrl); err != nil {
		return nil, err
	}
	return &f, nil
}

func (c *ClientTxn) InsertFeed(ctx context.Context, feedUrl, siteUrl, title string) error {
	_, err := c.ExecContext(ctx, c.sql(`
INSERT INTO feed 
    (url, siteurl, title, http_status, pubdate, cache)
VALUES
    (?, ?, ?, 0, CURRENT_TIMESTAMP, '{}')
    `), feedUrl, siteUrl, title)
	return err
}

func (c *ClientTxn) GetNextSerial(ctx context.Context, feedID uint64) (*uint64, error) {
	var f uint64
	if err := c.GetContext(ctx, &f, c.sql("SELECT next_serial FROM feed WHERE id = ?"), feedID); err != nil {
		return nil, err
	}
	return &f, nil
}

func (c *ClientTxn) UpdateNextSerial(ctx context.Context, feedID uint64) error {
	_, err := c.ExecContext(ctx, c.sql("UPDATE feed SET next_serial = next_serial + 1 WHERE id = ?"), feedID)
	return err
}

func (c *ClientTxn) UpdateFeed(ctx context.Context, item Feed) error {
	_, err := c.ExecContext(ctx, c.sql("UPDATE feed SET http_status = ?, revisit_at = ?, cache = ? WHERE id = ?"), item.HTTPStatus, item.GetRevisitAt(), item.Cache, item.ID)
	return err
}

func (c *ClientTxn) UpdateFeedRSSUrl(ctx context.Context, item Feed) error {
	_, err := c.ExecContext(ctx, c.sql("UPDATE feed SET url = ? WHERE id = ?"), item.URL, item.ID)
	return err
}

func (c *ClientTxn) UpdateFeedSiteUrl(ctx context.Context, item Feed) error {
	_, err := c.ExecContext(ctx, c.sql("UPDATE feed SET siteurl = ? WHERE id = ?"), item.SiteURL, item.ID)
	return err
}

func (c *ClientTxn) UpdateFeedWithLatestEntryAt(ctx context.Context, item Feed) error {
	_, err := c.ExecContext(ctx, c.sql(`
UPDATE feed
SET
    http_status = ?, latest_entry_at = ?, revisit_at = ?, cache = ?
WHERE
    id = ?
`), item.HTTPStatus, item.GetLatestEntryAt(), item.GetRevisitAt(), item.Cache, item.ID)

	return err
}

func (c *ClientTxn) DeleteFeed(ctx context.Context, feedID uint64) error {
	_, err := c.ExecContext(ctx, c.sql("DELETE FROM feed WHERE id = ?"), feedID)
	return err
}
