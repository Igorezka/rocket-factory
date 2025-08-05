package converter

import (
	"github.com/Igorezka/rocket-factory/order/internal/model"
	generatedOrderV1 "github.com/Igorezka/rocket-factory/shared/pkg/openapi/order/v1"
)

func OrderPayToModel(uuid string, paymentMethod generatedOrderV1.PaymentMethod) model.OrderPay {
	return model.OrderPay{
		Uuid:          uuid,
		PaymentMethod: model.PaymentMethod(paymentMethod),
	}
}

func OrderToOpenAPI(order model.Order) *generatedOrderV1.OrderDto {
	var transaction generatedOrderV1.OptString
	if order.TransactionUuid != nil {
		transaction = generatedOrderV1.OptString{Value: *order.TransactionUuid, Set: true}
	}

	var method generatedOrderV1.OptPaymentMethod
	if order.PaymentMethod != nil {
		method = generatedOrderV1.OptPaymentMethod{Value: generatedOrderV1.PaymentMethod(*order.PaymentMethod), Set: true}
	}
	return &generatedOrderV1.OrderDto{
		OrderUUID:       order.Uuid,
		UserUUID:        order.UserUuid,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transaction,
		PaymentMethod:   method,
		Status:          generatedOrderV1.OrderStatus(order.Status),
	}
}
