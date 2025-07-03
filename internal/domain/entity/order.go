package entity

import (
	"order-placement-system/internal/domain/value_object"
	"order-placement-system/pkg/errors"
	"order-placement-system/pkg/log"
)

type InputOrder struct {
	No                int                 `json:"no"`
	PlatformProductId string              `json:"platformProductId"`
	Qty               int                 `json:"qty"`
	UnitPrice         *value_object.Price `json:"unitPrice"`
	TotalPrice        *value_object.Price `json:"totalPrice"`
}

type CleanedOrder struct {
	No         int                 `json:"no"`
	ProductId  string              `json:"productId"`
	MaterialId string              `json:"materialId,omitempty"`
	ModelId    string              `json:"modelId,omitempty"`
	Qty        int                 `json:"qty"`
	UnitPrice  *value_object.Price `json:"unitPrice"`
	TotalPrice *value_object.Price `json:"totalPrice"`
}

type OrderBatch struct {
	Orders []InputOrder
}

func NewOrderBatch(orders []InputOrder) *OrderBatch {
	return &OrderBatch{
		Orders: orders,
	}
}

func (o *InputOrder) IsValid() error {
	if o.No <= 0 {
		log.Errorf("order number must be positive")
		return errors.ErrInvalidInput
	}

	if o.PlatformProductId == "" {
		log.Errorf("platform product id cannot be empty")
		return errors.ErrInvalidInput
	}

	if o.Qty <= 0 {
		log.Errorf("quantity must be positive")
		return errors.ErrInvalidInput
	}

	if o.UnitPrice == nil || o.UnitPrice.Amount() < 0 {
		log.Errorf("unit price cannot be negative")
		return errors.ErrInvalidInput
	}

	if o.TotalPrice == nil || o.TotalPrice.Amount() < 0 {
		log.Errorf("total price cannot be negative")
		return errors.ErrInvalidInput
	}

	return nil
}

func (c *CleanedOrder) IsValid() error {
	if c.No <= 0 {
		log.Errorf("order number must be positive")
		return errors.ErrInvalidInput
	}

	if c.ProductId == "" {
		log.Errorf("product id cannot be empty")
		return errors.ErrInvalidInput
	}

	if c.Qty <= 0 {
		log.Errorf("quantity must be positive")
		return errors.ErrInvalidInput
	}

	if c.UnitPrice == nil || c.UnitPrice.Amount() < 0 {
		log.Errorf("unit price cannot be negative")
		return errors.ErrInvalidInput
	}

	if c.TotalPrice == nil || c.TotalPrice.Amount() < 0 {
		log.Errorf("total price cannot be negative")
		return errors.ErrInvalidInput
	}

	return nil
}
