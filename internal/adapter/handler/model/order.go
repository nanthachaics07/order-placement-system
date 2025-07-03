package model

import (
	"order-placement-system/internal/domain/entity"
	"order-placement-system/pkg/errors"
	"order-placement-system/pkg/log"

	"github.com/gin-gonic/gin"
)

type InputOrder struct {
	No                int     `json:"no"`
	PlatformProductId string  `json:"platformProductId"`
	Qty               int     `json:"qty"`
	UnitPrice         float64 `json:"unitPrice"`
	TotalPrice        float64 `json:"totalPrice"`
}

func (r *InputOrder) Parse(c *gin.Context) ([]*InputOrder, error) {
	var inputOrders []*InputOrder

	if err := c.ShouldBindJSON(&inputOrders); err != nil {
		log.Errorf("failed to bind JSON", log.E(err))
		return nil, errors.ErrInvalidInput
	}

	for i, order := range inputOrders {
		if err := r.validateInputOrder(order, i); err != nil {
			log.Errorf("validation failed for order",
				log.S("index", string(rune(i))), log.E(err))
			return nil, err
		}
	}

	log.Infof("successfully parsed input orders",
		log.S("count", string(rune(len(inputOrders)))))

	return inputOrders, nil
}

func (r *InputOrder) validateInputOrder(order *InputOrder, index int) error {
	if order == nil {
		return errors.ErrInvalidInput
	}
	if order.PlatformProductId == "" {
		return errors.ErrInvalidInput
	}
	if order.Qty <= 0 {
		return errors.ErrInvalidInput
	}
	if order.UnitPrice < 0 {
		return errors.ErrInvalidInput
	}
	if order.TotalPrice < 0 {
		return errors.ErrInvalidInput
	}
	return nil
}

func ToEntity(models []*InputOrder) []*entity.InputOrder {
	entities := make([]*entity.InputOrder, 0, len(models))

	for _, model := range models {
		entity := &entity.InputOrder{
			No:                model.No,
			PlatformProductId: model.PlatformProductId,
			Qty:               model.Qty,
			UnitPrice:         model.UnitPrice,
			TotalPrice:        model.TotalPrice,
		}
		entities = append(entities, entity)
	}

	return entities
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

func (r *CleanedOrder) fromEntity(entity *entity.CleanedOrder) *CleanedOrder {
	return &CleanedOrder{
		No:         entity.No,
		ProductId:  entity.ProductId,
		MaterialId: entity.MaterialId,
		ModelId:    entity.ModelId,
		Qty:        entity.Qty,
		UnitPrice:  entity.UnitPrice,
		TotalPrice: entity.TotalPrice,
	}
}

func FromEntities(entities []*entity.CleanedOrder) []*CleanedOrder {
	models := make([]*CleanedOrder, 0, len(entities))

	for _, entity := range entities {
		model := &CleanedOrder{}
		models = append(models, model.fromEntity(entity))
	}

	return models
}
