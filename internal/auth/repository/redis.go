package auth_repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"market_auth/config"
	"market_auth/internal/auth"
	"market_auth/internal/auth/model"
	"market_auth/pkg/logger"
	"strings"
	"time"
)

type redisClient struct {
	redisClient *redis.Client
	config      *config.Config
	logger      logger.UC
	accessTime  time.Duration
	refreshTime time.Duration
}

func NewRedisClient(cfg *config.Config, logger logger.UC) auth.CacheRepository {
	c := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
	})
	if err := c.Ping().Err(); err != nil {
		panic("Unable to connect to redis " + err.Error())
	}

	return &redisClient{
		redisClient: c,
		config:      cfg,
		logger:      logger,
		accessTime:  time.Duration(cfg.Auth.AccessLifeTime),
		refreshTime: time.Duration(cfg.Auth.RefreshLifeTime),
	}
}

func (c *redisClient) GetKey(key string) (string, error) {
	val, err := c.redisClient.Get(key).Result()
	if err == redis.Nil || err != nil {
		return "", err
	}
	val = strings.Trim(val, "\"")
	return val, nil
}

func (c *redisClient) Keys(pattern string) ([]string, error) {
	val, err := c.redisClient.Keys(pattern).Result()
	if err == redis.Nil || err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	if len(val) == 0 {
		return nil, errors.New("not found redis keys by pattern " + pattern)
	}
	return val, nil
}

func (c *redisClient) Del(pattern string) (int64, error) {
	val, err := c.redisClient.Del(pattern).Result()
	if err == redis.Nil || err != nil {
		return 0, err
	}
	return val, nil
}

func (c *redisClient) SetKey(key string, value interface{}, expiration time.Duration) error {
	cacheEntry, err := json.Marshal(value)
	if err != nil {
		return err
	}
	key = strings.Trim(key, "\"")
	err = c.redisClient.Set(key, cacheEntry, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *redisClient) AddTokensInDB(params auth_model.AddTokenGatewayInput) error {
	err := c.DeleteDataByFingerKey(auth_model.DeleteDataGatewayInput{FingerKey: params.FingerKey})
	if err != nil {
		c.logger.Log(logger.Info, fmt.Sprintf("Redis AddTokensInDB try delete data by fingerkey. %s", err.Error()))
	}

	err = c.SetKey(params.FingerKey+":"+params.UserID+":"+"access", params.AccessToken, time.Minute*c.accessTime)
	if err != nil {
		return err
	}

	err = c.SetKey(params.FingerKey+":"+params.UserID+":"+"refresh", params.RefreshToken, time.Minute*c.refreshTime)
	if err != nil {
		return err
	}

	err = c.SetKey(params.FingerKey+":"+params.UserID+":"+"uid", params.UserID, time.Minute*c.refreshTime)
	if err != nil {
		return err
	}

	return nil
}

func (c *redisClient) DeleteDataByFingerKey(params auth_model.DeleteDataGatewayInput) error {
	resultsUid, err := c.Keys(params.FingerKey + ":*:*")
	if err != nil {
		return err
	}
	if len(resultsUid) == 0 {
		return errors.New("data by fingerKey not found")
	}
	for _, v := range resultsUid {
		_, err := c.Del(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *redisClient) IsValidRefreshToken(params auth_model.IsValidRefreshTokenGatewayInput) (bool, error) {
	if params.RefreshToken == "" {
		return false, errors.New("IsValidRefreshToken(): refreshToken empty: " + params.RefreshToken)
	}

	uid, err := c.GetUserIdByFingerKey(auth_model.GetUserIdByFingerKeyGatewayInput{FingerKey: params.FingerKey})
	if err != nil {
		return false, err
	}

	patternStr := params.FingerKey + ":" + uid + ":refresh"
	refreshT, err := c.GetKey(patternStr)
	if err != nil || refreshT == "" {
		return false, err
	}

	if refreshT != params.RefreshToken {
		return false, errors.New("invalid refresh token " + params.RefreshToken)
	}

	return true, nil
}

func (c *redisClient) UpdateTokenByFingerKey(params auth_model.UpdateTokenByFingerKeyGatewayInput) error {
	if params.AccessToken == "" {
		return errors.New("access token is empty: " + params.AccessToken)
	}
	uid, err := c.GetUserIdByFingerKey(auth_model.GetUserIdByFingerKeyGatewayInput{FingerKey: params.FingerKey})
	if err != nil {
		return err
	}

	err = c.SetKey(params.FingerKey+":"+uid+":"+"access", params.AccessToken, time.Minute*c.accessTime)
	if err != nil {
		return err
	}
	err = c.SetKey(params.FingerKey+":"+uid+":"+"refresh", params.RefreshToken, time.Minute*c.refreshTime)
	if err != nil {
		return err
	}

	return nil
}

func (c *redisClient) GetUserIdByFingerKey(params auth_model.GetUserIdByFingerKeyGatewayInput) (string, error) {
	if params.FingerKey == "" {
		c.logger.Log(logger.Error, "Empty finger key")
		return "", errors.New("fingerKey is empty")
	}
	resultsUid, err := c.Keys(params.FingerKey + ":*:uid")
	if len(resultsUid) == 0 {
		return "", err
	}
	if err != nil || (len(resultsUid) > 0 && resultsUid[0] == "") {
		return "", err
	}
	uid, err := c.GetKey(resultsUid[0])
	if err != nil || uid == "" {
		return "", err
	}

	return uid, nil
}

func (c *redisClient) GetDataByFingerKey(params auth_model.GetDataByFingerKeyGatewayInput) (*auth_model.TokensDataGatewayOutput, error) {
	var data auth_model.TokensDataGatewayOutput
	if params.FingerKey == "" {
		return nil, errors.New("fingerKey is empty")
	}
	uid, err := c.GetUserIdByFingerKey(auth_model.GetUserIdByFingerKeyGatewayInput{FingerKey: params.FingerKey})
	if err != nil {
		return nil, err
	}
	patternStr := params.FingerKey + ":" + fmt.Sprint(uid) + ":access"
	accessToken, err := c.GetKey(patternStr)
	if err != nil || accessToken == "" {
		return nil, err
	}
	patternStr = params.FingerKey + ":" + fmt.Sprint(uid) + ":refresh"
	refreshToken, err := c.GetKey(patternStr)
	if err != nil || refreshToken == "" {
		return nil, err
	}

	data.UserId = uid
	data.AccessToken = accessToken
	data.RefreshToken = refreshToken

	return &data, nil
}
