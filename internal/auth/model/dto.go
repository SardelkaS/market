package auth_model

type SignInToken struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	FingerKey    string   `json:"finger_key"`
	Role         []string `json:"role"`
}

type SignUpToken struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	FingerKey    string   `json:"finger_key"`
	Role         []string `json:"role"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	FingerKey    string `json:"finger_key"`
}

type SignInBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RefreshBody struct {
	FingerKey    string `json:"finger_key"`
	RefreshToken string `json:"refresh_token"`
}

type ChangePasswordBody struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type ChangeTimezoneBody struct {
	NewTimezone string `json:"new_timezone"`
}

type SignOutBody struct {
	FingerKey string `json:"finger_key"`
}

type ValidateBody struct {
	Token     string `json:"token"`
	FingerKey string `json:"finger_key"`
}

type SignUpBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
