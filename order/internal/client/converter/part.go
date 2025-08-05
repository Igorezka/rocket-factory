package converter

import (
	"time"

	"github.com/Igorezka/rocket-factory/order/internal/model"
	generatedInventoryV1 "github.com/Igorezka/rocket-factory/shared/pkg/proto/inventory/v1"
)

func PartsFilterToProto(filter model.PartsFilter) *generatedInventoryV1.PartsFilter {
	return &generatedInventoryV1.PartsFilter{
		Uuids:                 filter.Uuids,
		Names:                 filter.Names,
		Categories:            CategoriesToModel(filter.Categories),
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}

func CategoriesToModel(categories []model.Category) []generatedInventoryV1.Category {
	data := make([]generatedInventoryV1.Category, len(categories))
	for _, category := range categories {
		data = append(data, generatedInventoryV1.Category(category))
	}
	return data
}

func PartToModel(part *generatedInventoryV1.Part) model.Part {
	var updatedAt *time.Time
	if part.UpdatedAt != nil {
		tmp := part.UpdatedAt.AsTime()
		updatedAt = &tmp
	}

	return model.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      model.Category(part.Category),
		Dimensions:    PartDimensionsToModel(part.Dimensions),
		Manufacturer:  PartManufacturerToModel(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      PartMetadataToModel(part.Metadata),
		CreatedAt:     part.CreatedAt.AsTime(),
		UpdatedAt:     updatedAt,
	}
}

func PartDimensionsToModel(dimensions *generatedInventoryV1.Dimensions) model.Dimensions {
	return model.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func PartManufacturerToModel(manufacturer *generatedInventoryV1.Manufacturer) model.Manufacturer {
	return model.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func PartMetadataToModel(metadata map[string]*generatedInventoryV1.Value) model.Metadata {
	return model.Metadata{
		Creator: metadata["creator"].String(),
		Patent:  metadata["patent"].GetInt64Value(),
	}
}

func PartListToModel(parts []*generatedInventoryV1.Part) []model.Part {
	res := make([]model.Part, 0, len(parts))
	for _, part := range parts {
		res = append(res, PartToModel(part))
	}

	return res
}
