package inventory

import (
	"context"

	"github.com/Igorezka/rocket-factory/inventory/internal/model"
	repoConverter "github.com/Igorezka/rocket-factory/inventory/internal/repository/converter"
)

// Get возвращает деталь по uuid
func (r *repository) Get(_ context.Context, uuid string) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	part, ok := r.data[uuid]
	if !ok {
		return model.Part{}, model.ErrPartNotFound
	}

	return repoConverter.PartToModel(part), nil
}
