package entity

import (
	"order-placement-system/internal/domain/value_object"
	"order-placement-system/pkg/errors"
	"order-placement-system/pkg/log"
	"strings"
)

const (
	WipingClothProductId = "WIPING-CLOTH"
	CleanerSuffix        = "-CLEANNER"
)

type ComplementaryItem struct {
	ProductId string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

type ComplementaryCalculation struct {
	WipingCloth *ComplementaryItem            `json:"wipingCloth"`
	Cleaners    map[string]*ComplementaryItem `json:"cleaners"`
}

func NewComplementaryCalculation() *ComplementaryCalculation {
	return &ComplementaryCalculation{
		Cleaners: make(map[string]*ComplementaryItem),
	}
}

// adds a product to the complementary calc
func (c *ComplementaryCalculation) AddProduct(product *Product) error {
	if product == nil {
		log.Errorf("product cannot be nil")
		return errors.ErrInvalidInput
	}

	texture := product.GetTexture()
	if texture == "" {
		log.Errorf("product does not have a valid texture", log.S("productId", product.ProductId))
		return errors.ErrInvalidInput
	}

	if !IsValidTexture(texture) {
		log.Errorf("invalid texture", log.S("texture", texture))
		return errors.ErrInvalidInput
	}

	// add Wiping Cloth 1:1
	if c.WipingCloth == nil {
		c.WipingCloth = &ComplementaryItem{
			ProductId: WipingClothProductId,
			Quantity:  0,
		}
	}
	c.WipingCloth.Quantity += product.Quantity

	// add Cleaner based on texture
	cleanerId := generateCleanerId(texture)
	if c.Cleaners[texture] == nil {
		c.Cleaners[texture] = &ComplementaryItem{
			ProductId: cleanerId,
			Quantity:  0,
		}
	}
	c.Cleaners[texture].Quantity += product.Quantity

	return nil
}

// converts the complementary calculation to a list of cleaned orders
func (c *ComplementaryCalculation) ToCleanedOrders(startingNo int) []*CleanedOrder {
	var orders []*CleanedOrder
	currentNo := startingNo

	if c.WipingCloth != nil && c.WipingCloth.Quantity > 0 {
		orders = append(orders, &CleanedOrder{
			No:         currentNo,
			ProductId:  c.WipingCloth.ProductId,
			Qty:        c.WipingCloth.Quantity,
			UnitPrice:  value_object.ZeroPrice(),
			TotalPrice: value_object.ZeroPrice(),
		})
		currentNo++
	}

	// FIXME: Improve memory space usage
	textures := []string{"CLEAR", "MATTE", "PRIVACY"}
	for _, texture := range textures {
		if cleaner, exists := c.Cleaners[texture]; exists && cleaner.Quantity > 0 {
			orders = append(orders, &CleanedOrder{
				No:         currentNo,
				ProductId:  cleaner.ProductId,
				Qty:        cleaner.Quantity,
				UnitPrice:  value_object.ZeroPrice(),
				TotalPrice: value_object.ZeroPrice(),
			})
			currentNo++
		}
	}

	return orders
}

func (c *ComplementaryCalculation) GetTotalComplementaryValue(
	wipingClothPrice *value_object.Price,
	cleanerPrices map[string]*value_object.Price,
) (*value_object.Price, error) {
	totalValue := value_object.ZeroPrice()

	if c.WipingCloth != nil && c.WipingCloth.Quantity > 0 && wipingClothPrice != nil {
		wipingClothValue, err := wipingClothPrice.MultiplyByInt(c.WipingCloth.Quantity)
		if err != nil {
			log.Errorf("failed to calculate wiping cloth value", log.E(err))
			return nil, errors.ErrInvalidInput
		}
		totalValue, err = totalValue.Add(wipingClothValue)
		if err != nil {
			log.Errorf("failed to add wiping cloth value", log.E(err))
			return nil, errors.ErrInvalidInput
		}
	}

	if cleanerPrices != nil {
		for texture, cleaner := range c.Cleaners {
			if cleaner.Quantity > 0 && cleanerPrices[texture] != nil {
				cleanerValue, err := cleanerPrices[texture].MultiplyByInt(cleaner.Quantity)
				if err != nil {
					log.Errorf("failed to calculate %s cleaner value", texture, log.E(err))
					return nil, errors.ErrInvalidInput
				}
				totalValue, err = totalValue.Add(cleanerValue)
				if err != nil {
					log.Errorf("failed to add %s cleaner value", texture, log.E(err))
					return nil, errors.ErrInvalidInput
				}
			}
		}
	}

	return totalValue, nil
}

func generateCleanerId(texture string) string {
	return strings.ToUpper(texture) + CleanerSuffix
}

func IsValidTexture(texture string) bool {
	validTextures := map[string]bool{
		"CLEAR":   true,
		"MATTE":   true,
		"PRIVACY": true,
	}

	return validTextures[strings.ToUpper(texture)]
}

func CalculateComplementaryItems(products []*Product) (*ComplementaryCalculation, error) {
	calc := NewComplementaryCalculation()

	for _, product := range products {
		if err := calc.AddProduct(product); err != nil {
			log.Errorf("Failed to add product", log.E(err), log.S("productId", product.ProductId))
			return nil, errors.ErrInvalidInput
		}
	}

	return calc, nil
}
