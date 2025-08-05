package model

type Order struct {
	Uuid            string
	UserUuid        string
	PartUuids       []string
	TotalPrice      float64
	TransactionUuid *string
	PaymentMethod   *PaymentMethod
	Status          Status
}

type OrderUpdate struct {
	PartUuids       *[]string
	TotalPrice      *float64
	TransactionUuid *string
	PaymentMethod   *PaymentMethod
	Status          *Status
}

type (
	PaymentMethod string
	Status        string
)
