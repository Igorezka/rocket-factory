package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Igorezka/rocket-factory/inventory/internal/model"
	generatedInventoryV1 "github.com/Igorezka/rocket-factory/shared/pkg/proto/inventory/v1"
)

func PartsFilterToModel(filter *generatedInventoryV1.PartsFilter) model.PartsFilter {
	return model.PartsFilter{
		Uuids:                 filter.GetUuids(),
		Names:                 filter.GetNames(),
		Categories:            CategoriesToModel(filter.GetCategories()),
		ManufacturerCountries: filter.GetManufacturerCountries(),
		Tags:                  filter.GetTags(),
	}
}

func CategoriesToModel(categories []generatedInventoryV1.Category) []model.Category {
	data := make([]model.Category, len(categories))
	for _, category := range categories {
		data = append(data, model.Category(category))
	}
	return data
}

func PartToProto(part model.Part) *generatedInventoryV1.Part {
	var updatedAt *timestamppb.Timestamp
	if part.UpdatedAt != nil {
		updatedAt = timestamppb.New(*part.UpdatedAt)
	}

	return &generatedInventoryV1.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      generatedInventoryV1.Category(part.Category),
		Dimensions: &generatedInventoryV1.Dimensions{
			Length: part.Dimensions.Length,
			Width:  part.Dimensions.Width,
			Height: part.Dimensions.Height,
			Weight: part.Dimensions.Weight,
		},
		Manufacturer: &generatedInventoryV1.Manufacturer{
			Name:    part.Manufacturer.Name,
			Country: part.Manufacturer.Country,
			Website: part.Manufacturer.Website,
		},
		Tags: part.Tags,
		Metadata: map[string]*generatedInventoryV1.Value{
			"creator": {
				ValueType: &generatedInventoryV1.Value_StringValue{StringValue: part.Metadata.Creator},
			},
			"patent": {
				ValueType: &generatedInventoryV1.Value_Int64Value{Int64Value: part.Metadata.Patent},
			},
		},
		CreatedAt: timestamppb.New(part.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func PartsToProto(parts []model.Part) []*generatedInventoryV1.Part {
	data := make([]*generatedInventoryV1.Part, len(parts))
	for i, part := range parts {
		data[i] = PartToProto(part)
	}
	return data
}
