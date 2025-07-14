package order

import (
	"context"
	"time"

	"github.com/Igorezka/rocket-factory/order/internal/model"
)

func (s *service) Pay(ctx context.Context, orderPay model.OrderPay) (string, error) {
	order, err := s.orderRepository.Get(ctx, orderPay.Uuid)
	if err != nil {
		return "", err
	}

	// Проверяем статус заказа
	switch order.Status {
	case model.OrderStatusPaid:
		return "", model.ErrOrderAlreadyPaid
	case model.OrderStatusCancelled:
		return "", model.ErrOrderCancelled
	}

	// Создаем таймаут на обращение
	clientCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// Оплачиваем заказ через payment service
	transactionUuid, err := s.paymentClient.PayOrder(clientCtx, orderPay.Uuid, order.UserUuid, string(orderPay.PaymentMethod))
	if err != nil {
		return "", err
	}

	status := model.OrderStatusPaid
	// Обновляем платежную информацию
	err = s.orderRepository.Update(ctx, orderPay.Uuid, model.OrderUpdate{
		TransactionUuid: &transactionUuid,
		PaymentMethod:   &orderPay.PaymentMethod,
		Status:          &status,
	})
	if err != nil {
		return "", err
	}

	return transactionUuid, nil
}
