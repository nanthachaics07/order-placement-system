package entity

// TODO: Add Validation
type InputOrder struct {
	No                int     `json:"no"`
	PlatformProductId string  `json:"platformProductId"`
	Qty               int     `json:"qty"`
	UnitPrice         float64 `json:"unitPrice"`
	TotalPrice        float64 `json:"totalPrice"`
}

type CleanedOrder struct {
	No         int     `json:"no"`
	ProductId  string  `json:"productId"`
	MaterialId string  `json:"materialId,omitempty"`
	ModelId    string  `json:"modelId,omitempty"`
	Qty        int     `json:"qty"`
	UnitPrice  float64 `json:"unitPrice"`
	TotalPrice float64 `json:"totalPrice"`
}

type OrderBatch struct {
	Orders []InputOrder
}

func NewOrderBatch(orders []InputOrder) *OrderBatch {
	return &OrderBatch{
		Orders: orders,
	}
}
