package v1

import (
	"context"
	"net/http"

	"github.com/go-faster/errors"

	"github.com/Igorezka/rocket-factory/order/internal/model"
	generatedOrderV1 "github.com/Igorezka/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *generatedOrderV1.CreateOrderRequest) (generatedOrderV1.CreateOrderRes, error) {
	if len(req.PartUuids) == 0 {
		return &generatedOrderV1.InternalServerError{
			Code:    http.StatusBadRequest,
			Message: "Details not provided",
		}, nil
	}

	orderUuid, totalPrice, err := a.orderService.Create(ctx, model.OrderCreate{
		UserUuid:  req.GetUserUUID(),
		PartUuids: req.GetPartUuids(),
	})
	if err != nil {
		if errors.Is(err, model.ErrPartsNotFound) || errors.Is(err, model.ErrPartNotFound) {
			return &generatedOrderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			}, nil
		}

		return &generatedOrderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}, nil
	}

	return &generatedOrderV1.CreateOrderResponse{
		OrderUUID:  orderUuid,
		TotalPrice: totalPrice,
	}, nil
}
