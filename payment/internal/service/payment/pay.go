package payment

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/Igorezka/rocket-factory/payment/internal/model"
)

// PayOrder производит оплату и возвращает uuid транзакции
func (s *service) PayOrder(_ context.Context, _ model.PayOrder) (string, error) {
	u := uuid.NewString()

	log.Printf("Оплата прошла успешно, transaction_uuid: %s\n", u)

	return u, nil
}
