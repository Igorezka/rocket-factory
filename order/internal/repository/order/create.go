package order

import (
	"context"

	"github.com/Igorezka/rocket-factory/order/internal/model"
	"github.com/Igorezka/rocket-factory/order/internal/repository/converter"
)

func (r *repository) Create(_ context.Context, order model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[order.Uuid] = converter.OrderToRepoModel(order)

	return nil
}
