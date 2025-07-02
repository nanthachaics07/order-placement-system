package entity

import (
	"order-placement-system/pkg/errors"
	"order-placement-system/pkg/log"
)

type InputOrder struct {
	No                int     `json:"no"`
	PlatformProductId string  `json:"platform_product_id"`
	Qty               int     `json:"qty"`
	UnitPrice         float64 `json:"unit_price"`
	TotalPrice        float64 `json:"total_price"`
}

type CleanedOrder struct {
	No         int     `json:"no"`
	ProductId  string  `json:"product_id"`
	MaterialId string  `json:"material_id,omitempty"`
	ModelId    string  `json:"model_id,omitempty"`
	Qty        int     `json:"qty"`
	UnitPrice  float64 `json:"unit_price"`
	TotalPrice float64 `json:"total_price"`
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

	if o.UnitPrice < 0 {
		log.Errorf("unit price cannot be negative")
		return errors.ErrInvalidInput
	}

	if o.TotalPrice < 0 {
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

	if c.UnitPrice < 0 {
		log.Errorf("unit price cannot be negative")
		return errors.ErrInvalidInput
	}

	if c.TotalPrice < 0 {
		log.Errorf("total price cannot be negative")
		return errors.ErrInvalidInput
	}

	return nil
}
