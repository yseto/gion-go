package db

import (
	"context"
	"time"
)

type Websub struct {
	ID        uint64    `db:"id"`
	TopicURL  string    `db:"topic_url"`
	ExpiredAt time.Time `db:"expired_at"`
}

func (c *ClientTxn) Websubs(ctx context.Context) ([]*Websub, error) {
	w := []*Websub{}
	if err := c.SelectContext(ctx, &w, "SELECT * FROM websubs"); err != nil {
		return nil, err
	}
	return w, nil
}

func (c *ClientTxn) WebsubByID(ctx context.Context, id uint64) (*Websub, error) {
	w := Websub{}
	if err := c.GetContext(ctx, &w, "SELECT * FROM websubs WHERE id = ?", id); err != nil {
		return nil, err
	}
	return &w, nil
}

func (c *ClientTxn) WebsubByTopicURL(ctx context.Context, topicURL string) (*Websub, error) {
	w := Websub{}
	if err := c.GetContext(ctx, &w, "SELECT * FROM websubs WHERE topic_url = ?", topicURL); err != nil {
		return nil, err
	}
	return &w, nil
}

func (c *ClientTxn) InsertWebsub(ctx context.Context, topicURL string, expiredAt time.Time) error {
	_, err := c.ExecContext(ctx, "INSERT INTO websubs (id, topic_url, expired_at) VALUES (NULL, ?, ?)", topicURL, expiredAt)
	return err
}

func (c *ClientTxn) UpdateWebsub(ctx context.Context, item Websub) error {
	_, err := c.NamedExecContext(ctx, `UPDATE websubs SET expired_at = :expired_at WHERE id = :id`, item)
	return err
}

func (c *ClientTxn) DeleteWebsub(ctx context.Context, item Websub) error {
	_, err := c.ExecContext(ctx, `DELETE FROM websubs WHERE id = ?`, item.ID)
	return err
}
