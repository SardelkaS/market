package auth

import (
	"auth/internal/auth/model"
	"github.com/valyala/fasthttp"
)

type UC interface {
	SignUp(params auth_model.SignUpLogicInput) (*auth_model.SignUpToken, error)
	SignIn(params auth_model.SignInLogicInput) (*auth_model.SignInToken, error)
	SignOut(params auth_model.SignOutLogicInput) error

	GenerateRefresh(params auth_model.GenerateRefreshLogicInput) (*auth_model.Token, error)

	GetUserByFingerKey(fingerKey string) (*auth_model.User, error)

	GetUser(input auth_model.GetUserLogicInput) (*auth_model.GetUserLogicOutput, error)
	GetUserInfoById(id int64) (*auth_model.UserInfo, error)

	ValidateJWT(params auth_model.ValidateJWTLogicInput) (*auth_model.ValidateJWTLogicOutput, error)
	ChangePassword(input auth_model.ChangePasswordLogicInput) error
	ChangeTimezone(input auth_model.ChangeTimezoneLogicInput) error
	UpdateUserInfo(input auth_model.UpdateUserInfoBody) error

	Proxy(req *fasthttp.Request, res *fasthttp.Response) error
}
