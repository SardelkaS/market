package auth_usecase

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"market_auth/config"
	"market_auth/internal/auth"
	"market_auth/internal/auth/model"
	"market_auth/internal/failure"
	"market_auth/pkg/logger"
	"market_auth/pkg/secure"
	"strconv"
	"time"
)

type uc struct {
	repo   auth.Repository
	redis  auth.CacheRepository
	cfg    *config.Config
	logger logger.UC
}

func NewUC(postgres auth.Repository, redis auth.CacheRepository, cfg *config.Config, logger logger.UC) auth.UC {
	return &uc{
		repo:   postgres,
		redis:  redis,
		cfg:    cfg,
		logger: logger,
	}
}

func (u *uc) GenerateToken(duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["exp"] = time.Now().Add(time.Minute * duration).Unix()

	tokenString, err := token.SignedString([]byte(u.cfg.Auth.Secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (u *uc) SignIn(input auth_model.SignInLogicInput) (*auth_model.SignInToken, error) {
	userData, err := u.repo.GetUserByName(*input.Login)
	if err != nil || userData == nil {
		u.logger.Log(logger.Error, fmt.Sprintf("get auth by login (%s) error: %v", *input.Login, err))
		return nil, failure.ErrNotFound
	}

	err = secure.ComparePassword(*input.Password, *userData.Password)
	if err != nil {
		u.logger.Log(logger.Error, fmt.Sprintf("wrong password(%s) for auth %s", *input.Password, *input.Login))
		return nil, failure.ErrPasswordNotCorrect
	}

	accessToken, err := u.GenerateToken(time.Duration(u.cfg.Auth.AccessLifeTime))
	if err != nil {
		u.logger.Log(logger.Error, fmt.Sprintf("generate access token for auth %s error: %s", *input.Login, err.Error()))
		return nil, failure.ErrJWTGenerate
	}
	refreshToken, err := u.GenerateToken(time.Duration(u.cfg.Auth.RefreshLifeTime))
	if err != nil {
		u.logger.Log(logger.Error, fmt.Sprintf("generate refresh token for auth %s error: %s", *input.Login, err.Error()))
		return nil, failure.ErrJWTGenerate
	}
	finger := secure.GenerateRandomString(32)
	err = u.redis.AddTokensInDB(auth_model.AddTokenGatewayInput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       strconv.FormatInt(*userData.Id, 10),
		FingerKey:    finger,
	})
	if err != nil {
		u.logger.Log(logger.Error, fmt.Sprintf("add tokens for auth %s to redis error: %s", *input.Login, err.Error()))
		return nil, failure.ErrJWTGenerate
	}

	return &auth_model.SignInToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		FingerKey:    finger,
		Role:         []string{*userData.Role},
	}, nil
}

func (u *uc) SignOut(params auth_model.SignOutLogicInput) error {
	err := u.redis.DeleteDataByFingerKey(auth_model.DeleteDataGatewayInput{FingerKey: params.FingerKey})
	if err != nil {
		return failure.ErrInternal
	}

	return nil
}

func (u *uc) SignUp(params auth_model.SignUpLogicInput) (*auth_model.SignUpToken, error) {
	passwordHash, err := secure.HashPassword(params.Password)
	if err != nil {
		u.logger.Log(logger.Error, fmt.Sprintf("Hash password for auth(%s) error: %s", params.Login, err.Error()))
		return nil, failure.ErrHashingPassword
	}

	internalId := secure.CalcSignature(time.Now().String(), params.Login)

	ban := false
	defaultTimezone := "UTC"
	defaultRole := "user"
	err = u.repo.InsertUser(auth_model.User{
		Login:      &params.Login,
		Password:   &passwordHash,
		InternalId: &internalId,
		Role:       &defaultRole,
		Ban:        &ban,
		Timezone:   &defaultTimezone,
	})
	if err != nil {
		u.logger.Log(logger.Error, fmt.Sprintf("Insert auth %s error: %s", params.Login, err.Error()))
		return nil, failure.ErrInternal
	}
	userInfo, err := u.repo.GetUserByName(params.Login)
	if err != nil {
		u.logger.Log(logger.Error, fmt.Sprintf("Get auth %s by name error: %s", params.Login, err.Error()))
		return nil, failure.ErrInternal
	}

	accessToken, err := u.GenerateToken(time.Duration(u.cfg.Auth.AccessLifeTime))
	if err != nil {
		u.logger.Log(logger.Error, fmt.Sprintf("Generate access token for auth %s error: %s", params.Login, err.Error()))
		return nil, failure.ErrJWTGenerate
	}
	refreshToken, err := u.GenerateToken(time.Duration(u.cfg.Auth.RefreshLifeTime))
	if err != nil {
		u.logger.Log(logger.Error, fmt.Sprintf("Generate refresh token for auth %s error: %s", params.Login, err.Error()))
		return nil, failure.ErrJWTGenerate
	}
	finger := secure.GenerateRandomString(32)
	err = u.redis.AddTokensInDB(auth_model.AddTokenGatewayInput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       strconv.FormatInt(*userInfo.Id, 10),
		FingerKey:    finger,
	})
	if err != nil {
		u.logger.Log(logger.Error, fmt.Sprintf("Add tokens for auth %s to redis error: %s", params.Login, err.Error()))
		return nil, failure.ErrJWTGenerate
	}

	return &auth_model.SignUpToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		FingerKey:    finger,
		Role:         []string{*userInfo.Role},
	}, nil
}

