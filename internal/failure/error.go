package failure

import "errors"

var (
	ErrInput              = errors.New("bad request")
	ErrNotFound           = errors.New("not found")
	ErrPasswordNotCorrect = errors.New("incorrect password")
	ErrJWTGenerate        = errors.New("error to generate token")
	ErrInternal           = errors.New("internal error")
	ErrHashingPassword    = errors.New("error to hash password")
	ErrAuth               = errors.New("wrong token")
	ErrJWTNotValid        = errors.New("not walid token")
	ErrChangeTimezone     = errors.New("error to change timezone")
	ErrToGetUser          = errors.New("user not found")
)
