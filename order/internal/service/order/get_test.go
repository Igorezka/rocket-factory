package order

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/Igorezka/rocket-factory/order/internal/model"
)

func (s *ServiceSuite) TestGet() {
	type args struct {
		uuid string
	}

	var (
		uuid = gofakeit.UUID()

		order = model.Order{
			Uuid:       uuid,
			UserUuid:   gofakeit.UUID(),
			PartUuids:  []string{gofakeit.UUID(), gofakeit.UUID()},
			TotalPrice: gofakeit.Float64(),
			Status:     model.OrderStatusCancelled,
		}
	)

	tests := []struct {
		name                         string
		args                         args
		want                         model.Order
		err                          error
		orderRepositoryMockConfigure func()
	}{
		{
			name: "success case",
			args: args{
				uuid: uuid,
			},
			want: order,
			err:  nil,
			orderRepositoryMockConfigure: func() {
				s.orderRepository.On("Get", s.ctx, uuid).Return(order, nil).Once()
			},
		},
		{
			name: "error case",
			args: args{
				uuid: uuid,
			},
			want: model.Order{},
			err:  model.ErrOrderNotFound,
			orderRepositoryMockConfigure: func() {
				s.orderRepository.On("Get", s.ctx, uuid).Return(model.Order{}, model.ErrOrderNotFound).Once()
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.orderRepositoryMockConfigure()
			res, err := s.service.Get(s.ctx, tt.args.uuid)
			s.Require().Equal(tt.want, res)
			s.Require().Equal(tt.err, err)
		})
	}
}
