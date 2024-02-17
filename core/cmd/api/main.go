package main

import (
	"market_auth/api/http"
	"market_auth/config"
	"market_auth/internal"
	"market_auth/migrations"
)

func main() {
	cfg, err := config.LoadConfig("config/config.yml")
	if err != nil {
		panic(err)
	}
	err = migrations.Up(cfg)
	if err != nil {
		panic(err)
	}
	app := internal.NewApp(cfg)
	err = app.Init()
	if err != nil {
		panic(err)
	}

	httpServer := http.NewHttpServer(cfg)
	err = httpServer.Init()
	if err != nil {
		panic(err)
	}
	err = httpServer.MapHandlers(app)
	if err != nil {
		panic(err)
	}
	err = httpServer.Run()
	if err != nil {
		panic(err)
	}
}
