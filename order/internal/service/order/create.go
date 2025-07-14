package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Igorezka/rocket-factory/order/internal/model"
)

func (s *service) Create(ctx context.Context, orderCreate model.OrderCreate) (string, float64, error) {
	// Создаем таймаут на обращение
	clientCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// Получаем список запчастей по uuid
	res, err := s.inventoryClient.ListParts(clientCtx, model.PartsFilter{
		Uuids: orderCreate.PartUuids,
	})
	if err != nil {
		return "", 0, err
	}

	// Создаем базовую информацию о заказе
	order := model.Order{
		Uuid:     uuid.NewString(),
		UserUuid: orderCreate.UserUuid,
		Status:   model.OrderStatusPendingPayment,
	}

	// Проверяем на наличие всех необходимых запчастей, при нахождении добавляем в заказ и плюсуем цену,
	// при не находе падаем в ошибку
	for _, partUuid := range orderCreate.PartUuids {
		var part *model.Part
		for _, p := range res {
			if p.Uuid == partUuid && p.StockQuantity > 0 {
				part = &p
			}
		}

		if part == nil {
			return "", 0, fmt.Errorf("%s - %w", partUuid, model.ErrPartNotFound)
		}

		order.PartUuids = append(order.PartUuids, partUuid)
		order.TotalPrice += part.Price
	}

	// Сохраняем заказ
	err = s.orderRepository.Create(ctx, order)
	if err != nil {
		return "", 0, err
	}

	return order.Uuid, order.TotalPrice, nil
}
