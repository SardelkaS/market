package redis

import (
	"context"
	"github.com/go-redis/redis"
	"sync"
)

var (
	client Client
	once   = &sync.Once{}
	mux    = &sync.RWMutex{}
)

const Nil = redis.Nil

func Get() Client {
	mux.RLock()
	defer mux.RUnlock()

	return client
}

func Set(c Client) {
	mux.Lock()
	defer mux.Unlock()

	client = c
}

type Config struct {
	Host     string
	Port     string
	Password string
}

type Client struct {
	*redis.Client
}

func InitRedisClient(cfg Config, ctx context.Context) {
	once.Do(func() {
		c := redis.NewClient(&redis.Options{
			Addr:     cfg.Host + ":" + cfg.Port,
			Password: cfg.Password,
		})
		if err := c.Ping().Err(); err != nil {
			panic("Unable to connect to redis " + err.Error())
		}

		Set(Client{c})

		go func() {
			<-ctx.Done()
			_ = c.Close()
		}()
	})
}
