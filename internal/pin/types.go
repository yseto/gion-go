package pin

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type ReadFlag string

const (
	Unseen ReadFlag = "Unseen"
	Seen   ReadFlag = "Seen"
	Setpin ReadFlag = "Setpin"
)

func (c ReadFlag) String() string {
	switch c {
	case Unseen:
		return "Unseen"
	case Seen:
		return "Seen"
	case Setpin:
		return "Setpin"
	default:
		panic("Unknown")
	}
}

// for client response.
func (c ReadFlag) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

// for write to db
func (c ReadFlag) Value() (driver.Value, error) {
	switch c {
	case Unseen:
		return int64(0), nil
	case Seen:
		return int64(1), nil
	case Setpin:
		return int64(2), nil
	}
	return int64(0), errors.New("incompatible value for ReadFlag")
}

// for read from db
func (j *ReadFlag) Scan(src interface{}) error {
	v, ok := src.(int64)
	if !ok {
		return errors.New("incompatible type for ReadFlag")
	}

	switch v {
	case 0:
		*j = "Unseen"
	case 1:
		*j = "Seen"
	case 2:
		*j = "Setpin"
	default:
		return errors.New("incompatible type for ReadFlag")
	}
	return nil
}

type UpdateTime time.Time

func (h *UpdateTime) MarshalJSON() ([]byte, error) {
	ht := time.Time(*h)
	return []byte(`"` + ht.Format("2006-01-02 15:04:05") + `"`), nil
}
