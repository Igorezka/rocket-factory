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

// GetPart возвращает деталь по UUID
func (a *api) GetPart(ctx context.Context, req *generatedInventoryV1.GetPartRequest) (*generatedInventoryV1.GetPartResponse, error) {
	part, err := a.partService.Get(ctx, req.GetUuid())
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return nil, status.Errorf(codes.NotFound, "internal with UUID %s not found", req.GetUuid())
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &generatedInventoryV1.GetPartResponse{
		Part: converter.PartToProto(part),
	}, nil
}
