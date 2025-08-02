package db

import (
	"context"
	"errors"
)

type Subscription struct {
	ID         uint64 `db:"id"`
	CategoryID uint64 `db:"category_id"`
	FeedID     uint64 `db:"feed_id"`
	UserID     uint64 `db:"user_id"`
}

func (c *UserClientTxn) SubscriptionByFeedID(ctx context.Context, feedID uint64) (*Subscription, error) {
	s := Subscription{}
	err := c.GetContext(ctx, &s, c.sql("SELECT * FROM subscription WHERE user_id = ? AND feed_id = ?"), c.UserID, feedID)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (c *ClientTxn) SubscriptionByFeedID(ctx context.Context, feedID uint64) ([]*Subscription, error) {
	s := []*Subscription{}
	err := c.SelectContext(ctx, &s, c.sql("SELECT * FROM subscription WHERE feed_id = ?"), feedID)
	if err != nil {
		return nil, err
	}
	return s, nil
}

type SubscriptionForUser struct {
	FeedID     uint64 `db:"id"`
	FeedTitle  string `db:"title"`
	CategoryID uint64 `db:"category_id"`
	HTTPStatus string `db:"http_status"`
	SiteURL    string `db:"siteurl"`
}

func (c *UserClient) Subscriptions(ctx context.Context) ([]*SubscriptionForUser, error) {
	s := []*SubscriptionForUser{}
	err := c.SelectContext(ctx, &s, c.sql(`
SELECT
    feed.id,
    feed.title,
    subscription.category_id,
    feed.http_status,
    feed.siteurl
FROM subscription
INNER JOIN feed ON feed_id = feed.id
WHERE subscription.user_id = ?
ORDER BY title ASC
`), c.UserID)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (c *UserClient) FeedsByCategoryID(ctx context.Context, categoryID uint64) ([]*Feed, error) {
	f := []*Feed{}
	err := c.SelectContext(ctx, &f, c.sql(`
SELECT
    feed.* 
FROM subscription
INNER JOIN feed ON feed_id = feed.id
WHERE
    subscription.user_id = ?
AND
    subscription.category_id = ?
`), c.UserID, categoryID)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (c *UserClient) DeleteSubscription(ctx context.Context, feedID uint64) error {
	_, err := c.ExecContext(ctx, c.sql("DELETE FROM subscription WHERE feed_id = ? AND user_id = ?"), feedID, c.UserID)
	return err
}

var errCategoryNotOwned = errors.New("category is not owned")

func (c *UserClientTxn) InsertSubscription(ctx context.Context, feedID, categoryID uint64) error {
	var exists bool
	err := c.GetContext(ctx, &exists, c.sql("SELECT EXISTS(SELECT user_id FROM category WHERE id = ? AND user_id = ?)"), categoryID, c.UserID)
	if err != nil || !exists {
		return errCategoryNotOwned
	}

	_, err = c.ExecContext(ctx, c.sql("INSERT INTO subscription (category_id, feed_id, user_id) VALUES (?, ?, ?)"), categoryID, feedID, c.UserID)
	return err
}

func (c *UserClient) UpdateSubscription(ctx context.Context, feedID, categoryID uint64) error {
	var exists bool
	err := c.GetContext(ctx, &exists, c.sql("SELECT EXISTS(SELECT user_id FROM category WHERE id = ? AND user_id = ?)"), categoryID, c.UserID)
	if err != nil || !exists {
		return errCategoryNotOwned
	}

	_, err = c.ExecContext(ctx, c.sql("UPDATE subscription SET category_id = ? WHERE feed_id = ? AND user_id = ?"), categoryID, feedID, c.UserID)
	return err
}
