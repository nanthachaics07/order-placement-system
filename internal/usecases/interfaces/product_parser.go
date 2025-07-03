package interfaces

import (
	"order-placement-system/internal/domain/entity"
)

type ProductParser interface {
	Parse(platformProductId string, originalQty int, totalPrice float64) ([]*entity.ParsedProduct, error)
	CleanPrefix(productId string) string
	ExtractQuantity(productId string) (cleanId string, quantity int, hasQuantity bool)
	SplitBundle(productId string) []string
	ParseProductCode(productId string) (materialId, modelId string, err error)
	Validate(productId string) error
}

type OrderProcessorUseCase interface {
	ProcessOrders(inputOrders []*entity.InputOrder) ([]*entity.CleanedOrder, error)
}

type ComplementaryCalculator interface {
	CalculateWithStartingOrderNo(mainProducts []*entity.Product, startingOrderNo int) ([]*entity.CleanedOrder, error)
}
