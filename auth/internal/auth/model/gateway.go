package auth_model

type AddTokenGatewayInput struct {
	AccessToken  string
	RefreshToken string
	UserID       string
	FingerKey    string
}

type DeleteDataGatewayInput struct {
	FingerKey string
}

type IsValidRefreshTokenGatewayInput struct {
	FingerKey    string
	RefreshToken string
}

type UpdatePasswordGatewayInput struct {
	UserId         int64
	HashedPassword string
}

type UpdateTimezoneGatewayInput struct {
	UserId      int64
	NewTimezone string
}

type UpdateTokenByFingerKeyGatewayInput struct {
	AccessToken  string
	RefreshToken string
	FingerKey    string
}

type GetUserIdByFingerKeyGatewayInput struct {
	FingerKey string
}

type GetDataByFingerKeyGatewayInput struct {
	FingerKey string
}

type TokensDataGatewayOutput struct {
	AccessToken  string
	RefreshToken string
	UserId       string
}

type UpdateUserInfoGatewayInput struct {
	Id          *int64
	Name        *string
	Email       *string
	PhoneNumber *string
	BirthDate   *string
	ContactData *string
}
