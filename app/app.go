package app

import (
	"context"
	"github.com/smalls0098/gopkg/app/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func WithServer(servers ...server.Server) Option {
	return func(a *App) { a.servers = servers }
}

func WithName(name string) Option {
	return func(a *App) { a.name = name }
}

type Option func(a *App)

type App struct {
	name    string
	servers []server.Server
}

func New(opts ...Option) *App {
	a := &App{
		name:    "test",
		servers: make([]server.Server, 0, 1),
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func (a *App) Run(ctx context.Context) error {
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(signals)

	for _, srv := range a.servers {
		go func(srv server.Server) {
			err := srv.Start(ctx)
			if err != nil {
				log.Printf("Server start err: %v", err)
			}
		}(srv)
	}

	select {
	case <-signals:
		log.Println("Received termination signal") // 终止信号
	case <-ctx.Done():
		log.Println("Context canceled", ctx.Err()) // 取消
	}

	for _, srv := range a.servers {
		err := srv.Stop(ctx)
		if err != nil {
			log.Printf("Server stop err: %v", err)
		}
	}
	return nil
}
