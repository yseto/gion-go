package db

import (
	"fmt"

	"github.com/yseto/gion-go/internal/pin"
)

type ReadFlag uint64

const (
	Unseen ReadFlag = iota
	Seen
	SetPin
)

func (c ReadFlag) String() string {
	switch c {
	case Unseen:
		return "Unseen"
	case Seen:
		return "Seen"
	case SetPin:
		return "Setpin"
	default:
		panic("Unknown")
	}
}

func (c ReadFlag) ToPinReadFlag() pin.ReadFlag {
	switch c {
	case Unseen:
		return pin.Unseen
	case Seen:
		return pin.Seen
	case SetPin:
		return pin.Setpin
	default:
		panic("Unknown")
	}
}

type PinnedItem struct {
	Title         string `db:"title"`
	URL           string `db:"url"`
	EntrySerial   uint64 `db:"serial"`
	EntryFeedID   uint64 `db:"feed_id"`
	EntryUpdateAt MyTime `db:"update_at"`
}

func (c *UserClient) PinnedItems() ([]*PinnedItem, error) {
	items := []*PinnedItem{}
	err := c.Select(&items, `
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
`, c.UserID)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (c *UserClientTxn) UpdateEntry(feedID, serial uint64, readflag ReadFlag) error {
	if readflag == Unseen {
		return fmt.Errorf("UpdateEntry : Readflag is invalid. UserID: %d feedID: %d Serial: %d", c.UserID, feedID, serial)
	}
	_, err := c.Exec(`
UPDATE entry
SET
    readflag = ?,
    update_at = CURRENT_TIMESTAMP
WHERE user_id = ? AND serial = ? AND feed_id = ?
    `, readflag, c.UserID, serial, feedID)
	return err
}

func (c *UserClient) RemovePinnedItem() error {
	_, err := c.Exec(`
UPDATE entry
SET
    readflag = 1,
    update_at = CURRENT_TIMESTAMP
WHERE readflag = 2 AND user_id = ?
    `, c.UserID)
	return err
}
