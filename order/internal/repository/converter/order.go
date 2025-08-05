package converter

import (
	"github.com/samber/lo"

	"github.com/Igorezka/rocket-factory/order/internal/model"
	repoModel "github.com/Igorezka/rocket-factory/order/internal/repository/model"
)

func OrderToModel(order repoModel.Order) model.Order {
	var method *model.PaymentMethod
	if order.PaymentMethod != nil {
		method = lo.ToPtr(model.PaymentMethod(*order.PaymentMethod))
	}
	return model.Order{
		Uuid:            order.Uuid,
		UserUuid:        order.UserUuid,
		PartUuids:       order.PartUuids,
		TransactionUuid: order.TransactionUuid,
		PaymentMethod:   method,
		Status:          model.Status(order.Status),
	}
}

func OrderToRepoModel(order model.Order) repoModel.Order {
	var method *repoModel.PaymentMethod
	if order.PaymentMethod != nil {
		method = lo.ToPtr(repoModel.PaymentMethod(*order.PaymentMethod))
	}

	return repoModel.Order{
		Uuid:            order.Uuid,
		UserUuid:        order.UserUuid,
		PartUuids:       order.PartUuids,
		TransactionUuid: order.TransactionUuid,
		PaymentMethod:   method,
		Status:          repoModel.Status(order.Status),
	}
}
