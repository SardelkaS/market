package internal

import (
	"github.com/jmoiron/sqlx"
	"market_auth/config"
	"market_auth/internal/auth"
	auth_repository "market_auth/internal/auth/repository"
	auth_usecase "market_auth/internal/auth/usecase"
	"market_auth/internal/basket"
	basket_repository "market_auth/internal/basket/repository"
	basket_usecase "market_auth/internal/basket/usecase"
	"market_auth/internal/feedback"
	feedback_repository "market_auth/internal/feedback/repository"
	feedback_usecase "market_auth/internal/feedback/usecase"
	"market_auth/internal/order"
	order_repository "market_auth/internal/order/repository"
	order_usecase "market_auth/internal/order/usecase"
	"market_auth/internal/product"
	product_repository "market_auth/internal/product/repository"
	product_usecase "market_auth/internal/product/usecase"
	"market_auth/internal/tg_bot"
	tg_bot_usecase "market_auth/internal/tg_bot/usecase"
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
	a.Repo["basket"] = basket_repository.NewPostgresRepo(a.dbConnection["postgres"].(*sqlx.DB))
	a.Repo["order"] = order_repository.NewPostgresRepo(a.dbConnection["postgres"].(*sqlx.DB))
	a.Repo["product"] = product_repository.NewPostgresRepo(a.dbConnection["postgres"].(*sqlx.DB))
	a.Repo["feedback"] = feedback_repository.NewPostgresRepo(a.dbConnection["postgres"].(*sqlx.DB))

	a.UC["auth"] = auth_usecase.NewUC(a.Repo["authPostgres"].(auth.Repository), a.Repo["authRedis"].(auth.CacheRepository), a.cfg, a.UC["logger"].(logger.UC))
	a.UC["basket"] = basket_usecase.New(a.Repo["basket"].(basket.Repository), a.Repo["product"].(product.Repository))
	a.UC["tg_bot"], err = tg_bot_usecase.New(a.cfg)
	if err != nil {
		return err
	}
	a.UC["order"] = order_usecase.New(a.Repo["order"].(order.Repository), a.Repo["product"].(product.Repository), a.UC["tg_bot"].(tg_bot.UC))
	a.UC["product"] = product_usecase.New(a.Repo["product"].(product.Repository))
	a.UC["feedback"] = feedback_usecase.New(a.Repo["feedback"].(feedback.Repository), a.UC["product"].(product.UC))

	return nil
}
