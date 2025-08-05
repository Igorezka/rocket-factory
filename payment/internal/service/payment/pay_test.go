package payment

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/Igorezka/rocket-factory/payment/internal/model"
)

func (s *ServiceSuite) TestPay() {
	type args struct {
		req model.PayOrder
	}

	var (
		orderUuid     = uuid.NewString()
		userUuid      = uuid.NewString()
		paymentMethod = gofakeit.IntRange(0, 4)

		req = model.PayOrder{
			OrderUuid:     orderUuid,
			UserUuid:      userUuid,
			PaymentMethod: model.PaymentMethod(paymentMethod),
		}
	)

	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "success case",
			args: args{req: req},
			err:  nil,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			res, err := s.service.PayOrder(s.ctx, tt.args.req)
			s.Require().NoError(err)
			s.Require().NotEmpty(res)
		})
	}
}
