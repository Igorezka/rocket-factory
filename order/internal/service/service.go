package service

import (
	"context"

	"github.com/Igorezka/rocket-factory/order/internal/model"
)

type OrderService interface {
	Get(ctx context.Context, uuid string) (model.Order, error)
	Create(ctx context.Context, orderCreate model.OrderCreate) (string, float64, error)
	Pay(ctx context.Context, orderPay model.OrderPay) (string, error)
	Cancel(ctx context.Context, uuid string) error
}
