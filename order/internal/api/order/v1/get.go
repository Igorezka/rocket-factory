package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/Igorezka/rocket-factory/order/internal/converter"
	"github.com/Igorezka/rocket-factory/order/internal/model"
	generatedOrderV1 "github.com/Igorezka/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrderByUUID(ctx context.Context, params generatedOrderV1.GetOrderByUUIDParams) (generatedOrderV1.GetOrderByUUIDRes, error) {
	order, err := a.orderService.Get(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &generatedOrderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "Order by UUID " + params.OrderUUID + " not found",
			}, nil
		}

		return &generatedOrderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}, nil
	}

	return converter.OrderToOpenAPI(order), nil
}
