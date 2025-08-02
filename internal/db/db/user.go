package db

import (
	"context"
	"time"
)

type UserProfile struct {
	AutoSeen           bool   `db:"autoseen"`
	EntryCount         uint64 `db:"numentry"`
	OnLoginSkipPinList bool   `db:"nopinlist"`
	SubstringLength    uint64 `db:"numsubstr"`
}

func (c *UserClient) Profile(ctx context.Context) (*UserProfile, error) {
	p := UserProfile{}
	err := c.GetContext(ctx, &p, c.sql("SELECT autoseen, numentry, nopinlist, numsubstr FROM users WHERE id = ?"), c.UserID)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (c *UserClient) UpdateProfile(ctx context.Context, item UserProfile) error {
	_, err := c.ExecContext(ctx, c.sql("UPDATE users SET autoseen = ?, numentry = ?, nopinlist = ?, numsubstr = ? WHERE id = ?"),
		item.AutoSeen,
		item.EntryCount,
		item.OnLoginSkipPinList,
		item.SubstringLength,
		c.UserID)
	return err
}

type User struct {
	ID        uint64    `db:"id"`
	Name      string    `db:"name"`
	Digest    string    `db:"digest"`
	LastLogin time.Time `db:"last_login"`
	UserProfile
}

func (c *UserClient) User(ctx context.Context) (*User, error) {
	u := User{}
	err := c.GetContext(ctx, &u, c.sql("SELECT id, name, digest, last_login, autoseen, numentry, nopinlist, numsubstr FROM users WHERE id = ?"), c.UserID)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (c *Client) UserByName(ctx context.Context, name string) (*User, error) {
	u := User{}
	err := c.GetContext(ctx, &u, c.sql("SELECT id, name, digest, last_login, autoseen, numentry, nopinlist, numsubstr FROM users WHERE name = ?"), name)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (c *UserClient) UpdateUserDigest(ctx context.Context, digest string) error {
	_, err := c.ExecContext(ctx, c.sql("UPDATE users SET digest = ?  WHERE id = ?"), digest, c.UserID)
	return err
}

func (c *Client) UpdateUserLastLogin(ctx context.Context, userID uint64) error {
	_, err := c.ExecContext(ctx, c.sql("UPDATE users SET last_login = CURRENT_TIMESTAMP WHERE id = ?"), userID)
	return err
}

func (c *Client) InsertUser(ctx context.Context, username, password string) error {
	_, err := c.ExecContext(ctx, c.sql("INSERT INTO users (name, digest) VALUES (?, ?)"), username, password)
	return err
}
