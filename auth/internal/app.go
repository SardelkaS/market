package internal

import (
	"auth/config"
	"auth/pkg/db/postgres"
	"auth/pkg/db/redis"
	"context"
	"fmt"
)

func Init(ctx context.Context) error {
	cfg := config.Get()

	connectionUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s pool_max_conns=%s pool_min_conns=%s pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s pool_health_check_period=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
		cfg.Postgres.SSLMode,
		cfg.Postgres.MaxConns,
		cfg.Postgres.MinConns,
		cfg.Postgres.MaxConnLifetime,
		cfg.Postgres.MaxConnIdleTime,
		cfg.Postgres.HealthCheckDuration)
	postgres.InitPsqlDB(connectionUrl, ctx)

	redis.InitRedisClient(redis.Config{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
	}, ctx)

	return nil
}
