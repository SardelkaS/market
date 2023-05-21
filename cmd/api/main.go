package main

import (
	"fmt"
	"market_auth/config"
	"market_auth/internal"
	"market_auth/internal/server"
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
	fmt.Println(app)
	httpServer := server.NewHttpServer(cfg)
	err = app.Init()
	if err != nil {
		panic(err)
	}
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
