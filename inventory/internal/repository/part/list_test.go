package inventory

import (
	"slices"

	"github.com/google/uuid"

	"github.com/Igorezka/rocket-factory/inventory/internal/model"
	"github.com/Igorezka/rocket-factory/inventory/internal/repository/converter"
)

func (s *RepositorySuite) TestList() {
	type args struct {
		filter model.PartsFilter
	}

	var (
		uuids         []string
		names         []string
		categories    []model.Category
		countries     []string
		expectedParts []model.Part
	)

	for _, part := range s.generatedParts {
		tmp := converter.PartToModel(part)
		expectedParts = append(expectedParts, tmp)
		uuids = append(uuids, tmp.Uuid)
		names = append(names, tmp.Name)
		if !slices.Contains(categories, tmp.Category) {
			categories = append(categories, tmp.Category)
		}
		if !slices.Contains(countries, tmp.Manufacturer.Country) {
			countries = append(countries, tmp.Manufacturer.Country)
		}
	}

	tests := []struct {
		name string
		args args
		want []model.Part
		err  error
	}{
		{
			name: "success by uuids case",
			args: args{
				filter: model.PartsFilter{
					Uuids: uuids,
				},
			},
			want: expectedParts,
			err:  nil,
		},
		{
			name: "success by names case",
			args: args{
				filter: model.PartsFilter{
					Names: names,
				},
			},
			want: expectedParts,
			err:  nil,
		},
		{
			name: "success by categories case",
			args: args{
				filter: model.PartsFilter{
					Categories: categories,
				},
			},
			want: expectedParts,
			err:  nil,
		},
		{
			name: "success by countries case",
			args: args{
				filter: model.PartsFilter{
					ManufacturerCountries: countries,
				},
			},
			want: expectedParts,
			err:  nil,
		},
		{
			name: "error case",
			args: args{
				filter: model.PartsFilter{
					Uuids: []string{uuid.NewString()},
				},
			},
			want: nil,
			err:  model.ErrPartsNotFound,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			res, err := s.repository.List(s.ctx, tt.args.filter)
			if err != nil {
				s.Require().Nil(res)
			} else {
				for _, item := range res {
					s.Require().Contains(tt.want, item)
				}
			}
			s.Require().Equal(tt.err, err)
		})
	}
}
