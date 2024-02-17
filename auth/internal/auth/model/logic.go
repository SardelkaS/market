package auth_model

type ValidateJWTLogicInput struct {
	Token     string
	FingerKey string
}

type ValidateJWTLogicOutput struct {
	Role   *string
	UserId *int64
}

type GenerateRefreshLogicInput struct {
	FingerKey    string
	RefreshToken string
}

type SignInLogicInput struct {
	Login    *string
	Password *string
}

type SignOutLogicInput struct {
	FingerKey string
}

type SignUpLogicInput struct {
	Login      string
	Password   string
	InviteCode string
}

type ChangePasswordLogicInput struct {
	FingerKey   string
	OldPassword string
	NewPassword string
}

type ChangeTimezoneLogicInput struct {
	FingerKey   string
	NewTimezone string
}

type GetUserLogicInput struct {
	FingerKey *string
	Id        *int64
	Login     *string
}

type GetUserLogicOutput struct {
	Result *User
}
