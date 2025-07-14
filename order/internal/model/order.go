package model

const (
	PaymentMethodCard PaymentMethod = "PAYMENT_METHOD_CARD"

	OrderStatusPendingPayment Status = "PENDING_PAYMENT"
	OrderStatusPaid           Status = "PAID"
	OrderStatusCancelled      Status = "CANCELLED"
)

type Order struct {
	Uuid            string
	UserUuid        string
	PartUuids       []string
	TotalPrice      float64
	TransactionUuid *string
	PaymentMethod   *PaymentMethod
	Status          Status
}

type OrderCreate struct {
	UserUuid  string
	PartUuids []string
}

type OrderUpdate struct {
	PartUuids       *[]string
	TotalPrice      *float64
	TransactionUuid *string
	PaymentMethod   *PaymentMethod
	Status          *Status
}

type OrderPay struct {
	Uuid          string
	PaymentMethod PaymentMethod
}

type (
	PaymentMethod string
	Status        string
)
