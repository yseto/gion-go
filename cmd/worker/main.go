package main

import (
	"context"
	"log"

	"github.com/hibiken/asynq"

	"github.com/yseto/gion-go/cmd/worker/cleaner"
	"github.com/yseto/gion-go/cmd/worker/crawler"
	"github.com/yseto/gion-go/config"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.RedisAddr},
		asynq.Config{
			BaseContext: func() context.Context {
				ctx := context.Background()
				return context.WithValue(ctx, "cfg", cfg)
			},
			Concurrency: 1,
			Queues: map[string]int{
				"default": 1,
			},
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(crawler.TypeCrawler, crawler.HandleCrawlerTask)
	mux.HandleFunc(cleaner.TypeCleaner, cleaner.HandleCleanerTask)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
