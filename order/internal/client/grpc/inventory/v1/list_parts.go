package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	clientConverter "github.com/Igorezka/rocket-factory/order/internal/client/converter"
	"github.com/Igorezka/rocket-factory/order/internal/model"
	generatedInventoryV1 "github.com/Igorezka/rocket-factory/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	parts, err := c.generatedClient.ListParts(ctx, &generatedInventoryV1.ListPartsRequest{
		Filter: clientConverter.PartsFilterToProto(filter),
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, model.ErrPartsNotFound
		}

		return nil, err
	}

	return clientConverter.PartListToModel(parts.GetParts()), nil
}
