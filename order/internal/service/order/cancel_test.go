package order

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/Igorezka/rocket-factory/order/internal/model"
)

func (s *ServiceSuite) TestCancel() {
	type args struct {
		uuid string
	}

	var (
		uuid  = gofakeit.UUID()
		order = model.Order{
			Uuid:      uuid,
			UserUuid:  gofakeit.UUID(),
			PartUuids: []string{gofakeit.UUID(), gofakeit.UUID()},
			Status:    model.OrderStatusPendingPayment,
		}

		orderPaid = model.Order{
			Uuid:      uuid,
			UserUuid:  gofakeit.UUID(),
			PartUuids: []string{gofakeit.UUID(), gofakeit.UUID()},
			Status:    model.OrderStatusPaid,
		}

		orderCanceled = model.Order{
			Uuid:      uuid,
			UserUuid:  gofakeit.UUID(),
			PartUuids: []string{gofakeit.UUID(), gofakeit.UUID()},
			Status:    model.OrderStatusCancelled,
		}

		status      = model.OrderStatusCancelled
		orderUpdate = model.OrderUpdate{
			Status: &status,
		}
	)

	tests := []struct {
		name                         string
		args                         args
		err                          error
		orderRepositoryMockConfigure func()
	}{
		{
			name: "success case",
			args: args{
				uuid: uuid,
			},
			err: nil,
			orderRepositoryMockConfigure: func() {
				s.orderRepository.On("Get", s.ctx, uuid).Return(order, nil).Once()
				s.orderRepository.On("Update", s.ctx, uuid, orderUpdate).Return(nil).Once()
			},
		},
		{
			name: "not found error case",
			args: args{
				uuid: uuid,
			},
			err: model.ErrOrderNotFound,
			orderRepositoryMockConfigure: func() {
				s.orderRepository.On("Get", s.ctx, uuid).Return(model.Order{}, model.ErrOrderNotFound).Once()
			},
		},
		{
			name: "order already paid error case",
			args: args{
				uuid: uuid,
			},
			err: model.ErrOrderAlreadyPaid,
			orderRepositoryMockConfigure: func() {
				s.orderRepository.On("Get", s.ctx, uuid).Return(orderPaid, nil).Once()
			},
		},
		{
			name: "order already canceled error case",
			args: args{
				uuid: uuid,
			},
			err: model.ErrOrderCancelled,
			orderRepositoryMockConfigure: func() {
				s.orderRepository.On("Get", s.ctx, uuid).Return(orderCanceled, nil).Once()
			},
		},
		{
			name: "update not found error case",
			args: args{
				uuid: uuid,
			},
			err: model.ErrOrderNotFound,
			orderRepositoryMockConfigure: func() {
				s.orderRepository.On("Get", s.ctx, uuid).Return(order, nil).Once()
				s.orderRepository.On("Update", s.ctx, uuid, orderUpdate).Return(model.ErrOrderNotFound).Once()
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.orderRepositoryMockConfigure()

			err := s.service.Cancel(s.ctx, tt.args.uuid)
			s.Require().Equal(tt.err, err)
		})
	}
}
