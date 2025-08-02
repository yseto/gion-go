package cleaner

import (
	"context"
	"fmt"

	"github.com/yseto/gion-go/config"
	"github.com/yseto/gion-go/db"
)

func Task(ctx context.Context) error {
	cfg, _ := config.ReadConfig()

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
