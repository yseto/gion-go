package websub

import (
	"context"
	"log/slog"
	"time"

	"github.com/yseto/gion-go/internal/config"
	"github.com/yseto/gion-go/internal/db"
)

// TODO 削除中の取り合いに注意
func Task(ctx context.Context) error {
	logger := slog.Default().With("task", "websub")

	cfg, _ := config.ReadConfig()

	dbConn, err := db.Open(cfg)
	if err != nil {
		return err
	}
	defer dbConn.Close()

	dbc := db.New(dbConn)

	tx := dbc.MustBegin()
	items, err := tx.Websubs(ctx)
	if err != nil {
		return err
	}
	tx.Commit()

	logger.Info("websub re-subscribe start")

	nextRun := time.Now().Add(24 * time.Hour)

	for i := range items {
		if items[i].ExpiredAt.Before(nextRun) {
			if err := Subscribe(ctx, items[i].TopicURL); err != nil {
				slog.Error("re-subscribe: Subscribe", slog.String("error", err.Error()))
				continue
			}
			slog.Info("re-subscribe", slog.String("topicURL", items[i].TopicURL))
		}
	}

	logger.Info("websub re-subscribe done.")
	return nil
}
