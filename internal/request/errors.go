package request

import "errors"

var ErrInvalidBool = errors.New("не корректное значение для булевого типа")

type BadRequestError struct {
	Status int
	Err    error
}

func NewBadRequestErorr(status int, err error) BadRequestError {
	return BadRequestError{
		Status: status,
		Err:    err,
	}
}

func (e BadRequestError) Error() string {
	return e.Err.Error()
}

func (e BadRequestError) Unwrap() error {
	return e.Err
}
