package service

import (
	"order-placement-system/internal/domain/entity"
	"order-placement-system/internal/domain/value_object"
)

type ProductParser interface {
	Parse(platformProductId string, originalQty int, totalPrice *value_object.Price) ([]*entity.ParsedProduct, error)
	ParseFromFloat64(platformProductId string, originalQty int, totalPrice float64) ([]*entity.ParsedProduct, error)
	CleanPrefix(productId string) string
	ExtractQuantity(productId string) (cleanId string, quantity int, hasQuantity bool)
	SplitBundle(productId string) []string
	ParseProductCode(productId string) (materialId, modelId string, err error)
	Validate(productId string) error
}
