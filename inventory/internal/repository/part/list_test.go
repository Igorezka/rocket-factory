package inventory

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/Igorezka/rocket-factory/inventory/internal/model"
	"github.com/Igorezka/rocket-factory/inventory/internal/repository/converter"
)

func (s *RepositorySuite) TestList() {
	type args struct {
		filter model.PartsFilter
	}

	var (
		filter        model.PartsFilter
		expectedParts []model.Part
	)

	for _, part := range s.generatedParts {
		if gofakeit.Bool() {
			tmp := converter.PartToModel(part)
			expectedParts = append(expectedParts, tmp)
			filter.Names = append(filter.Names, tmp.Name)
		}
	}

	tests := []struct {
		name string
		args args
		want []model.Part
		err  error
	}{
		{
			name: "success",
			args: args{
				filter: model.PartsFilter{
					Names: filter.Names,
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
