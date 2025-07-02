package entity

import (
	"order-placement-system/pkg/errors"
	"order-placement-system/pkg/log"
	"strings"
)

type ParsedProduct struct {
	CleanProductId string  `json:"clean_product_id"`
	Quantity       int     `json:"quantity"`
	OriginalQty    int     `json:"original_qty"`
	UnitPrice      float64 `json:"unit_price"`
	TotalPrice     float64 `json:"total_price"`
}

type Product struct {
	ProductId  string  `json:"product_id"`
	MaterialId string  `json:"material_id"`
	ModelId    string  `json:"model_id"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	TotalPrice float64 `json:"total_price"`
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

// "material-model-texture"
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
