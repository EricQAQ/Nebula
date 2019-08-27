package err

import (
	"fmt"
)

type NebulaError struct {
	Code   int
	StrErr string
	Args   []interface{}
}

func CreateNebulaError(code int, err string, args []interface{}) *NebulaError {
	return &NebulaError{
		Code:   code,
		StrErr: err,
		Args:   args,
	}
}

func (err *NebulaError) getMessage() string {
	if len(err.Args) > 0 {
		return fmt.Sprintf(err.StrErr, err.Args...)
	}
	return err.StrErr
}

func (err *NebulaError) Error() string {
	return fmt.Sprintf("[ErrCode: %d]: %s", err.Code, err.getMessage())
}

func (err *NebulaError) FastGen(args ...interface{}) error {
	e := *err
	e.Args = args
	return &e
}


