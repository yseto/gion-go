package main

import (
	"context"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/lestrrat-go/accesslog"

	"github.com/yseto/gion-go/internal/config"
	"github.com/yseto/gion-go/internal/db"
	"github.com/yseto/gion-go/internal/reader"
	"github.com/yseto/gion-go/internal/worker"
)

type serveAndShutdown interface {
	Serve() error
	Shutdown(ctx context.Context) error
}

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	dbConn, err := db.Open(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	al := accesslog.New()

	var listener net.Listener
	if cfg.SystemdSocket {
		fdNum, err := strconv.Atoi(os.Getenv("LISTEN_FDS"))
		if err != nil || fdNum == 0 {
			log.Fatal("Socket not found")
		}
		const fdsStart = 3
		for fd := fdsStart; fd < fdsStart+fdNum; fd++ {
			syscall.CloseOnExec(fd)
			f := os.NewFile(uintptr(fd), "")

			if pc, err := net.FileListener(f); err == nil {
				listener = pc
				f.Close()
				break
			}
		}
	} else {
		listener, err = net.Listen("tcp", net.JoinHostPort(cfg.HTTPHost, cfg.HTTPPort))
		if err != nil {
			log.Fatal(err)
		}
	}

	reader := &httpServer{
		s: &http.Server{
			Handler: al.Wrap(reader.CreateHandler(cfg, dbConn)),
		},
		l: listener,
	}

	wCrawler := worker.NewCrawler(cfg)

	wCleaner := worker.NewCleaner(cfg)

	idleShutdown := make(chan struct{})

	var srvs = []serveAndShutdown{reader, wCrawler, wCleaner}

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var wg sync.WaitGroup
		for _, s := range srvs {
			wg.Add(1)
			go func(ctx context.Context, s serveAndShutdown) {
				defer wg.Done()
				if err := s.Shutdown(ctx); err != nil {
					log.Fatal(err)
				}
			}(ctx, s)
		}
		wg.Wait()

		close(idleShutdown)
	}()

	for _, s := range srvs {
		go func(s serveAndShutdown) {
			if err := s.Serve(); err != nil {
				log.Fatal("shutting down the server")
			}
		}(s)
	}

	slog.Info("initialized.")
	<-idleShutdown
}

type httpServer struct {
	s *http.Server
	l net.Listener
}

func (s *httpServer) Serve() error {
	if err := s.s.Serve(s.l); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *httpServer) Shutdown(ctx context.Context) error {
	return s.s.Shutdown(ctx)
}
