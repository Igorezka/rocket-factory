package converter

import (
	"github.com/Igorezka/rocket-factory/inventory/internal/model"
	repoModel "github.com/Igorezka/rocket-factory/inventory/internal/repository/model"
)

func PartsFilterToRepoModel(filter model.PartsFilter) repoModel.PartsFilter {
	return repoModel.PartsFilter{
		Uuids:                 filter.Uuids,
		Names:                 filter.Names,
		Categories:            CategoriesToRepoModel(filter.Categories),
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}

func CategoriesToRepoModel(categories []model.Category) []repoModel.Category {
	data := make([]repoModel.Category, len(categories))
	for _, category := range categories {
		data = append(data, repoModel.Category(category))
	}
	return data
}

func PartToModel(part repoModel.Part) model.Part {
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
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

func PartDimensionsToModel(dimensions repoModel.Dimensions) model.Dimensions {
	return model.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func PartManufacturerToModel(manufacturer repoModel.Manufacturer) model.Manufacturer {
	return model.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func PartMetadataToModel(metadata repoModel.Metadata) model.Metadata {
	return model.Metadata{
		Creator: metadata.Creator,
		Patent:  metadata.Patent,
	}
}
