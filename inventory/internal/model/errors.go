package model

import "errors"

var (
	ErrPartNotFound  = errors.New("internal not found")
	ErrPartsNotFound = errors.New("parts not found")
)
