package cleaner

import (
	"context"
	"log/slog"

	"github.com/yseto/gion-go/internal/config"
	"github.com/yseto/gion-go/internal/db"
)

func Task(ctx context.Context) error {
	logger := slog.Default().With("task", "cleaner")

	cfg, _ := config.ReadConfig()

	dbConn, err := db.Open(cfg)
	if err != nil {
		return err
	}
	defer dbConn.Close()

	dbc := db.New(dbConn)
	if err = dbc.PurgeReadEntry(ctx); err != nil {
		return err
	}

	feeds, err := dbc.Feeds(ctx)
	if err != nil {
		return err
	}

	tx := dbc.MustBegin()
	for i := range feeds {
		subsByFeed, err := tx.SubscriptionByFeedID(ctx, feeds[i].ID)
		if err != nil {
			return err
		}
		if len(subsByFeed) == 0 {
			if err = tx.DeleteFeed(ctx, feeds[i].ID); err != nil {
				return err
			}
		}
	}
	tx.Commit()

	tx = dbc.MustBegin()
	stories, err := tx.Stories(ctx)
	if err != nil {
		return err
	}
	for i := range stories {
		exist, err := tx.ExistEntry(ctx, stories[i].FeedID, stories[i].Serial)
		if err != nil {
			return err
		}
		if exist == 0 {
			if err = tx.DeleteStory(ctx, stories[i].FeedID, stories[i].Serial); err != nil {
				return err
			}
		}
	}
	tx.Commit()
	logger.Info("Cleanup done.")
	return nil
}
