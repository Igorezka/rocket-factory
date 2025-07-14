package inventory

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/samber/lo"

	repoModel "github.com/Igorezka/rocket-factory/inventory/internal/repository/model"
)

// fillTestData генерирует тестовые данные
func fillTestData(count int) map[string]repoModel.Part {
	data := make(map[string]repoModel.Part)

	for i := 0; i < count; i++ {
		id := uuid.NewString()

		part := repoModel.Part{
			Uuid:          id,
			Name:          gofakeit.Name(),
			Description:   gofakeit.Name(),
			Price:         gofakeit.Float64(),
			StockQuantity: gofakeit.Int64(),
			Category:      repoModel.Category(gofakeit.IntRange(0, 4)), //nolint:gosec // safe: gofakeit.IntRange returns 1..4
			Dimensions: repoModel.Dimensions{
				Length: gofakeit.Float64(),
				Width:  gofakeit.Float64(),
				Height: gofakeit.Float64(),
				Weight: gofakeit.Float64(),
			},
			Manufacturer: repoModel.Manufacturer{
				Name:    gofakeit.Company(),
				Country: gofakeit.Country(),
				Website: gofakeit.URL(),
			},
			Tags: []string{gofakeit.Name(), gofakeit.Company(), gofakeit.Country()},
			Metadata: repoModel.Metadata{
				Creator: gofakeit.Name(),
				Patent:  gofakeit.Int64(),
			},
			CreatedAt: time.Now(),
			UpdatedAt: lo.ToPtr(time.Now().Add(2 * time.Hour)),
		}
		data[id] = part
	}

	return data
}
