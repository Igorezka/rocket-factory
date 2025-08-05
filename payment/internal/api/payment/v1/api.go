package v1

import (
	"github.com/Igorezka/rocket-factory/payment/internal/service"
	generatedPaymentV1 "github.com/Igorezka/rocket-factory/shared/pkg/proto/payment/v1"
)

type api struct {
	generatedPaymentV1.UnimplementedPaymentServiceServer

	paymentService service.PaymentService
}

func NewAPI(paymentService service.PaymentService) *api {
	return &api{
		paymentService: paymentService,
	}
}
