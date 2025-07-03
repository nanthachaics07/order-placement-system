package implementation

import (
	"order-placement-system/internal/domain/entity"
	"order-placement-system/internal/usecases/interfaces"
	"order-placement-system/pkg/log"
)

type complementaryCalculatorUseCase struct{}

func NewComplementaryCalculator() interfaces.ComplementaryCalculator {
	return &complementaryCalculatorUseCase{}
}

func (uc *complementaryCalculatorUseCase) CalculateWithStartingOrderNo(mainProducts []*entity.Product, startingOrderNo int) ([]*entity.CleanedOrder, error) {
	if len(mainProducts) == 0 {
		return []*entity.CleanedOrder{}, nil
	}

	calculation := entity.NewComplementaryCalculation()

	for _, product := range mainProducts {
		if err := calculation.AddProduct(product); err != nil {
			log.Errorf("failed to add product to calculation", log.S("product_id", product.ProductId), log.E(err))
			return nil, err
		}
	}

	complementaryOrders := calculation.ToCleanedOrders(startingOrderNo)

	return complementaryOrders, nil
}
