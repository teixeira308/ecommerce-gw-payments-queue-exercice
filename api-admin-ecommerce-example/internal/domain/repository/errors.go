package repository

import (
	"errors"
	"fmt"
)

var (
	ErrPaymentNotFound = errors.New("payment not found")
	ErrInternal        = errors.New("internal server error")
)

type ErrNotFound struct {
	Message string
}

func (e *ErrNotFound) Error() string {
	return e.Message
}

func NewErrNotFound(id string) error {
	return fmt.Errorf("%w: %s", ErrPaymentNotFound, id)
}
