package auth

import "market_auth/internal/auth/model"

type UC interface {
	SignUp(params auth_model.SignUpLogicInput) (*auth_model.SignUpToken, error)
	SignIn(params auth_model.SignInLogicInput) (*auth_model.SignInToken, error)
	SignOut(params auth_model.SignOutLogicInput) error

	GenerateRefresh(params auth_model.GenerateRefreshLogicInput) (*auth_model.Token, error)

	GetUser(input auth_model.GetUserLogicInput) (*auth_model.GetUserLogicOutput, error)

	ValidateJWT(params auth_model.ValidateJWTLogicInput) (*auth_model.ValidateJWTLogicOutput, error)
	ChangePassword(input auth_model.ChangePasswordLogicInput) error
	ChangeTimezone(input auth_model.ChangeTimezoneLogicInput) error
}
