package main

import (
	"context"
	"errors"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/config"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/httpapi"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/repository"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/service"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/storage"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/telemetry"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/worker"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	if e := run(ctx); e != nil {
		slog.Error("server stopped", "error", e)
		os.Exit(1)
	}
}
func run(ctx context.Context) error {
	cfg := config.Load()
	if e := cfg.Validate(); e != nil {
		return e
	}
	db, e := storage.Open(ctx, cfg.DatabaseURL)
	if e != nil {
		return e
	}
	defer db.Close()
	repos := repository.New(db)
	svc := service.New(repos)
	api := httpapi.New(svc)
	srv := &http.Server{Addr: ":" + cfg.Port, Handler: api.Routes(), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second}
	client := &telemetry.MemoryClient{}
	outbox := &worker.Outbox{Repos: repos, Client: client, Log: slog.Default()}
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
		_ = client.Close()
	}()
	go func() { <-outbox.Run(ctx, 20) }()
	err := srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
