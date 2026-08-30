package cleaner

import (
	"context"
	"log/slog"

	"github.com/yseto/gion-go/internal/config"
	"github.com/yseto/gion-go/internal/db"
	"github.com/yseto/gion-go/internal/websub"
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
			tx.Rollback()
			return err
		}
		if len(subsByFeed) == 0 {
			if err = tx.DeleteFeed(ctx, feeds[i].ID); err != nil {
				tx.Rollback()
				return err
			}
			if id := feeds[i].WebsubID; id != nil {
				item, err := tx.WebsubByID(ctx, *id)
				if err != nil {
					tx.Rollback()
					return err
				}
				if err = tx.DeleteWebsub(ctx, *item); err != nil {
					tx.Rollback()
					return err
				}
				if err = websub.UnSubscribe(ctx, item.TopicURL); err != nil {
					return err
				}
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
