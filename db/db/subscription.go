package db

type Subscription struct {
	ID         uint64 `db:"id"`
	CategoryID uint64 `db:"category_id"`
	FeedID     uint64 `db:"feed_id"`
	UserID     uint64 `db:"user_id"`
}

func (c *UserClientTxn) SubscriptionByFeedID(feedID uint64) (*Subscription, error) {
	s := Subscription{}
	err := c.Get(&s, "SELECT * FROM subscription WHERE user_id = ? AND feed_id = ?", c.UserID, feedID)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (c *ClientTxn) SubscriptionByFeedID(feedID uint64) ([]*Subscription, error) {
	s := []*Subscription{}
	err := c.Select(&s, "SELECT * FROM subscription WHERE feed_id = ?", feedID)
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

func (c *UserClient) Subscriptions() ([]*SubscriptionForUser, error) {
	s := []*SubscriptionForUser{}
	err := c.Select(&s, `
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
`, c.UserID)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (c *UserClient) FeedsByCategoryID(categoryID uint64) ([]*Feed, error) {
	f := []*Feed{}
	err := c.Select(&f, `
SELECT
    feed.* 
FROM subscription
INNER JOIN feed ON feed_id = feed.id
WHERE
    subscription.user_id = ?
AND
    subscription.category_id = ?
`, c.UserID, categoryID)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (c *UserClient) DeleteSubscription(feedID uint64) error {
	_, err := c.Exec("DELETE FROM subscription WHERE feed_id = ? AND user_id = ?", feedID, c.UserID)
	return err
}

func (c *UserClientTxn) InsertSubscription(feedID, categoryID uint64) error {
	_, err := c.Exec("INSERT INTO subscription (category_id, feed_id, user_id) VALUES (?, ?, ?)", categoryID, feedID, c.UserID)
	return err
}

func (c *UserClient) UpdateSubscription(feedID, categoryID uint64) error {
	_, err := c.Exec("UPDATE subscription SET category_id = ? WHERE feed_id = ? AND user_id = ?", categoryID, feedID, c.UserID)
	return err
}
