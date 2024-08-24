package main

import (
	"auth/api/http"
	"auth/config"
	"auth/internal"
	"auth/migrations"
	"context"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfig("config/config.yml")
	if err != nil {
		panic(err)
	}
	err = migrations.Up(cfg)
	if err != nil {
		panic(err)
	}
	err = internal.Init(ctx)
	if err != nil {
		panic(err)
	}

	httpServer := http.NewHttpServer(cfg)
	err = httpServer.Init()
	if err != nil {
		panic(err)
	}
	err = httpServer.MapHandlers()
	if err != nil {
		panic(err)
	}
	err = httpServer.Run()
	if err != nil {
		panic(err)
	}
}
