package db

import (
	"context"
	"fmt"
	"time"

	"github.com/yseto/gion-go/internal/pin"
)

type PinnedItem struct {
	Title         string    `db:"title"`
	URL           string    `db:"url"`
	EntrySerial   uint64    `db:"serial"`
	EntryFeedID   uint64    `db:"feed_id"`
	EntryUpdateAt time.Time `db:"update_at"`
}

func (c *UserClient) PinnedItems(ctx context.Context) ([]*PinnedItem, error) {
	items := []*PinnedItem{}
	err := c.SelectContext(ctx, &items, c.sql(`
SELECT
    story.title,
    story.url,
    entry.serial,
    entry.feed_id,
    entry.update_at
FROM entry
INNER JOIN
    story ON story.serial = entry.serial
AND
    story.feed_id = entry.feed_id
WHERE
    entry.readflag = 2
AND
    entry.user_id = ?
ORDER BY pubdate DESC
`), c.UserID)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (c *UserClientTxn) UpdateEntry(ctx context.Context, feedID, serial uint64, readflag pin.ReadFlag) error {
	if readflag == pin.Unseen {
		return fmt.Errorf("UpdateEntry : Readflag is invalid. UserID: %d feedID: %d Serial: %d", c.UserID, feedID, serial)
	}
	_, err := c.ExecContext(ctx, c.sql(`
UPDATE entry
SET
    readflag = ?,
    update_at = CURRENT_TIMESTAMP
WHERE user_id = ? AND serial = ? AND feed_id = ?
    `), readflag, c.UserID, serial, feedID)
	return err
}

func (c *UserClient) RemovePinnedItem(ctx context.Context) error {
	_, err := c.ExecContext(ctx, c.sql(`
UPDATE entry
SET
    readflag = 1,
    update_at = CURRENT_TIMESTAMP
WHERE readflag = 2 AND user_id = ?
    `), c.UserID)
	return err
}
