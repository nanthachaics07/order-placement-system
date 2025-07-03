package interfaces

import (
	"order-placement-system/internal/domain/entity"
)

type OrderProcessorUseCase interface {
	ProcessOrders(inputOrders []*entity.InputOrder) ([]*entity.CleanedOrder, error)
}

type ComplementaryCalculator interface {
	CalculateWithStartingOrderNo(mainProducts []*entity.Product, startingOrderNo int) ([]*entity.CleanedOrder, error)
}
