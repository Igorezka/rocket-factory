package repository

import (
	"context"

	"github.com/Igorezka/rocket-factory/order/internal/model"
)

type OrderRepository interface {
	Get(ctx context.Context, uuid string) (model.Order, error)
	Create(ctx context.Context, order model.Order) error
	Update(ctx context.Context, uuid string, updateOrder model.OrderUpdate) error
}
