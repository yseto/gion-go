package db

import (
	"context"
)

func (c *ClientTxn) InsertStory(ctx context.Context, feedID, serial uint64, title, description, url string) error {
	_, err := c.ExecContext(ctx, c.sql(`
INSERT INTO story
(feed_id, serial, title, description, url)
VALUES
(?, ?, ?, ?, ?)
`), feedID, serial, title, description, url)
	return err
}

type Story struct {
	FeedID uint64 `db:"feed_id"`
	Serial uint64 `db:"serial"`
	URL    string `db:"url"`
}

func (c *ClientTxn) Stories(ctx context.Context) ([]*Story, error) {
	s := []*Story{}
	err := c.SelectContext(ctx, &s, "SELECT feed_id, serial, url FROM story")
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (c *ClientTxn) DeleteStory(ctx context.Context, feedID, serial uint64) error {
	_, err := c.ExecContext(ctx, c.sql("DELETE FROM story WHERE feed_id = ? AND serial = ?"), feedID, serial)
	return err
}

func (c *ClientTxn) FindStory(ctx context.Context, feedID uint64, url string) (*Story, error) {
	s := Story{}
	err := c.GetContext(ctx, &s, c.sql(`SELECT feed_id, serial, url FROM story WHERE feed_id = ? AND url = ?`), feedID, url)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (c *ClientTxn) UpdateStory(ctx context.Context, s Story, title, description string) error {
	_, err := c.ExecContext(ctx, c.sql(`UPDATE story SET title = ?, description = ? WHERE feed_id = ? AND serial = ? AND url = ?`),
		title,
		description,
		s.FeedID,
		s.Serial,
		s.URL,
	)
	return err
}
