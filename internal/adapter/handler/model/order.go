package model

import (
	"order-placement-system/internal/domain/entity"
	"order-placement-system/internal/domain/value_object"
	"order-placement-system/pkg/errors"
	"order-placement-system/pkg/log"

	"github.com/gin-gonic/gin"
)

type InputOrder struct {
	No                int     `json:"no" binding:"required,min=1"`
	PlatformProductId string  `json:"platformProductId" binding:"required"`
	Qty               int     `json:"qty" binding:"required,min=1"`
	UnitPrice         float64 `json:"unitPrice" binding:"required,min=0"`
	TotalPrice        float64 `json:"totalPrice" binding:"required,min=0"`
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

func (o *InputOrder) Parse(c *gin.Context) ([]*InputOrder, error) {
	var orders []*InputOrder

	if err := c.ShouldBindJSON(&orders); err != nil {
		log.Errorf("failed to bind JSON", log.E(err))
		return nil, errors.ErrInvalidInput
	}

	if len(orders) == 0 {
		log.Error("empty orders array")
		return nil, errors.ErrInvalidInput
	}

	return orders, nil
}

func (o *InputOrder) ToEntity() (*entity.InputOrder, error) {
	unitPrice, err := value_object.NewPrice(o.UnitPrice)
	if err != nil {
		return nil, errors.ErrInvalidInput
	}

	totalPrice, err := value_object.NewPrice(o.TotalPrice)
	if err != nil {
		return nil, errors.ErrInvalidInput
	}

	return &entity.InputOrder{
		No:                o.No,
		PlatformProductId: o.PlatformProductId,
		Qty:               o.Qty,
		UnitPrice:         unitPrice,
		TotalPrice:        totalPrice,
	}, nil
}

func ToEntity(models []*InputOrder) ([]*entity.InputOrder, error) {
	entities := make([]*entity.InputOrder, len(models))
	for i, model := range models {
		entity, err := model.ToEntity()
		if err != nil {
			return nil, err
		}
		entities[i] = entity
	}
	return entities, nil
}

func FromEntity(e *entity.CleanedOrder) *CleanedOrder {
	return &CleanedOrder{
		No:         e.No,
		ProductId:  e.ProductId,
		MaterialId: e.MaterialId,
		ModelId:    e.ModelId,
		Qty:        e.Qty,
		UnitPrice:  e.UnitPrice,
		TotalPrice: e.TotalPrice,
	}
}

func FromEntities(entities []*entity.CleanedOrder) []*CleanedOrder {
	models := make([]*CleanedOrder, len(entities))
	for i, e := range entities {
		models[i] = FromEntity(e)
	}
	return models
}

func (o *InputOrder) Validate() error {
	if o.No <= 0 {
		return errors.ErrInvalidInput
	}

	if o.PlatformProductId == "" {
		return errors.ErrInvalidInput
	}

	if o.Qty <= 0 {
		return errors.ErrInvalidInput
	}

	if o.UnitPrice < 0 {
		return errors.ErrInvalidInput
	}

	if o.TotalPrice < 0 {
		return errors.ErrInvalidInput
	}

	return nil
}
