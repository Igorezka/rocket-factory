package model

type PayOrder struct {
	OrderUuid     string
	UserUuid      string
	PaymentMethod PaymentMethod
}

type PaymentMethod int32
