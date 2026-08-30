package worker

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/yseto/gion-go/internal/config"
	"github.com/yseto/gion-go/internal/websub"
)

type wWebSub struct {
	cfg *config.Config

	wg       sync.WaitGroup
	shutdown chan struct{}
}

func NewWebSub(cfg *config.Config) *wWebSub {
	return &wWebSub{
		cfg:      cfg,
		shutdown: make(chan struct{}),
	}
}

func (w *wWebSub) Serve() error {
	ticker := time.NewTicker(24 * time.Hour)
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
		fmt.Println("ticker stopped")
	}()

	for {
		select {
		case <-ticker.C:
			w.wg.Add(1)
			defer w.wg.Done()

			// println("Tick at:", t)
			err := websub.Task(ctx)
			if err != nil {
				slog.Error("websubWorker", "error", err.Error())
			}
			// println("quit tick")
		case <-quit:
			return nil
		}
	}
}

func (w *wWebSub) Shutdown(_ context.Context) error {
	slog.Info("wWebSub shutdown...")
	close(w.shutdown)
	w.wg.Wait()
	slog.Info("wWebSub shutdown")
	return nil
}
