package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Igorezka/rocket-factory/payment/internal/converter"
	generatedPaymentV1 "github.com/Igorezka/rocket-factory/shared/pkg/proto/payment/v1"
)

// PayOrder производит оплату и возвращает uuid транзакции
func (a *api) PayOrder(ctx context.Context, req *generatedPaymentV1.PayOrderRequest) (*generatedPaymentV1.PayOrderResponse, error) {
	transactionUuid, err := a.paymentService.PayOrder(ctx, converter.PayOrderToModel(req))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &generatedPaymentV1.PayOrderResponse{
		TransactionUuid: transactionUuid,
	}, nil
}
