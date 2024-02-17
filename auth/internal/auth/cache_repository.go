package auth

import "auth/internal/auth/model"

type CacheRepository interface {
	AddTokensInDB(params auth_model.AddTokenGatewayInput) error
	DeleteDataByFingerKey(params auth_model.DeleteDataGatewayInput) error
	IsValidRefreshToken(params auth_model.IsValidRefreshTokenGatewayInput) (bool, error)
	UpdateTokenByFingerKey(params auth_model.UpdateTokenByFingerKeyGatewayInput) error
	GetUserIdByFingerKey(params auth_model.GetUserIdByFingerKeyGatewayInput) (string, error)
	GetDataByFingerKey(params auth_model.GetDataByFingerKeyGatewayInput) (*auth_model.TokensDataGatewayOutput, error)
}
