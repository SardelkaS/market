package internal

import (
	"auth/config"
	"auth/internal/auth"
	auth_repository "auth/internal/auth/repository"
	auth_usecase "auth/internal/auth/usecase"
	"auth/pkg/db"
	"auth/pkg/logger"
	"fmt"
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

	connectionUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s pool_max_conns=%s pool_min_conns=%s pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s pool_health_check_period=%s",
		a.cfg.Postgres.Host,
		a.cfg.Postgres.Port,
		a.cfg.Postgres.User,
		a.cfg.Postgres.Password,
		a.cfg.Postgres.DBName,
		a.cfg.Postgres.SSLMode,
		a.cfg.Postgres.MaxConns,
		a.cfg.Postgres.MinConns,
		a.cfg.Postgres.MaxConnLifetime,
		a.cfg.Postgres.MaxConnIdleTime,
		a.cfg.Postgres.HealthCheckDuration)

	a.dbConnection["postgres"], err = db.InitPsqlDB(connectionUrl)
	if err != nil {
		return err
	}

	a.UC["logger"] = logger.New()

	a.Repo["auth_postgres"] = auth_repository.NewPostgresRepo(a.dbConnection["postgres"].(db.Connection))
	a.Repo["auth_redis"] = auth_repository.NewRedisClient(a.cfg, a.UC["logger"].(logger.UC))

	a.UC["auth"] = auth_usecase.NewUC(a.Repo["auth_postgres"].(auth.Repository), a.Repo["auth_redis"].(auth.CacheRepository), a.cfg, a.UC["logger"].(logger.UC))

	return nil
}
