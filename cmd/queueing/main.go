package main

import (
	"flag"
	"log"

	"github.com/hibiken/asynq"

	"github.com/yseto/gion-go/cmd/worker/cleaner"
	"github.com/yseto/gion-go/cmd/worker/crawler"
	"github.com/yseto/gion-go/config"
	"github.com/yseto/gion-go/db"
	dbType "github.com/yseto/gion-go/db/db"
)

var redisAddr string

func main() {
	mode := flag.String("mode", "", "cleanup or crawler")
	term := flag.Uint64("term", 0, "0=all, =1,2,3...")
	flag.Parse()

	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	client := asynq.NewClient(asynq.RedisClientOpt{Addr: cfg.RedisAddr})
	defer client.Close()

	switch {
	case *mode == "crawler":
		doCrawler(client, cfg, term)
	case *mode == "cleanup":
		doCleaner(client)
	default:
		log.Fatal("mode is empty")
	}

}

func doCleaner(client *asynq.Client) {
	task, err := cleaner.NewCleanerTask()
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}
	info, err := client.Enqueue(task)
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
	return
}

func doCrawler(client *asynq.Client, cfg *config.Config, term *uint64) {
	dbConn, err := db.Open(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	dbc := db.New(dbConn)

	var feeds []*dbType.Feed
	if *term == 0 {
		feeds, err = dbc.Feeds()
	} else {
		feeds, err = dbc.FeedsByTerm(*term)
	}
	if err != nil {
		log.Fatal(err)
	}

	for i := range feeds {
		task, err := crawler.NewCrawlerTask(feeds[i].ID)
		if err != nil {
			log.Fatalf("could not create task: %v", err)
		}
		info, err := client.Enqueue(task)
		if err != nil {
			log.Fatalf("could not enqueue task: %v", err)
		}
		log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
	}
}
