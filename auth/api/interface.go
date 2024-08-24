package api

type Server interface {
	Init() error
	MapHandlers() error
	Run() error
}
