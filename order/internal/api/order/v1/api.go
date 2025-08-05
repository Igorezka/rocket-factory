package v1

import (
	"context"
	"net/http"

	"github.com/Igorezka/rocket-factory/order/internal/service"
	generatedOrderV1 "github.com/Igorezka/rocket-factory/shared/pkg/openapi/order/v1"
)

type api struct {
	orderService service.OrderService
}

func NewAPI(orderService service.OrderService) *api {
	return &api{
		orderService: orderService,
	}
}

func (a *api) NewError(_ context.Context, err error) *generatedOrderV1.GenericErrorStatusCode {
	return &generatedOrderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: generatedOrderV1.GenericError{
			Code:    generatedOrderV1.NewOptInt(http.StatusInternalServerError),
			Message: generatedOrderV1.NewOptString(err.Error()),
		},
	}
}