func (u *uc) ValidateJWT(params auth_model.ValidateJWTLogicInput) (*auth_model.ValidateJWTLogicOutput, error) {
	if len(params.Token) <= 0 || len(params.FingerKey) <= 0 {
		u.logger.Log(logger.Error, fmt.Sprintf("Wrong params for finger key %s", params.FingerKey))
		return nil, failure.ErrAuth
	}
	token, err := jwt.Parse(params.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(u.cfg.Auth.Secret), nil
	})
	if err != nil || token == nil || !token.Valid {
		u.logger.Log(logger.Error, fmt.Sprintf("Not valid token for finger key %s; error: %v", params.FingerKey, err))
		return nil, failure.ErrAuth
	}

	userData, err := u.GetUserByFingerKey(params.FingerKey)
	if err != nil {
		return nil, err
	}

	return &auth_model.ValidateJWTLogicOutput{
		Role:   userData.Role,
		UserId: userData.Id,
	}, nil
}

func (u *uc) GenerateRefresh(params auth_model.GenerateRefreshLogicInput) (*auth_model.Token, error) {
	isValid, err := u.redis.IsValidRefreshToken(auth_model.IsValidRefreshTokenGatewayInput{
		FingerKey:    params.FingerKey,
		RefreshToken: params.RefreshToken,
	})
	if !isValid || err != nil {
		u.logger.Log(logger.Error, fmt.Sprintf("Not valid token for finger key %s error: %s", params.FingerKey, err.Error()))
		return nil, failure.ErrJWTNotValid
	}

	accessToken, err := u.GenerateToken(time.Duration(u.cfg.Auth.AccessLifeTime))
	if err != nil {
		u.logger.Log(logger.Error, fmt.Sprintf("Generate access token for finger key %s error: %s", params.FingerKey, err.Error()))
		return nil, failure.ErrJWTGenerate
	}
	refreshToken, err := u.GenerateToken(time.Duration(u.cfg.Auth.RefreshLifeTime))
	if err != nil {
		u.logger.Log(logger.Error, fmt.Sprintf("Generate refresh token for finger key %s error: %s", params.FingerKey, err.Error()))
		return nil, failure.ErrJWTGenerate
	}
	err = u.redis.UpdateTokenByFingerKey(auth_model.UpdateTokenByFingerKeyGatewayInput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		FingerKey:    params.FingerKey,
	})
	if err != nil {
		u.logger.Log(logger.Error, fmt.Sprintf("Add tokens for finger key %s to redis error: %s", params.FingerKey, err.Error()))
		return nil, failure.ErrJWTGenerate
	}

	return &auth_model.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		FingerKey:    params.FingerKey,
	}, nil
}

