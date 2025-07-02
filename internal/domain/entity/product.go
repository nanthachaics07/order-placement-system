package entity

import (
	"order-placement-system/pkg/errors"
	"order-placement-system/pkg/log"
	"strings"
)

type ParsedProduct struct {
	CleanProductId string  `json:"cleanProductId"`
	Quantity       int     `json:"quantity"`
	OriginalQty    int     `json:"originalQty"`
	UnitPrice      float64 `json:"unitPrice"`
	TotalPrice     float64 `json:"totalPrice"`
}

type Product struct {
	ProductId  string  `json:"productId"`
	MaterialId string  `json:"materialId"`
	ModelId    string  `json:"modelId"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unitPrice"`
	TotalPrice float64 `json:"totalPrice"`
}

func NewProduct(productId string, quantity int, unitPrice, totalPrice float64) (*Product, error) {
	materialId, modelId, err := parseProductCode(productId)
	if err != nil {
		log.Errorf("Failed to parse product code", log.E(err), productId)
		return nil, errors.ErrInvalidInput
	}

	return &Product{
		ProductId:  productId,
		MaterialId: materialId,
		ModelId:    modelId,
		Quantity:   quantity,
		UnitPrice:  unitPrice,
		TotalPrice: totalPrice,
	}, nil
}

// material-model-texture
func parseProductCode(productId string) (materialId, modelId string, err error) {
	if productId == "" {
		log.Error("Product ID is empty")
		return "", "", errors.ErrInvalidInput
	}

	parts := strings.Split(productId, "-")
	if len(parts) < 3 {
		log.Errorf("Invalid product ID format", log.S("productId", productId))
		return "", "", errors.ErrInvalidInput
	}

	materialId = strings.Join(parts[:2], "-")
	modelId = strings.Join(parts[2:], "-")

	return materialId, modelId, nil
}

func (p *Product) GetTexture() string {
	parts := strings.Split(p.MaterialId, "-")
	if len(parts) >= 2 {
		return parts[1]
	}
	return ""
}

func (p *Product) ToCleanedOrder(orderNo int) *CleanedOrder {
	return &CleanedOrder{
		No:         orderNo,
		ProductId:  p.ProductId,
		MaterialId: p.MaterialId,
		ModelId:    p.ModelId,
		Qty:        p.Quantity,
		UnitPrice:  p.UnitPrice,
		TotalPrice: p.TotalPrice,
	}
}

func (p *Product) IsValid() error {
	if p.ProductId == "" {
		log.Error("Product ID cannot be empty")
		return errors.ErrInvalidInput
	}

	if p.MaterialId == "" {
		log.Error("Material ID cannot be empty")
		return errors.ErrInvalidInput
	}

	if p.ModelId == "" {
		log.Error("Model ID cannot be empty")
		return errors.ErrInvalidInput
	}

	if p.Quantity <= 0 {
		log.Error("Quantity must be positive")
		return errors.ErrInvalidInput
	}

	if p.UnitPrice < 0 {
		log.Error("Unit price cannot be negative")
		return errors.ErrInvalidInput
	}

	if p.TotalPrice < 0 {
		log.Error("Total price cannot be negative")
		return errors.ErrInvalidInput
	}

	return nil
}

func (c *CleanedOrder) IsMainProduct() bool {
	return c.MaterialId != "" && c.ModelId != ""
}

func (c *CleanedOrder) IsComplementaryProduct() bool {
	return !c.IsMainProduct()
}
