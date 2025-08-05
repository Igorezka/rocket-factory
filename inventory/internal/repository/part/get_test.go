package inventory

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/Igorezka/rocket-factory/inventory/internal/model"
	"github.com/Igorezka/rocket-factory/inventory/internal/repository/converter"
)

func (s *RepositorySuite) TestGet() {
	type args struct {
		uuid string
	}

	var (
		randomPartIndex = gofakeit.IntRange(0, len(s.generatedParts)-1)
		part            = s.generatedParts[randomPartIndex]

		expectedPart = converter.PartToModel(part)
	)

	tests := []struct {
		name string
		args args
		want model.Part
		err  error
	}{
		{
			name: "success case",
			args: args{
				uuid: part.Uuid,
			},
			want: expectedPart,
			err:  nil,
		},
		{
			name: "error case",
			args: args{
				uuid: uuid.NewString(),
			},
			want: model.Part{},
			err:  model.ErrPartNotFound,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			res, err := s.repository.Get(s.ctx, tt.args.uuid)
			s.Require().Equal(tt.want, res)
			s.Require().Equal(tt.err, err)
		})
	}
}