func (u *uc) ChangePassword(input auth_model.ChangePasswordLogicInput) error {
	userData, err := u.GetUserByFingerKey(input.FingerKey)
	if err != nil {
		return err
	}

	err = secure.ComparePassword(input.OldPassword, *userData.Password)
	if err != nil {
		u.logger.Log(logger.Error, fmt.Sprintf("Wrong password for auth %s", *userData.Login))
		return failure.ErrAuth
	}

	passwordHash, err := secure.HashPassword(input.NewPassword)
	if err != nil {
		u.logger.Log(logger.Error, fmt.Sprintf("Hash password for auth(%d) error: %s", *userData.Id, err.Error()))
		return failure.ErrHashingPassword
	}
	err = u.repo.UpdatePassword(auth_model.UpdatePasswordGatewayInput{
		UserId:         *userData.Id,
		HashedPassword: passwordHash,
	})
	return err
}

func (u *uc) ChangeTimezone(input auth_model.ChangeTimezoneLogicInput) error {
	userData, err := u.GetUserByFingerKey(input.FingerKey)
	if err != nil {
		return err
	}

	err = u.repo.UpdateTimezone(auth_model.UpdateTimezoneGatewayInput{
		UserId:      *userData.Id,
		NewTimezone: input.NewTimezone,
	})
	if err != nil {
		return failure.ErrChangeTimezone
	}
	return nil
}

func (u *uc) GetUserByFingerKey(fingerKey string) (*auth_model.User, error) {
	data, err := u.redis.GetDataByFingerKey(auth_model.GetDataByFingerKeyGatewayInput{FingerKey: fingerKey})
	if err != nil {
		u.logger.Log(logger.Error, fmt.Sprintf("Get info for finger key %s from redis error: %s", fingerKey, err.Error()))
		return nil, failure.ErrAuth
	}

	userId, err := strconv.ParseInt(data.UserId, 10, 64)
	if err != nil {
		u.logger.Log(logger.Error, fmt.Sprintf("Wrong id format for finger key %s; error: %s", fingerKey, err.Error()))
		return nil, failure.ErrInternal
	}

	userData, err := u.repo.GetUserById(userId)
	if err != nil {
		u.logger.Log(logger.Error, fmt.Sprintf("User not found(%d) error: %v", userId, err))
		return nil, failure.ErrNotFound
	}
	return userData, nil
}

func (u *uc) GetUser(input auth_model.GetUserLogicInput) (*auth_model.GetUserLogicOutput, error) {
	var userData *auth_model.User
	var err error

	if input.FingerKey != nil {
		userData, err = u.GetUserByFingerKey(*input.FingerKey)
		if err != nil {
			return nil, err
		}
		return &auth_model.GetUserLogicOutput{
			Result: userData,
		}, nil
	}
	if input.Id != nil {
		userData, err = u.repo.GetUserById(*input.Id)
		if err != nil {
			u.logger.Log(logger.Error, fmt.Sprintf("get user info by id(%d) error: %v", *input.Id, err))
			return nil, failure.ErrNotFound
		}
		return &auth_model.GetUserLogicOutput{
			Result: userData,
		}, nil
	}
	userData, err = u.repo.GetUserByName(*input.Login)
	if err != nil {
		u.logger.Log(logger.Error, fmt.Sprintf("get user by name(%s) error: %v", *input.Login, err))
		return nil, failure.ErrNotFound
	}

	return &auth_model.GetUserLogicOutput{
		Result: userData,
	}, nil
}
