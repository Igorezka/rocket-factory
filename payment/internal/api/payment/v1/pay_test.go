package v1

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Igorezka/rocket-factory/payment/internal/model"
	generatedPaymentV1 "github.com/Igorezka/rocket-factory/shared/pkg/proto/payment/v1"
)

func (s *APISuite) TestPay() {
	type args struct {
		req *generatedPaymentV1.PayOrderRequest
	}

	var (
		transactionUuid = uuid.NewString()
		orderUuid       = uuid.NewString()
		userUuid        = uuid.NewString()
		paymentMethod   = gofakeit.IntRange(0, 4)

		serviceErr = status.Error(codes.Internal, "internal error")

		req = &generatedPaymentV1.PayOrderRequest{
			OrderUuid:     orderUuid,
			UserUuid:      userUuid,
			PaymentMethod: generatedPaymentV1.PaymentMethod(paymentMethod),
		}

		modelPayOrder = model.PayOrder{
			OrderUuid:     orderUuid,
			UserUuid:      userUuid,
			PaymentMethod: model.PaymentMethod(paymentMethod),
		}

		res = &generatedPaymentV1.PayOrderResponse{
			TransactionUuid: transactionUuid,
		}
	)

	tests := []struct {
		name                        string
		args                        args
		want                        *generatedPaymentV1.PayOrderResponse
		err                         error
		paymentServiceMockConfigure func()
	}{
		{
			name: "success case",
			args: args{req: req},
			want: res,
			err:  nil,
			paymentServiceMockConfigure: func() {
				s.paymentService.On("PayOrder", s.ctx, modelPayOrder).Return(transactionUuid, nil).Once()
			},
		},
		{
			name: "error case",
			args: args{req: req},
			want: nil,
			err:  serviceErr,
			paymentServiceMockConfigure: func() {
				s.paymentService.On("PayOrder", s.ctx, modelPayOrder).Return("", serviceErr)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.paymentServiceMockConfigure()
			res, err := s.api.PayOrder(s.ctx, tt.args.req)
			s.Require().Equal(tt.want, res)
			s.Require().Equal(tt.err, err)
		})
	}
}
