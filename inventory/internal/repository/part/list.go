package inventory

import (
	"context"
	"slices"

	"github.com/Igorezka/rocket-factory/inventory/internal/model"
	repoConverter "github.com/Igorezka/rocket-factory/inventory/internal/repository/converter"
	repoModel "github.com/Igorezka/rocket-factory/inventory/internal/repository/model"
)

// filterFunc функция по которой фильтруется деталь
type filterFunc func(part repoModel.Part) bool

// List возвращает детали отфильтрованные в соответствии с переданным фильтром
func (r *repository) List(_ context.Context, filter model.PartsFilter) ([]model.Part, error) {
	// Копируем детали в локальную переменную
	r.mu.RLock()
	parts := make(map[string]repoModel.Part)
	for k, v := range r.data {
		parts[k] = v
	}
	r.mu.RUnlock()

	// Конвертируем модель фильтра и инициализируем список
	repoFilter := repoConverter.PartsFilterToRepoModel(filter)
	filterFuncs := newFilterFuncs(repoFilter)

	filteredParts := make([]model.Part, 0)
	for _, part := range parts {
		match := true
		for _, f := range filterFuncs {
			if !f(part) {
				match = false
				break
			}
		}
		if match {
			filteredParts = append(filteredParts, repoConverter.PartToModel(part))
		}
	}

	if len(filteredParts) == 0 {
		return nil, model.ErrPartsNotFound
	}

	return filteredParts, nil
}

// newFilterFuncs создает список функций по которым будут фильтроваться детали
func newFilterFuncs(filter repoModel.PartsFilter) []filterFunc {
	return []filterFunc{
		func(part repoModel.Part) bool {
			if len(filter.Uuids) == 0 {
				return true
			}
			return slices.Contains(filter.Uuids, part.Uuid)
		},
		func(part repoModel.Part) bool {
			if len(filter.Names) == 0 {
				return true
			}
			return slices.Contains(filter.Names, part.Name)
		},
		func(part repoModel.Part) bool {
			if len(filter.Categories) == 0 {
				return true
			}
			return slices.Contains(filter.Categories, part.Category)
		},
		func(part repoModel.Part) bool {
			if len(filter.ManufacturerCountries) == 0 {
				return true
			}
			return slices.Contains(filter.ManufacturerCountries, part.Manufacturer.Country)
		},
		func(part repoModel.Part) bool {
			if len(filter.Tags) == 0 {
				return true
			}

			if len(part.Tags) == 0 {
				return false
			}
			for _, tag := range filter.Tags {
				if slices.Contains(part.Tags, tag) {
					return true
				}
			}
			return false
		},
	}
}
