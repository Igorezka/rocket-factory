package order

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/Igorezka/rocket-factory/order/internal/model"
)

func (s *RepositorySuite) TestUpdate() {
	type args struct {
		uuid        string
		orderUpdate model.OrderUpdate
	}

	var (
		uuid            string
		partUuids       = []string{gofakeit.UUID(), gofakeit.UUID(), gofakeit.UUID()}
		totalPrice      = gofakeit.Float64()
		transactionUuid = gofakeit.UUID()
		paymentMethod   = model.PaymentMethodCard
		status          = model.OrderStatusPaid

		orderUpdate = model.OrderUpdate{
			PartUuids:       &partUuids,
			TotalPrice:      &totalPrice,
			TransactionUuid: &transactionUuid,
			PaymentMethod:   &paymentMethod,
			Status:          &status,
		}
	)

	k := gofakeit.Number(0, len(s.repository.data)-1)
	i := 0
	for _, v := range s.repository.data {
		if i == k {
			uuid = v.Uuid
			break
		}
		i++
	}

	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "success case",
			args: args{
				uuid:        uuid,
				orderUpdate: orderUpdate,
			},
			err: nil,
		},
		{
			name: "error case",
			args: args{
				uuid:        gofakeit.UUID(),
				orderUpdate: orderUpdate,
			},
			err: model.ErrOrderNotFound,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			err := s.repository.Update(s.ctx, tt.args.uuid, tt.args.orderUpdate)
			s.Require().Equal(tt.err, err)
		})
	}
}
