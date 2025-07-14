package order

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/Igorezka/rocket-factory/order/internal/model"
	"github.com/Igorezka/rocket-factory/order/internal/repository/converter"
)

func (s *RepositorySuite) TestGet() {
	type args struct {
		uuid string
	}

	var (
		uuid         string
		expectedPart model.Order
	)

	k := gofakeit.Number(0, len(s.repository.data)-1)
	i := 0
	for _, v := range s.repository.data {
		if i == k {
			uuid = v.Uuid
			expectedPart = converter.OrderToModel(v)
			break
		}
		i++
	}

	tests := []struct {
		name string
		args args
		want model.Order
		err  error
	}{
		{
			name: "success case",
			args: args{
				uuid: uuid,
			},
			want: expectedPart,
			err:  nil,
		},
		{
			name: "error case",
			args: args{
				uuid: gofakeit.UUID(),
			},
			want: model.Order{},
			err:  model.ErrOrderNotFound,
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
