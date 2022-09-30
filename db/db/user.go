package db

type UserProfile struct {
	AutoSeen        bool   `db:"autoseen" json:"autoseen"`
	EntryCount      uint64 `db:"numentry" json:"numentry"`
	NoPinList       bool   `db:"nopinlist" json:"nopinlist"`
	SubstringLength uint64 `db:"numsubstr" json:"numsubstr"`
}

func (c *UserClient) Profile() (*UserProfile, error) {
	p := UserProfile{}
	err := c.Get(&p, "SELECT autoseen, numentry, nopinlist, numsubstr FROM users WHERE id = ?", c.UserID)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (c *UserClient) UpdateProfile(item UserProfile) error {
	_, err := c.Exec("UPDATE users SET autoseen = ?, numentry = ?, nopinlist = ?, numsubstr = ? WHERE id = ?",
		item.AutoSeen,
		item.EntryCount,
		item.NoPinList,
		item.SubstringLength,
		c.UserID)
	return err
}

type User struct {
	ID        uint64 `db:"id"`
	Name      string `db:"name"`
	Digest    string `db:"digest"`
	LastLogin MyTime `db:"last_login"`
	UserProfile
}

func (c *UserClient) User() (*User, error) {
	u := User{}
	err := c.Get(&u, "SELECT id, `name`, `digest`, last_login, autoseen, numentry, nopinlist, numsubstr FROM users WHERE id = ?", c.UserID)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (c *Client) UserByName(name string) (*User, error) {
	u := User{}
	err := c.Get(&u, "SELECT id, `name`, `digest`, last_login, autoseen, numentry, nopinlist, numsubstr FROM users WHERE name = ?", name)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (c *UserClient) UpdateUserDigest(digest string) error {
	_, err := c.Exec("UPDATE users SET digest = ?  WHERE id = ?", digest, c.UserID)
	return err
}

func (c *Client) UpdateUserLastLogin(userID uint64) error {
	_, err := c.Exec("UPDATE users SET last_login = CURRENT_TIMESTAMP WHERE id = ?", userID)
	return err
}

func (c *Client) InsertUser(username, password string) error {
	_, err := c.Exec("INSERT INTO users (name, digest) VALUES (?, ?)", username, password)
	return err
}
