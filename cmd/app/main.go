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
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/lestrrat-go/accesslog"

	"github.com/yseto/gion-go/internal/config"
	"github.com/yseto/gion-go/internal/db"
	"github.com/yseto/gion-go/internal/reader"
	"github.com/yseto/gion-go/internal/websub"
	"github.com/yseto/gion-go/internal/worker"
)

type serveAndShutdown interface {
	Serve() error
	Shutdown(ctx context.Context) error
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer stop()

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

	var listener = make(map[string]net.Listener)
	if cfg.SystemdSocket {
		if _, ok := os.LookupEnv("LISTEN_FDNAMES"); !ok {
			log.Fatal("LISTEN_FDNAMES not found")
		}
		fdnames := strings.Split(os.Getenv("LISTEN_FDNAMES"), ":")

		fdNum, err := strconv.Atoi(os.Getenv("LISTEN_FDS"))
		if err != nil || fdNum == 0 {
			log.Fatal("Socket not found")
		}
		const fdsStart = 3
		for fd := fdsStart; fd < fdsStart+fdNum; fd++ {
			syscall.CloseOnExec(fd)
			f := os.NewFile(uintptr(fd), "")

			if pc, err := net.FileListener(f); err == nil {
				listener[fdnames[fd-fdsStart]] = pc
				f.Close()
			}
		}
	} else {
		listener["web"], err = net.Listen("tcp", net.JoinHostPort(cfg.HTTPHost, cfg.HTTPPort))
		if err != nil {
			log.Fatal(err)
		}
		listener["websub"], err = net.Listen("tcp", net.JoinHostPort(cfg.HTTPHost, cfg.WebsubPort))
		if err != nil {
			log.Fatal(err)
		}
	}

	reader := &httpServer{
		s: &http.Server{
			Handler:     al.Wrap(reader.CreateHandler(cfg, dbConn)),
			BaseContext: func(net.Listener) context.Context { return ctx },
		},
		l: listener["web"],
	}

	websub := &httpServer{
		s: &http.Server{
			Handler:     al.Wrap(websub.CreateHandler(cfg, dbConn)),
			BaseContext: func(net.Listener) context.Context { return ctx },
		},
		l: listener["websub"],
	}

	wCrawler := worker.NewCrawler(cfg, func() context.Context { return ctx })

	wCleaner := worker.NewCleaner(cfg, func() context.Context { return ctx })

	wWebSub := worker.NewWebSub(cfg, func() context.Context { return ctx })

	var srvs = []serveAndShutdown{reader, websub, wCrawler, wCleaner, wWebSub}

	for _, s := range srvs {
		go func(s serveAndShutdown) {
			if err := s.Serve(); err != nil {
				log.Fatal("shutting down the server")
			}
		}(s)
	}

	slog.Info("initialized.")
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	for _, s := range srvs {
		wg.Add(1)
		go func(ctx context.Context, s serveAndShutdown) {
			defer wg.Done()
			if err := s.Shutdown(ctx); err != nil {
				log.Fatal(err)
			}
		}(shutdownCtx, s)
	}
	wg.Wait()
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
