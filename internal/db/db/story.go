package db

func (c *ClientTxn) InsertStory(feedID, serial uint64, title, description, url string) error {
	_, err := c.Exec(c.sql(`
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

func (c *ClientTxn) Stories() ([]*Story, error) {
	s := []*Story{}
	err := c.Select(&s, "SELECT feed_id, serial, url FROM story")
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (c *ClientTxn) DeleteStory(feedID, serial uint64) error {
	_, err := c.Exec(c.sql("DELETE FROM story WHERE feed_id = ? AND serial = ?"), feedID, serial)
	return err
}
