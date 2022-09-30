package db

import "fmt"

type Category struct {
	ID     uint64 `db:"id" json:"id"`
	UserID uint64 `db:"user_id" json:"-"`
	Name   string `db:"name" json:"name"`
}

func (c *UserClient) CategoryByID(id uint64) (*Category, error) {
	category := Category{}
	err := c.Get(&category, "SELECT * FROM category WHERE user_id = ? AND id = ?", c.UserID, id)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (c *UserClientTxn) CategoryByName(name string) (*Category, error) {
	category := Category{}
	err := c.Get(&category, "SELECT * FROM category WHERE user_id = ? AND name = ?", c.UserID, name)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (c *UserClient) Category() ([]*Category, error) {
	categories := []*Category{}
	err := c.Select(&categories, "SELECT * FROM category WHERE user_id = ? ORDER BY name ASC", c.UserID)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

type UnreadEntryCount struct {
	Count uint64 `db:"count" json:"count"`
	ID    uint64 `db:"id" json:"id"`
	Name  string `db:"name" json:"name"`
}

func (c *UserClient) CategoryAndUnreadEntryCount() ([]*UnreadEntryCount, error) {
	counts := []*UnreadEntryCount{}
	err := c.Select(&counts, `
SELECT
    COUNT(0) AS count,
    category.id AS id,
    category.name AS name
FROM entry
INNER JOIN subscription ON entry.subscription_id = subscription.id
INNER JOIN category ON subscription.category_id = category.id
WHERE readflag <> 1
    AND category.user_id = ?
GROUP BY category.id
ORDER BY category.name ASC
    `, c.UserID)
	if err != nil {
		return nil, err
	}
	return counts, nil
}

func (c *UserClient) DeleteCategory(id uint64) error {
	res, err := c.Exec("DELETE FROM category WHERE id = ? AND user_id = ?", id, c.UserID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("DeleteCategory : missing delete line. userid: %d, id:%d", c.UserID, id)
	}
	return nil
}

func (c *UserClientTxn) InsertCategory(name string) error {
	_, err := c.Exec("INSERT INTO category (user_id, name) VALUES (?, ?)", c.UserID, name)
	if err != nil {
		return err
	}
	return nil
}
