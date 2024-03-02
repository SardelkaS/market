package internal

import (
	"core/config"
	auth_usecase "core/internal/auth/usecase"
	"core/internal/basket"
	basket_repository "core/internal/basket/repository"
	basket_usecase "core/internal/basket/usecase"
	"core/internal/feedback"
	feedback_repository "core/internal/feedback/repository"
	feedback_usecase "core/internal/feedback/usecase"
	"core/internal/order"
	order_repository "core/internal/order/repository"
	order_usecase "core/internal/order/usecase"
	"core/internal/product"
	product_repository "core/internal/product/repository"
	product_usecase "core/internal/product/usecase"
	"core/internal/tg_bot"
	tg_bot_usecase "core/internal/tg_bot/usecase"
	"core/pkg/db"
	"core/pkg/logger"
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

	a.Repo["basket"] = basket_repository.NewPostgresRepo(a.dbConnection["postgres"].(db.Connection))
	a.Repo["order"] = order_repository.NewPostgresRepo(a.dbConnection["postgres"].(db.Connection))
	a.Repo["product"] = product_repository.NewPostgresRepo(a.dbConnection["postgres"].(db.Connection))
	a.Repo["feedback"] = feedback_repository.NewPostgresRepo(a.dbConnection["postgres"].(db.Connection))

	a.UC["basket"] = basket_usecase.New(a.Repo["basket"].(basket.Repository), a.Repo["product"].(product.Repository))
	a.UC["tg_bot"], err = tg_bot_usecase.New(a.cfg)
	if err != nil {
		return err
	}
	a.UC["order"] = order_usecase.New(a.Repo["order"].(order.Repository), a.Repo["product"].(product.Repository), a.UC["tg_bot"].(tg_bot.UC))
	a.UC["product"] = product_usecase.New(a.Repo["product"].(product.Repository))
	a.UC["feedback"] = feedback_usecase.New(a.Repo["feedback"].(feedback.Repository), a.UC["product"].(product.UC))
	a.UC["auth"] = auth_usecase.New(a.cfg)

	return nil
}
