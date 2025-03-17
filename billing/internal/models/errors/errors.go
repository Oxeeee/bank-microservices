package custerrors

import "errors"

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrUserNotFound        = errors.New("user not found")
	ErrPaymentNotFound     = errors.New("payment not found")
)
