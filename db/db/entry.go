package db

import (
	"fmt"
	"time"
)

type EntryDetail struct {
	EntrySerial    uint64    `db:"serial"`
	EntryFeedID    uint64    `db:"feed_id"`
	Title          string    `db:"title"`
	Description    string    `db:"description"`
	PubDate        time.Time `db:"pubdate"`
	ReadFlag       ReadFlag  `db:"readflag"`
	URL            string    `db:"url"`
	SubscriptionID uint64    `db:"subscription_id"`
	SiteTitle      string    `db:"site_title"`
}

func (c *UserClient) UnreadEntryByCategory(categoryID uint64) ([]*EntryDetail, error) {
	e := []*EntryDetail{}
	err := c.Select(&e, `
SELECT
    entry.serial,
    entry.feed_id,
    story.title,
    description,
    entry.pubdate,
    readflag,
    story.url,
    subscription_id,
    feed.title AS site_title
FROM entry
INNER JOIN subscription ON subscription_id = subscription.id
INNER JOIN feed ON subscription.feed_id = feed.id
INNER JOIN story ON story.serial = entry.serial AND story.feed_id = entry.feed_id
WHERE subscription.category_id = ?
    AND readflag <> 1
    AND entry.user_id = ?
ORDER BY entry.pubdate DESC
`, categoryID, c.UserID)

	if err != nil {
		return nil, err
	}
	return e, nil
}

func (c *UserClient) UpdateEntrySeen(feedID, serial uint64) error {
	_, err := c.Exec(`
UPDATE entry
SET
    readflag = 1,
    update_at = CURRENT_TIMESTAMP
WHERE readflag = 0
    AND user_id = ?
    AND feed_id = ?
    AND serial = ?
    `, c.UserID, feedID, serial)
	return err
}

func (c *ClientTxn) InsertEntry(userID, feedID, serial, subscriptionID uint64, pubdate time.Time) error {
	_, err := c.Exec(`
INSERT INTO entry
(user_id, feed_id, serial, subscription_id, pubdate, readflag,  update_at)
VALUES (?, ?, ?, ?, ?, 0, CURRENT_TIMESTAMP)
    `, userID, feedID, serial, subscriptionID, pubdate)
	return err
}

func (c *ClientTxn) ExistEntry(feedID, serial uint64) (uint64, error) {
	var count uint64
	err := c.Get(&count, "SELECT COUNT(*) FROM entry WHERE feed_id = ? AND serial = ?", feedID, serial)
	return count, err
}

func (c *Client) PurgeReadEntry() error {
	var err error
	switch c.DriverName() {
	case "sqlite3":
		_, err = c.Exec("DELETE FROM entry WHERE readflag = 1 AND update_at < DATETIME('NOW', '-1 DAY')")
	case "mysql":
		_, err = c.Exec("DELETE FROM entry WHERE readflag = 1 AND update_at < DATE_ADD(CURRENT_TIMESTAMP, INTERVAL -1 DAY)")
	default:
		err = fmt.Errorf("invalid DB Driver: %s", c.DriverName())
	}
	return err
}
