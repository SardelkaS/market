package auth

import "market_auth/internal/auth/model"

type Repository interface {
	GetUserByName(name string) (*auth_model.User, error)
	GetUserById(id int64) (*auth_model.User, error)
	GetUserByInternalId(internalId string) (*auth_model.User, error)
	InsertUser(user auth_model.User) error

	UpdateUserInfo(input auth_model.UpdateUserInfoGatewayInput) error
	GetUserInfoById(id int64) (*auth_model.UserInfo, error)

	UpdateTimezone(input auth_model.UpdateTimezoneGatewayInput) error
	UpdatePassword(input auth_model.UpdatePasswordGatewayInput) error
}
