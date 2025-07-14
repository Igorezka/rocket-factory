package service

import (
	"context"

	"github.com/Igorezka/rocket-factory/payment/internal/model"
)

type PaymentService interface {
	PayOrder(ctx context.Context, payOrder model.PayOrder) (string, error)
}
