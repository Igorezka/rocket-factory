package grpc

import (
	"context"

	"github.com/Igorezka/rocket-factory/order/internal/model"
)

type InventoryClient interface {
	ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUuid, userUuid, paymentMethod string) (string, error)
}
