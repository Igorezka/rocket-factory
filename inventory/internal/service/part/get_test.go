package inventory

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"

	"github.com/Igorezka/rocket-factory/inventory/internal/model"
)

func (s *ServiceSuite) TestGet() {
	type args struct {
		uuid string
	}

	var (
		uuid = gofakeit.UUID()

		part = model.Part{
			Uuid:          uuid,
			Name:          gofakeit.Name(),
			Description:   gofakeit.Name(),
			Price:         gofakeit.Float64(),
			StockQuantity: gofakeit.Int64(),
			Category:      model.Category(gofakeit.IntRange(0, 4)), //nolint:gosec // safe: gofakeit.IntRange returns 1..4
			Dimensions: model.Dimensions{
				Length: gofakeit.Float64(),
				Width:  gofakeit.Float64(),
				Height: gofakeit.Float64(),
				Weight: gofakeit.Float64(),
			},
			Manufacturer: model.Manufacturer{
				Name:    gofakeit.Company(),
				Country: gofakeit.Country(),
				Website: gofakeit.URL(),
			},
			Tags: []string{gofakeit.Name(), gofakeit.Company(), gofakeit.Country()},
			Metadata: model.Metadata{
				Creator: gofakeit.Name(),
				Patent:  gofakeit.Int64(),
			},
			CreatedAt: time.Now(),
			UpdatedAt: lo.ToPtr(time.Now().Add(2 * time.Hour)),
		}
	)

	tests := []struct {
		name                        string
		args                        args
		want                        model.Part
		err                         error
		partRepositoryMockConfigure func()
	}{
		{
			name: "success case",
			args: args{
				uuid: uuid,
			},
			want: part,
			err:  nil,
			partRepositoryMockConfigure: func() {
				s.partRepository.On("Get", s.ctx, uuid).Return(part, nil).Once()
			},
		},
		{
			name: "error case",
			args: args{
				uuid: uuid,
			},
			want: model.Part{},
			err:  model.ErrPartNotFound,
			partRepositoryMockConfigure: func() {
				s.partRepository.On("Get", s.ctx, uuid).Return(model.Part{}, model.ErrPartNotFound).Once()
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.partRepositoryMockConfigure()
			res, err := s.service.Get(s.ctx, tt.args.uuid)
			s.Require().Equal(tt.want, res)
			s.Require().Equal(tt.err, err)
		})
	}
}
