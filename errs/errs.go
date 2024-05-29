package errs

// Error instance
type Error struct {
	err        error
	message    string
	statusCode int
}

// NewError returns a new strandard error
func NewError(err error, statusCode int, message string, moreInfos ...string) error {
	return &Error{
		err:        err,
		statusCode: statusCode,
		message:    message,
	}
}

func NewErrsWithCode(err error, codeErr CodeErr, moreInfos ...string) error {
	return &Error{
		err:        err,
		statusCode: codeErr.StatusCode(),
		message:    codeErr.Message(),
	}
}

// ParseError returns an instance of Error
func ParseError(err error) *Error {
	switch r := err.(type) {
	case *Error:
		return r
	default:
		return nil
	}
}

func (e *Error) Error() string {
	return e.err.Error()
}

func (e *Error) Message() string {
	return e.message
}

func (e *Error) StatusCode() int {
	return e.statusCode
}
