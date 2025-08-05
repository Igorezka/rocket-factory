package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/Igorezka/rocket-factory/order/internal/converter"
	"github.com/Igorezka/rocket-factory/order/internal/model"
	generatedOrderV1 "github.com/Igorezka/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *generatedOrderV1.PayOrderRequest, params generatedOrderV1.PayOrderParams) (generatedOrderV1.PayOrderRes, error) {
	transactionUuid, err := a.orderService.Pay(ctx, converter.OrderPayToModel(params.OrderUUID, req.PaymentMethod))
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &generatedOrderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "Order by UUID " + params.OrderUUID + " not found",
			}, nil
		}
		if errors.Is(err, model.ErrOrderAlreadyPaid) {
			return &generatedOrderV1.NotFoundError{
				Code:    http.StatusConflict,
				Message: "Order UUID " + params.OrderUUID + " already paid",
			}, nil
		}
		if errors.Is(err, model.ErrOrderCancelled) {
			return &generatedOrderV1.NotFoundError{
				Code:    http.StatusConflict,
				Message: "Order UUID " + params.OrderUUID + " canceled",
			}, nil
		}

		return &generatedOrderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}, nil
	}

	return &generatedOrderV1.PayOrderResponse{
		TransactionUUID: transactionUuid,
	}, nil
}
