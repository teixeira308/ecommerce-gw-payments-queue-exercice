package errors

import "errors"

var (
	ErrPaymentNotFound      = errors.New("payment not found")
	ErrOrderAlreadyProcessed = errors.New("order already processed")
	ErrInvalidPaymentData   = errors.New("invalid payment data")
	ErrInternal             = errors.New("internal server error")
)
