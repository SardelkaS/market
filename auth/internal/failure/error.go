package failure

type LogicError struct {
	code        int
	err         error
	description string
}

func (l LogicError) Code() int {
	return l.code
}

func (l LogicError) Description() string {
	return l.description
}

func (l LogicError) Error() string {
	return l.err.Error()
}

func (l LogicError) Wrap(err error) LogicError {
	l.err = err
	return l
}

var (
	ErrInput              = LogicError{code: 500, description: "bad request"}
	ErrNotFound           = LogicError{code: 404, description: "not found"}
	ErrPasswordNotCorrect = LogicError{code: 500, description: "incorrect password"}
	ErrJWTGenerate        = LogicError{code: 500, description: "error to generate token"}
	ErrInternal           = LogicError{code: 500, description: "internal error"}
	ErrHashingPassword    = LogicError{code: 500, description: "error to hash password"}
	ErrAuth               = LogicError{code: 500, description: "wrong token"}
	ErrJWTNotValid        = LogicError{code: 500, description: "not walid token"}
	ErrChangeTimezone     = LogicError{code: 500, description: "error to change timezone"}
	ErrToGetUser          = LogicError{code: 404, description: "user not found"}
	ErrToUpdateUser       = LogicError{code: 500, description: "error to update user"}
)
