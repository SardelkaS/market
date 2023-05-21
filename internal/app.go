package internal

import (
	"github.com/jmoiron/sqlx"
	"market_auth/config"
	"market_auth/internal/api"
	"market_auth/internal/auth"
	auth_repository "market_auth/internal/auth/repository"
	auth_usecase "market_auth/internal/auth/usecase"
	proxy_usecase "market_auth/internal/proxy/usecase"
	"market_auth/pkg/db"
	"market_auth/pkg/logger"
)

type App struct {
	cfg          *config.Config
	UC           map[string]interface{}
	Repo         map[string]interface{}
	dbConnection map[string]interface{}
}

func NewApp(cfg *config.Config) *App {
	return &App{
		UC:           make(map[string]interface{}),
		Repo:         make(map[string]interface{}),
		dbConnection: make(map[string]interface{}),
		cfg:          cfg,
	}
}

func (a *App) Init() error {
	var err error
	a.dbConnection["postgres"], err = db.InitPsqlDB(a.cfg)
	if err != nil {
		return err
	}

	a.UC["logger"] = logger.New()

	a.Repo["authPostgres"] = auth_repository.NewPostgresRepo(a.dbConnection["postgres"].(*sqlx.DB))
	a.Repo["authRedis"] = auth_repository.NewRedisClient(a.cfg, a.UC["logger"].(logger.UC))

	a.UC["auth"] = auth_usecase.NewUC(a.Repo["authPostgres"].(auth.Repository), a.Repo["authRedis"].(auth.CacheRepository), a.cfg, a.UC["logger"].(logger.UC))
	a.UC["api"] = api.New(a.cfg)
	a.UC["proxy"] = proxy_usecase.New(a.UC["auth"].(auth.UC), a.UC["api"].(api.UC), a.UC["logger"].(logger.UC))

	return nil
}
