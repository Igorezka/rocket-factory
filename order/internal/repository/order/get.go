package order

import (
	"context"

	"github.com/Igorezka/rocket-factory/order/internal/model"
	"github.com/Igorezka/rocket-factory/order/internal/repository/converter"
)

func (r *repository) Get(_ context.Context, uuid string) (model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.data[uuid]
	if !ok {
		return model.Order{}, model.ErrOrderNotFound
	}

	return converter.OrderToModel(order), nil
}
