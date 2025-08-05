package inventory

import (
	"context"

	"github.com/Igorezka/rocket-factory/inventory/internal/model"
)

// Get получает деталь из хранилища по uuid и возвращает её
func (s *service) Get(ctx context.Context, uuid string) (model.Part, error) {
	part, err := s.partRepository.Get(ctx, uuid)
	if err != nil {
		return model.Part{}, err
	}

	return part, nil
}
