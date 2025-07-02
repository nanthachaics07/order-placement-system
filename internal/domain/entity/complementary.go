package entity

import (
	"order-placement-system/pkg/errors"
	"order-placement-system/pkg/log"
	"strings"
)

const (
	WipingClothProductId = "WIPING-CLOTH"
	CleanerSuffix        = "-CLEANNER"
)

type ComplementaryItem struct {
	ProductId string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type ComplementaryCalculation struct {
	WipingCloth *ComplementaryItem            `json:"wiping_cloth"`
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
			UnitPrice:  0.00,
			TotalPrice: 0.00,
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
				UnitPrice:  0.00,
				TotalPrice: 0.00,
			})
			currentNo++
		}
	}

	return orders
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
