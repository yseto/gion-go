package worker

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/yseto/gion-go/internal/cleaner"
	"github.com/yseto/gion-go/internal/config"
)

type wCleaner struct {
	baseContext func() context.Context

	cfg *config.Config

	wg       sync.WaitGroup
	shutdown chan struct{}
}

func NewCleaner(cfg *config.Config, baseContext func() context.Context) *wCleaner {
	return &wCleaner{
		baseContext: baseContext,

		cfg:      cfg,
		shutdown: make(chan struct{}),
	}
}

func (w *wCleaner) Serve() error {
	ticker := time.NewTicker(12 * time.Hour)
	defer ticker.Stop()

	baseCtx := context.Background()
	if w.baseContext != nil {
		baseCtx = w.baseContext()
		if baseCtx == nil {
			panic("BaseContext returned a nil context")
		}
	}

	ctx, cancel := context.WithCancel(baseCtx)
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
			err := cleaner.Task(ctx)
			if err != nil {
				slog.Error("cleaner", "error", err.Error())
			}
			// println("quit tick")
		case <-quit:
			return nil
		}
	}
}

func (w *wCleaner) Shutdown(_ context.Context) error {
	slog.Info("wCleaner shutdown...")
	close(w.shutdown)
	w.wg.Wait()
	slog.Info("wCleaner shutdown")
	return nil
}
