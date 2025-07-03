package service

import "order-placement-system/internal/domain/value_object"

type PriceCalculator interface {
	CalculateUnitPrice(totalPrice *value_object.Price, quantity int) (*value_object.Price, error)
	CalculateTotalPrice(unitPrice *value_object.Price, quantity int) (*value_object.Price, error)
	DividePriceEqually(totalPrice *value_object.Price, parts int) (*value_object.Price, error)
	SumPrices(prices ...*value_object.Price) (*value_object.Price, error)
}
