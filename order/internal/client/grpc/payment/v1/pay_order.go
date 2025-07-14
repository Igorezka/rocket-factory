package v1

import (
	"context"

	generatedPaymentV1 "github.com/Igorezka/rocket-factory/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, orderUuid, userUuid, paymentMethod string) (string, error) {
	res, err := c.generatedClient.PayOrder(ctx, &generatedPaymentV1.PayOrderRequest{
		OrderUuid:     orderUuid,
		UserUuid:      userUuid,
		PaymentMethod: generatedPaymentV1.PaymentMethod(generatedPaymentV1.PaymentMethod_value[paymentMethod]),
	})
	if err != nil {
		return "", err
	}

	return res.TransactionUuid, err
}
