package inventory

import (
	"context"

	"github.com/Igorezka/rocket-factory/inventory/internal/model"
)

// List получает детали из хранилища отфильтрованные согласно переданному фильтру
func (s *service) List(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	parts, err := s.partRepository.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	return parts, nil
}
