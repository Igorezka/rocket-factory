package model

import "errors"

var (
	ErrOrderNotFound    = errors.New("order not found")
	ErrPartsNotFound    = errors.New("parts not found")
	ErrPartNotFound     = errors.New("part not found")
	ErrOrderAlreadyPaid = errors.New("order already paid")
	ErrOrderCancelled   = errors.New("order cancelled")
)
