package order

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
	"github.com/stretchr/testify/mock"

	"github.com/Igorezka/rocket-factory/order/internal/model"
)

func (s *ServiceSuite) TestCreate() {
	type args struct {
		orderCreate model.OrderCreate
	}

	var (
		userUuid  = gofakeit.UUID()
		partUuids = []string{gofakeit.UUID(), gofakeit.UUID(), gofakeit.UUID()}

		totalPrice float64
		parts      []model.Part

		orderCreate = model.OrderCreate{
			UserUuid:  userUuid,
			PartUuids: partUuids,
		}

		partsFilter = model.PartsFilter{
			Uuids: partUuids,
		}
	)

	for _, v := range partUuids {
		tmp := model.Part{
			Uuid:          v,
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
		totalPrice += tmp.Price
		parts = append(parts, tmp)
	}

	tests := []struct {
		name           string
		args           args
		want           float64
		err            error
		mocksConfigure func()
	}{
		{
			name: "success case",
			args: args{
				orderCreate: orderCreate,
			},
			want: totalPrice,
			err:  nil,
			mocksConfigure: func() {
				s.inventoryClient.On("ListParts", mock.Anything, partsFilter).Return(parts, nil).Once()
				s.orderRepository.On("Create", s.ctx, mock.Anything).Return(nil).Once()
			},
		},
		{
			name: "parts not found error case",
			args: args{
				orderCreate: orderCreate,
			},
			want: 0,
			err:  model.ErrPartsNotFound,
			mocksConfigure: func() {
				s.inventoryClient.On("ListParts", mock.Anything, partsFilter).Return(nil, model.ErrPartsNotFound).Once()
			},
		},
		{
			name: "part not found error case",
			args: args{
				orderCreate: orderCreate,
			},
			want: 0,
			err:  model.ErrPartNotFound,
			mocksConfigure: func() {
				s.inventoryClient.On("ListParts", mock.Anything, partsFilter).Return(parts[:len(parts)-1], nil).Once()
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mocksConfigure()
			uuid, totalPrice, err := s.service.Create(s.ctx, tt.args.orderCreate)
			if err != nil {
				s.Require().Empty(uuid)
				s.Require().ErrorContains(err, tt.err.Error())
			} else {
				s.Require().NotEmpty(uuid)
				s.Require().Nil(err)
			}
			s.Require().Equal(tt.want, totalPrice)
		})
	}
}
