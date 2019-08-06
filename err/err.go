package err

import (
	"fmt"
)

type TraedError struct {
	Code   int
	StrErr string
	Args   []interface{}
}

func CreateTraedError(code int, err string, args []interface{}) *TraedError {
	return &TraedError{
		Code:   code,
		StrErr: err,
		Args:   args,
	}
}

func (err *TraedError) getMessage() string {
	if len(err.Args) > 0 {
		return fmt.Sprintf(err.StrErr, err.Args...)
	}
	return err.StrErr
}

func (err *TraedError) Error() string {
	return fmt.Sprintf("[ErrCode: %d]: %s", err.Code, err.getMessage())
}

func (err *TraedError) FastGen(args ...interface{}) error {
	e := *err
	e.Args = args
	return &e
}


