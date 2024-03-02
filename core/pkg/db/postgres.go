package db

import (
	"context"
	"errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	_ "github.com/jackc/pgx/stdlib" // pgx driver
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"log"
)

const (
	_connectionAttempts = 10
)

type connection struct {
	db *pgxpool.Pool
}

type Connection interface {
	Select(dest interface{}, query string, args ...interface{}) error
	Get(dest interface{}, query string, args ...interface{}) error
	Queryx(query string, args ...interface{}) (pgx.Rows, error)
	QueryRowx(query string, args ...interface{}) pgx.Row
	Exec(query string, args ...interface{}) (pgconn.CommandTag, error)
}

func InitPsqlDB(connectionUrl string) (Connection, error) {
	cfg, err := pgxpool.ParseConfig(connectionUrl)
	if err != nil {
		return nil, err
	}

	var pool *pgxpool.Pool
	for i := 0; i < _connectionAttempts; i++ {
		pool, err = pgxpool.NewWithConfig(context.Background(), cfg)
		if err != nil {
			log.Printf("ATTEMPT %d ERROR: %s", i+1, err.Error())
			pool = nil
		} else {
			break
		}
	}

	if pool == nil {
		return nil, errors.New("cannot connect to postgres")
	}

	return &connection{
		db: pool,
	}, nil
}

func (c *connection) Select(dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Select(context.Background(), c.db, dest, query, args[:]...)
}

func (c *connection) Get(dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Get(context.Background(), c.db, dest, query, args[:]...)
}

func (c *connection) Queryx(query string, args ...interface{}) (pgx.Rows, error) {
	return c.db.Query(context.Background(), query, args[:]...)
}

func (c *connection) QueryRowx(query string, args ...interface{}) pgx.Row {
	return c.db.QueryRow(context.Background(), query, args[:]...)
}

func (c *connection) Exec(query string, args ...interface{}) (pgconn.CommandTag, error) {
	return c.db.Exec(context.Background(), query, args[:]...)
}
