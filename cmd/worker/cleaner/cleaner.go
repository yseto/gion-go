package cleaner

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"

	"github.com/yseto/gion-go/config"
	"github.com/yseto/gion-go/db"
)

const (
	TypeCleaner = "cleaner"
)

type CleanerPayload struct {
}

func NewCleanerTask() (*asynq.Task, error) {
	payload, err := json.Marshal(CleanerPayload{})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeCleaner, payload), nil
}

func HandleCleanerTask(ctx context.Context, t *asynq.Task) error {
	cfg := ctx.Value("cfg").(*config.Config)

	var p CleanerPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	dbConn, err := db.Open(cfg)
	if err != nil {
		return err
	}
	defer dbConn.Close()

	dbc := db.New(dbConn)
	if err = dbc.PurgeReadEntry(); err != nil {
		return err
	}

	feeds, err := dbc.Feeds()
	if err != nil {
		return err
	}

	tx := dbc.MustBegin()
	for i := range feeds {
		subsByFeed, err := tx.SubscriptionByFeedID(feeds[i].ID)
		if err != nil {
			return err
		}
		if len(subsByFeed) == 0 {
			if err = tx.DeleteFeed(feeds[i].ID); err != nil {
				return err
			}
		}
	}
	tx.Commit()

	tx = dbc.MustBegin()
	stories, err := tx.Stories()
	if err != nil {
		return err
	}
	for i := range stories {
		exist, err := tx.ExistEntry(stories[i].FeedID, stories[i].Serial)
		if err != nil {
			return err
		}
		if exist == 0 {
			if err = tx.DeleteStory(stories[i].FeedID, stories[i].Serial); err != nil {
				return err
			}
		}
	}
	tx.Commit()
	fmt.Println("Cleanup done.")
	return nil
}
