package implementation

import (
	"strconv"

	"order-placement-system/internal/domain/entity"
	"order-placement-system/internal/domain/service"
	usecase "order-placement-system/internal/usecases/interfaces"
	"order-placement-system/pkg/errors"
	"order-placement-system/pkg/log"
)

type orderProcessorUseCase struct {
	productParser           service.ProductParser
	complementaryCalculator usecase.ComplementaryCalculator
}

func NewOrderProcessor(
	parser service.ProductParser,
	complementaryCalculator usecase.ComplementaryCalculator,
) usecase.OrderProcessorUseCase {
	return &orderProcessorUseCase{
		productParser:           parser,
		complementaryCalculator: complementaryCalculator,
	}
}

func (uc *orderProcessorUseCase) ProcessOrders(inputOrders []*entity.InputOrder) ([]*entity.CleanedOrder, error) {
	if len(inputOrders) == 0 {
		return []*entity.CleanedOrder{}, nil
	}

	if err := uc.validateInputOrders(inputOrders); err != nil {
		log.Errorf("invalid input orders", log.E(err))
		return nil, err
	}

	var allMainProducts []*entity.Product
	var allCleanedOrders []*entity.CleanedOrder
	currentOrderNo := 1

	// Process each input order
	for _, inputOrder := range inputOrders {
		parsedProducts, err := uc.productParser.Parse(
			inputOrder.PlatformProductId,
			inputOrder.Qty,
			inputOrder.TotalPrice,
		)
		if err != nil {
			log.Errorf("failed to parse product id", log.S("product_id", inputOrder.PlatformProductId), log.E(err))
			return nil, err
		}

		for _, parsedProduct := range parsedProducts {
			product, err := uc.createProductFromParsed(parsedProduct)
			if err != nil {
				log.Errorf("failed to create product from parsed data", log.S("product_id", parsedProduct.CleanProductId), log.E(err))
				return nil, err
			}

			allMainProducts = append(allMainProducts, product)

			cleanedOrder := product.ToCleanedOrder(currentOrderNo)
			allCleanedOrders = append(allCleanedOrders, cleanedOrder)

			currentOrderNo++
		}
	}

	complementaryOrders, err := uc.complementaryCalculator.CalculateWithStartingOrderNo(allMainProducts, currentOrderNo)
	if err != nil {
		log.Errorf("failed to calculate complementary items", log.E(err))
		return nil, err
	}

	allCleanedOrders = append(allCleanedOrders, complementaryOrders...)

	if err := uc.validateCleanedOrders(allCleanedOrders); err != nil {
		log.Errorf("invalid cleaned orders", log.E(err))
		return nil, err
	}

	return allCleanedOrders, nil
}

func (uc *orderProcessorUseCase) createProductFromParsed(parsedProduct *entity.ParsedProduct) (*entity.Product, error) {
	materialId, modelId, err := uc.productParser.ParseProductCode(parsedProduct.CleanProductId)
	if err != nil {
		log.Errorf("failed to parse product code", log.S("product_code", parsedProduct.CleanProductId), log.E(err))
		return nil, err
	}

	product := &entity.Product{
		ProductId:  parsedProduct.CleanProductId,
		MaterialId: materialId,
		ModelId:    modelId,
		Quantity:   parsedProduct.Quantity,
		UnitPrice:  parsedProduct.UnitPrice,
		TotalPrice: parsedProduct.TotalPrice,
	}

	if err := product.IsValid(); err != nil {
		log.Errorf("invalid product", log.S("product_id", product.ProductId), log.E(err))
		return nil, err
	}

	return product, nil
}

func (uc *orderProcessorUseCase) validateInputOrders(inputOrders []*entity.InputOrder) error {
	for i, order := range inputOrders {
		if order == nil {
			log.Errorf("input order at index is nil", log.S("index", strconv.Itoa(i)))
			return errors.ErrInvalidInput
		}

		if err := order.IsValid(); err != nil {
			log.Errorf("input order is invalid", log.S("order_no", strconv.Itoa(order.No)), log.E(err))
			return err
		}
	}

	return nil
}

func (uc *orderProcessorUseCase) validateCleanedOrders(cleanedOrders []*entity.CleanedOrder) error {
	for i, order := range cleanedOrders {
		if order == nil {
			log.Errorf("cleaned order at index is nil", log.S("index", strconv.Itoa(i)))
			return errors.ErrInvalidInput
		}

		if err := order.IsValid(); err != nil {
			log.Errorf("cleaned order is invalid", log.S("order_no", strconv.Itoa(order.No)), log.E(err))
			return err
		}
	}

	return nil
}
