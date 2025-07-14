package order

import (
	"context"

	"github.com/Igorezka/rocket-factory/order/internal/model"
)

func (s *service) Cancel(ctx context.Context, uuid string) error {
	order, err := s.orderRepository.Get(ctx, uuid)
	if err != nil {
		return err
	}

	// Проверяем статус заказа
	switch order.Status {
	case model.OrderStatusPaid:
		return model.ErrOrderAlreadyPaid
	case model.OrderStatusCancelled:
		return model.ErrOrderCancelled
	}

	// Обновляем статус заказа
	status := model.OrderStatusCancelled
	err = s.orderRepository.Update(ctx, uuid, model.OrderUpdate{
		Status: &status,
	})
	if err != nil {
		return err
	}

	return nil
}
