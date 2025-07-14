package converter

import (
	"github.com/Igorezka/rocket-factory/payment/internal/model"
	generatedPaymentV1 "github.com/Igorezka/rocket-factory/shared/pkg/proto/payment/v1"
)

func PayOrderToModel(req *generatedPaymentV1.PayOrderRequest) model.PayOrder {
	return model.PayOrder{
		OrderUuid:     req.GetOrderUuid(),
		UserUuid:      req.GetUserUuid(),
		PaymentMethod: model.PaymentMethod(req.GetPaymentMethod()),
	}
}
