package v1

import (
	"context"
	"net/http"

	"github.com/go-faster/errors"

	"github.com/Igorezka/rocket-factory/order/internal/model"
	generatedOrderV1 "github.com/Igorezka/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params generatedOrderV1.CancelOrderParams) (generatedOrderV1.CancelOrderRes, error) {
	err := a.orderService.Cancel(ctx, params.OrderUUID)
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
				Message: "Order UUID " + params.OrderUUID + " already paid and cannot be cancelled",
			}, nil
		}
		if errors.Is(err, model.ErrOrderCancelled) {
			return &generatedOrderV1.NotFoundError{
				Code:    http.StatusConflict,
				Message: "Order UUID " + params.OrderUUID + " has already been canceled",
			}, nil
		}

		return &generatedOrderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}, nil
	}

	return &generatedOrderV1.CancelOrderNoContent{}, nil
}
