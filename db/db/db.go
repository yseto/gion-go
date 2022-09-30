package db

import (
	"github.com/jmoiron/sqlx"
)

type UserClient struct {
	*sqlx.DB
	UserID uint64
}

type UserClientTxn struct {
	*sqlx.Tx
	UserID uint64
}

func (c *UserClient) MustBegin() *UserClientTxn {
	return &UserClientTxn{c.DB.MustBegin(), c.UserID}
}

type Client struct {
	*sqlx.DB
}

type ClientTxn struct {
	*sqlx.Tx
}

func (c *Client) MustBegin() *ClientTxn {
	return &ClientTxn{c.DB.MustBegin()}
}
