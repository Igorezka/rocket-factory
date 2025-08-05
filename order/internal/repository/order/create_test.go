package order

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/Igorezka/rocket-factory/order/internal/model"
)

func (s *RepositorySuite) TestCreate() {
	type args struct {
		part model.Order
	}

	part := model.Order{
		Uuid:       gofakeit.UUID(),
		UserUuid:   gofakeit.UUID(),
		PartUuids:  []string{gofakeit.UUID(), gofakeit.UUID()},
		TotalPrice: gofakeit.Float64(),
		Status:     model.OrderStatusPendingPayment,
	}

	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "success case",
			args: args{
				part: part,
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			err := s.repository.Create(s.ctx, tt.args.part)
			s.Require().Equal(tt.err, err)
		})
	}
}
