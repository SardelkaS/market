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
	ErrInput           = LogicError{code: 500, description: "bad request"}
	ErrNotFound        = LogicError{code: 404, description: "not found"}
	ErrAuth            = LogicError{code: 401, description: "unauthorized"}
	ErrToGetUser       = LogicError{code: 500, description: "user not found"}
	ErrServiceNotFound = LogicError{code: 404, description: "service not found"}
)
