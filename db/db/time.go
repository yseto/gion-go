package db

import "time"

type MyTime time.Time

func (h *MyTime) MarshalJSON() ([]byte, error) {
	ht := time.Time(*h)
	return []byte(`"` + ht.Format("2006-01-02 15:04:05") + `"`), nil
}

func (h *MyTime) Sub(u time.Time) time.Duration {
	return time.Time(*h).Sub(u)
}
