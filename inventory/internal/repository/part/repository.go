package inventory

import (
	"sync"

	def "github.com/Igorezka/rocket-factory/inventory/internal/repository"
	repoModel "github.com/Igorezka/rocket-factory/inventory/internal/repository/model"
)

var _ def.PartRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data map[string]repoModel.Part
}

func NewRepository() *repository {
	return &repository{
		data: fillTestData(4),
	}
}
