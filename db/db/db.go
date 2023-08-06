package db

import (
	"github.com/jmoiron/sqlx"
)

func placeHolder(driverName string) int {
	if driverName == "postgres" {
		return sqlx.DOLLAR
	}
	return sqlx.QUESTION
}

type UserClient struct {
	*sqlx.DB
	UserID uint64
}

func (c *UserClient) sql(query string) string {
	return sqlx.Rebind(placeHolder(c.DriverName()), query)
}

type UserClientTxn struct {
	*sqlx.Tx
	UserID uint64
}

func (c *UserClient) MustBegin() *UserClientTxn {
	return &UserClientTxn{c.DB.MustBegin(), c.UserID}
}

func (c *UserClientTxn) sql(query string) string {
	return sqlx.Rebind(placeHolder(c.DriverName()), query)
}

type Client struct {
	*sqlx.DB
}

func (c *Client) sql(query string) string {
	return sqlx.Rebind(placeHolder(c.DriverName()), query)
}

type ClientTxn struct {
	*sqlx.Tx
}

func (c *Client) MustBegin() *ClientTxn {
	return &ClientTxn{c.DB.MustBegin()}
}

func (c *ClientTxn) sql(query string) string {
	return sqlx.Rebind(placeHolder(c.DriverName()), query)
}
