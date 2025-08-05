package v1

import (
	"github.com/Igorezka/rocket-factory/inventory/internal/service"
	generatedInventoryV1 "github.com/Igorezka/rocket-factory/shared/pkg/proto/inventory/v1"
)

type api struct {
	generatedInventoryV1.UnimplementedInventoryServiceServer

	partService service.PartService
}

func NewAPI(partService service.PartService) *api {
	return &api{
		partService: partService,
	}
}
