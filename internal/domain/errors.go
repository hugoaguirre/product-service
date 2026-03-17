package domain

import "errors"

var (
	ErrProductNotFound = errors.New("product not found")
	ErrInvalidID       = errors.New("invalid product id")
)
