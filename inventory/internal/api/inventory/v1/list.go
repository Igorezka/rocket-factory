package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Igorezka/rocket-factory/inventory/internal/converter"
	"github.com/Igorezka/rocket-factory/inventory/internal/model"
	generatedInventoryV1 "github.com/Igorezka/rocket-factory/shared/pkg/proto/inventory/v1"
)

// ListParts возвращает список деталей соответствующих переданным фильтрам
// или возвращает все детали если фильтры не переданы
func (a *api) ListParts(ctx context.Context, req *generatedInventoryV1.ListPartsRequest) (*generatedInventoryV1.ListPartsResponse, error) {
	parts, err := a.partService.List(ctx, converter.PartsFilterToModel(req.GetFilter()))
	if err != nil {
		if errors.Is(err, model.ErrPartsNotFound) {
			return nil, status.Error(codes.NotFound, "no parts found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &generatedInventoryV1.ListPartsResponse{
		Parts: converter.PartsToProto(parts),
	}, nil
}
