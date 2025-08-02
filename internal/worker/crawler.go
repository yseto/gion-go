package worker

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/yseto/gion-go/internal/config"
	"github.com/yseto/gion-go/internal/crawler"
	"github.com/yseto/gion-go/internal/db"
)

type wCrawler struct {
	cfg *config.Config

	wg       sync.WaitGroup
	shutdown chan struct{}
}

func NewCrawler(cfg *config.Config) *wCrawler {
	return &wCrawler{
		cfg:      cfg,
		shutdown: make(chan struct{}),
	}
}

func (w *wCrawler) Serve() error {
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan struct{})

	go func() {
		<-w.shutdown
		defer close(quit)

		cancel()
		ticker.Stop()

		// cleanup task...
		slog.Debug("ticker stopped")
	}()

	for {
		select {
		case <-ticker.C:
			w.wg.Add(1)
			defer w.wg.Done()

			// println("Tick at:", t)
			if err := doCrawl(ctx, w.cfg); err != nil {
				return err
			}
			// println("quit tick")
		case <-quit:
			return nil
		}
	}
}

func (w *wCrawler) Shutdown(_ context.Context) error {
	slog.Info("wCrawler shutdown...")
	close(w.shutdown)
	w.wg.Wait()
	slog.Info("wCrawler shutdown")
	return nil
}

func doCrawl(ctx context.Context, cfg *config.Config) error {
	dbConn, err := db.Open(cfg)
	if err != nil {
		return err
	}
	defer dbConn.Close()

	feeds, err := db.New(dbConn).Feeds(ctx)
	if err != nil {
		return err
	}

	var crawlFeedIDs []uint64
	now := time.Now()
	for i := range feeds {
		if ts := feeds[i].GetRevisitAt(); ts != nil && now.Before(*ts) {
			continue
		}
		crawlFeedIDs = append(crawlFeedIDs, feeds[i].ID)
	}
	dbConn.Close() // release db conn.

	// fmt.Println("crawlFeedIDs:", crawlFeedIDs)
	for _, id := range crawlFeedIDs {
		// TODO parallel execution
		err := crawler.Task(ctx, id)
		if err != nil && errors.Is(err, context.Canceled) {
			break
		}
		if err != nil {
			slog.Error("crawler", "error", err.Error())
		}
	}

	return nil
}
