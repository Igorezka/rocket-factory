package order

import (
	"context"

	"github.com/Igorezka/rocket-factory/order/internal/model"
	repoModel "github.com/Igorezka/rocket-factory/order/internal/repository/model"
)

func (r *repository) Update(_ context.Context, uuid string, updateOrder model.OrderUpdate) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, ok := r.data[uuid]
	if !ok {
		return model.ErrOrderNotFound
	}

	if updateOrder.PartUuids != nil {
		order.PartUuids = *updateOrder.PartUuids
	}

	if updateOrder.TotalPrice != nil {
		order.TotalPrice = *updateOrder.TotalPrice
	}

	if updateOrder.TransactionUuid != nil {
		order.TransactionUuid = updateOrder.TransactionUuid
	}

	if updateOrder.PaymentMethod != nil {
		method := repoModel.PaymentMethod(*updateOrder.PaymentMethod)
		order.PaymentMethod = &method
	}

	if updateOrder.Status != nil {
		order.Status = repoModel.Status(*updateOrder.Status)
	}

	r.data[uuid] = order

	return nil
}
