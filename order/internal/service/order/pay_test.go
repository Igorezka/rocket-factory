package order

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"

	"github.com/Igorezka/rocket-factory/order/internal/model"
)

func (s *ServiceSuite) TestPay() {
	type args struct {
		orderPay model.OrderPay
	}

	var (
		uuid            = gofakeit.UUID()
		transactionUuid = gofakeit.UUID()

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

		orderPay = model.OrderPay{
			Uuid:          uuid,
			PaymentMethod: model.PaymentMethodCard,
		}

		status      = model.OrderStatusPaid
		orderUpdate = model.OrderUpdate{
			TransactionUuid: &transactionUuid,
			PaymentMethod:   &orderPay.PaymentMethod,
			Status:          &status,
		}
	)

	tests := []struct {
		name           string
		args           args
		want           string
		err            error
		mocksConfigure func()
	}{
		{
			name: "success case",
			args: args{
				orderPay: orderPay,
			},
			want: transactionUuid,
			err:  nil,
			mocksConfigure: func() {
				s.orderRepository.On("Get", s.ctx, uuid).Return(order, nil).Once()
				s.paymentClient.On("PayOrder", mock.Anything, orderPay.Uuid, order.UserUuid, string(orderPay.PaymentMethod)).
					Return(transactionUuid, nil).Once()
				s.orderRepository.On("Update", s.ctx, uuid, orderUpdate).Return(nil).Once()
			},
		},
		{
			name: "not found error case",
			args: args{
				orderPay: orderPay,
			},
			want: "",
			err:  model.ErrOrderNotFound,
			mocksConfigure: func() {
				s.orderRepository.On("Get", s.ctx, uuid).Return(model.Order{}, model.ErrOrderNotFound).Once()
			},
		},
		{
			name: "order already paid error case",
			args: args{
				orderPay: orderPay,
			},
			err: model.ErrOrderAlreadyPaid,
			mocksConfigure: func() {
				s.orderRepository.On("Get", s.ctx, uuid).Return(orderPaid, nil).Once()
			},
		},
		{
			name: "order already canceled error case",
			args: args{
				orderPay: orderPay,
			},
			err: model.ErrOrderCancelled,
			mocksConfigure: func() {
				s.orderRepository.On("Get", s.ctx, uuid).Return(orderCanceled, nil).Once()
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mocksConfigure()
			res, err := s.service.Pay(s.ctx, tt.args.orderPay)
			s.Require().Equal(tt.want, res)
			s.Require().Equal(tt.err, err)
		})
	}
}
