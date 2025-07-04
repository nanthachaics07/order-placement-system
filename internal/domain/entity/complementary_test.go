package entity_test

import (
	"strings"
	"testing"

	"order-placement-system/internal/domain/entity"
	"order-placement-system/internal/domain/value_object"
	"order-placement-system/pkg/errors"
	"order-placement-system/pkg/log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	log.Init("dev")
}

func TestNewComplementaryCalculation(t *testing.T) {
	calc := entity.NewComplementaryCalculation()

	assert.NotNil(t, calc)
	assert.Nil(t, calc.WipingCloth)
	assert.NotNil(t, calc.Cleaners)
	assert.Empty(t, calc.Cleaners)
}

func TestComplementaryCalculation_AddProduct(t *testing.T) {
	tests := []struct {
		name                string
		products            []*entity.Product
		expectedWipingCloth int
		expectedCleaners    map[string]int
		expectError         bool
	}{
		{
			name: "Add single CLEAR product",
			products: []*entity.Product{
				createValidProduct(t, "FG0A-CLEAR-IPHONE16PROMAX", 2),
			},
			expectedWipingCloth: 2,
			expectedCleaners: map[string]int{
				"CLEAR": 2,
			},
			expectError: false,
		},
		{
			name: "Add single MATTE product",
			products: []*entity.Product{
				createValidProduct(t, "FG05-MATTE-OPPOA3", 3),
			},
			expectedWipingCloth: 3,
			expectedCleaners: map[string]int{
				"MATTE": 3,
			},
			expectError: false,
		},
		{
			name: "Add single PRIVACY product",
			products: []*entity.Product{
				createValidProduct(t, "FG0A-PRIVACY-SAMSUNGS25", 1),
			},
			expectedWipingCloth: 1,
			expectedCleaners: map[string]int{
				"PRIVACY": 1,
			},
			expectError: false,
		},
		{
			name: "Add multiple products with same texture",
			products: []*entity.Product{
				createValidProduct(t, "FG0A-CLEAR-IPHONE16PROMAX", 2),
				createValidProduct(t, "FG0A-CLEAR-OPPOA3", 3),
			},
			expectedWipingCloth: 5,
			expectedCleaners: map[string]int{
				"CLEAR": 5,
			},
			expectError: false,
		},
		{
			name: "Add multiple products with different textures",
			products: []*entity.Product{
				createValidProduct(t, "FG0A-CLEAR-IPHONE16PROMAX", 2),
				createValidProduct(t, "FG05-MATTE-OPPOA3", 3),
				createValidProduct(t, "FG0A-PRIVACY-SAMSUNGS25", 1),
			},
			expectedWipingCloth: 6,
			expectedCleaners: map[string]int{
				"CLEAR":   2,
				"MATTE":   3,
				"PRIVACY": 1,
			},
			expectError: false,
		},
		{
			name: "Add multiple products with mixed textures",
			products: []*entity.Product{
				createValidProduct(t, "FG0A-CLEAR-IPHONE16PROMAX", 2),
				createValidProduct(t, "FG05-MATTE-OPPOA3", 1),
				createValidProduct(t, "FG0A-CLEAR-OPPOA3", 1),
				createValidProduct(t, "FG05-MATTE-IPHONE16PROMAX", 2),
			},
			expectedWipingCloth: 6,
			expectedCleaners: map[string]int{
				"CLEAR": 3,
				"MATTE": 3,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc := entity.NewComplementaryCalculation()

			for _, product := range tt.products {
				err := calc.AddProduct(product)
				if tt.expectError {
					assert.Error(t, err)
					return
				}
				assert.NoError(t, err)
			}

			if !tt.expectError {
				// Check wiping cloth quantity
				assert.NotNil(t, calc.WipingCloth)
				assert.Equal(t, tt.expectedWipingCloth, calc.WipingCloth.Quantity)
				assert.Equal(t, "WIPING-CLOTH", calc.WipingCloth.ProductId)

				// Check cleaners
				assert.Equal(t, len(tt.expectedCleaners), len(calc.Cleaners))
				for texture, expectedQty := range tt.expectedCleaners {
					assert.Contains(t, calc.Cleaners, texture)
					assert.Equal(t, expectedQty, calc.Cleaners[texture].Quantity)
					assert.Equal(t, texture+"-CLEANNER", calc.Cleaners[texture].ProductId)
				}
			}
		})
	}
}

func TestComplementaryCalculation_AddProduct_ErrorCases(t *testing.T) {
	tests := []struct {
		name    string
		product *entity.Product
	}{
		{
			name:    "Nil product",
			product: nil,
		},
		{
			name:    "Product with invalid texture",
			product: createInvalidTextureProduct(t),
		},
		{
			name:    "Product with empty material ID",
			product: createProductWithEmptyMaterialId(t),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc := entity.NewComplementaryCalculation()

			err := calc.AddProduct(tt.product)

			assert.Error(t, err)
			assert.Equal(t, errors.ErrInvalidInput, err)
		})
	}
}

func TestComplementaryCalculation_ToCleanedOrders(t *testing.T) {
	tests := []struct {
		name         string
		products     []*entity.Product
		startingNo   int
		expectedLen  int
		expectedNos  []int
		expectedIds  []string
		expectedQtys []int
	}{
		{
			name: "Single CLEAR product",
			products: []*entity.Product{
				createValidProduct(t, "FG0A-CLEAR-IPHONE16PROMAX", 2),
			},
			startingNo:   3,
			expectedLen:  2,
			expectedNos:  []int{3, 4},
			expectedIds:  []string{"WIPING-CLOTH", "CLEAR-CLEANNER"},
			expectedQtys: []int{2, 2},
		},
		{
			name: "Single MATTE product",
			products: []*entity.Product{
				createValidProduct(t, "FG05-MATTE-OPPOA3", 1),
			},
			startingNo:   1,
			expectedLen:  2,
			expectedNos:  []int{1, 2},
			expectedIds:  []string{"WIPING-CLOTH", "MATTE-CLEANNER"},
			expectedQtys: []int{1, 1},
		},
		{
			name: "Multiple products with different textures",
			products: []*entity.Product{
				createValidProduct(t, "FG0A-CLEAR-IPHONE16PROMAX", 2),
				createValidProduct(t, "FG05-MATTE-OPPOA3", 3),
				createValidProduct(t, "FG0A-PRIVACY-SAMSUNGS25", 1),
			},
			startingNo:   10,
			expectedLen:  4,
			expectedNos:  []int{10, 11, 12, 13},
			expectedIds:  []string{"WIPING-CLOTH", "CLEAR-CLEANNER", "MATTE-CLEANNER", "PRIVACY-CLEANNER"},
			expectedQtys: []int{6, 2, 3, 1},
		},
		{
			name: "Mixed textures with same type",
			products: []*entity.Product{
				createValidProduct(t, "FG0A-CLEAR-IPHONE16PROMAX", 2),
				createValidProduct(t, "FG0A-CLEAR-OPPOA3", 1),
				createValidProduct(t, "FG05-MATTE-OPPOA3", 2),
			},
			startingNo:   5,
			expectedLen:  3,
			expectedNos:  []int{5, 6, 7},
			expectedIds:  []string{"WIPING-CLOTH", "CLEAR-CLEANNER", "MATTE-CLEANNER"},
			expectedQtys: []int{5, 3, 2},
		},
		{
			name:         "Empty products",
			products:     []*entity.Product{},
			startingNo:   1,
			expectedLen:  0,
			expectedNos:  []int{},
			expectedIds:  []string{},
			expectedQtys: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc := entity.NewComplementaryCalculation()

			for _, product := range tt.products {
				err := calc.AddProduct(product)
				require.NoError(t, err)
			}

			orders := calc.ToCleanedOrders(tt.startingNo)

			assert.Equal(t, tt.expectedLen, len(orders))

			for i, order := range orders {
				assert.Equal(t, tt.expectedNos[i], order.No)
				assert.Equal(t, tt.expectedIds[i], order.ProductId)
				assert.Equal(t, tt.expectedQtys[i], order.Qty)
				assert.Equal(t, 0.0, order.UnitPrice.Amount())
				assert.Equal(t, 0.0, order.TotalPrice.Amount())
				assert.Empty(t, order.MaterialId)
				assert.Empty(t, order.ModelId)
			}
		})
	}
}

func TestComplementaryCalculation_GetTotalComplementaryValue(t *testing.T) {
	tests := []struct {
		name             string
		products         []*entity.Product
		wipingClothPrice float64
		cleanerPrices    map[string]float64
		expectedTotal    float64
		expectError      bool
	}{
		{
			name: "Zero prices (default behavior)",
			products: []*entity.Product{
				createValidProduct(t, "FG0A-CLEAR-IPHONE16PROMAX", 2),
			},
			wipingClothPrice: 0.0,
			cleanerPrices:    map[string]float64{"CLEAR": 0.0},
			expectedTotal:    0.0,
			expectError:      false,
		},
		{
			name: "Non-zero prices",
			products: []*entity.Product{
				createValidProduct(t, "FG0A-CLEAR-IPHONE16PROMAX", 2),
				createValidProduct(t, "FG05-MATTE-OPPOA3", 1),
			},
			wipingClothPrice: 1.0,
			cleanerPrices:    map[string]float64{"CLEAR": 2.0, "MATTE": 3.0},
			expectedTotal:    10.0, // (2+1)*1.0 + 2*2.0 + 1*3.0 = 3 + 4 + 3 = 10
			expectError:      false,
		},
		{
			name: "Missing cleaner prices",
			products: []*entity.Product{
				createValidProduct(t, "FG0A-CLEAR-IPHONE16PROMAX", 2),
			},
			wipingClothPrice: 1.0,
			cleanerPrices:    nil,
			expectedTotal:    2.0, // Only wiping cloth: 2*1.0 = 2
			expectError:      false,
		},
		{
			name: "Partial cleaner prices",
			products: []*entity.Product{
				createValidProduct(t, "FG0A-CLEAR-IPHONE16PROMAX", 2),
				createValidProduct(t, "FG05-MATTE-OPPOA3", 1),
			},
			wipingClothPrice: 1.0,
			cleanerPrices:    map[string]float64{"CLEAR": 2.0}, // Missing MATTE price
			expectedTotal:    7.0,                              // 3*1.0 + 2*2.0 = 3 + 4 = 7
			expectError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc := entity.NewComplementaryCalculation()

			for _, product := range tt.products {
				err := calc.AddProduct(product)
				require.NoError(t, err)
			}

			var wipingClothPrice *value_object.Price
			if tt.wipingClothPrice >= 0 {
				wipingClothPrice = value_object.MustNewPrice(tt.wipingClothPrice)
			}

			var cleanerPrices map[string]*value_object.Price
			if tt.cleanerPrices != nil {
				cleanerPrices = make(map[string]*value_object.Price)
				for texture, price := range tt.cleanerPrices {
					cleanerPrices[texture] = value_object.MustNewPrice(price)
				}
			}

			total, err := calc.GetTotalComplementaryValue(wipingClothPrice, cleanerPrices)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, total)
				assert.Equal(t, errors.ErrInvalidInput, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, total)
				assert.Equal(t, tt.expectedTotal, total.Amount())
			}
		})
	}
}

func TestGenerateCleanerId(t *testing.T) {
	tests := []struct {
		name     string
		texture  string
		expected string
	}{
		{
			name:     "CLEAR texture",
			texture:  "clear",
			expected: "CLEAR-CLEANNER",
		},
		{
			name:     "MATTE texture",
			texture:  "matte",
			expected: "MATTE-CLEANNER",
		},
		{
			name:     "PRIVACY texture",
			texture:  "privacy",
			expected: "PRIVACY-CLEANNER",
		},
		{
			name:     "Already uppercase",
			texture:  "CLEAR",
			expected: "CLEAR-CLEANNER",
		},
		{
			name:     "Mixed case",
			texture:  "MaTtE",
			expected: "MATTE-CLEANNER",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc := entity.NewComplementaryCalculation()
			product := createValidProductWithTexture(t, tt.texture, 1)

			err := calc.AddProduct(product)
			require.NoError(t, err)

			upperTexture := strings.ToUpper(tt.texture)
			assert.Contains(t, calc.Cleaners, upperTexture)
			assert.Equal(t, tt.expected, calc.Cleaners[upperTexture].ProductId)
		})
	}
}

func TestIsValidTexture(t *testing.T) {
	tests := []struct {
		name     string
		texture  string
		expected bool
	}{
		{
			name:     "Valid CLEAR texture",
			texture:  "CLEAR",
			expected: true,
		},
		{
			name:     "Valid MATTE texture",
			texture:  "MATTE",
			expected: true,
		},
		{
			name:     "Valid PRIVACY texture",
			texture:  "PRIVACY",
			expected: true,
		},
		{
			name:     "Valid lowercase clear",
			texture:  "clear",
			expected: true,
		},
		{
			name:     "Valid lowercase matte",
			texture:  "matte",
			expected: true,
		},
		{
			name:     "Valid lowercase privacy",
			texture:  "privacy",
			expected: true,
		},
		{
			name:     "Valid mixed case",
			texture:  "MaTtE",
			expected: true,
		},
		{
			name:     "Invalid texture",
			texture:  "INVALID",
			expected: false,
		},
		{
			name:     "Empty string",
			texture:  "",
			expected: false,
		},
		{
			name:     "Whitespace only",
			texture:  "   ",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := entity.IsValidTexture(tt.texture)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateComplementaryItems(t *testing.T) {
	tests := []struct {
		name                string
		products            []*entity.Product
		expectedWipingCloth int
		expectedCleaners    map[string]int
		expectError         bool
	}{
		{
			name: "Single product",
			products: []*entity.Product{
				createValidProduct(t, "FG0A-CLEAR-IPHONE16PROMAX", 2),
			},
			expectedWipingCloth: 2,
			expectedCleaners: map[string]int{
				"CLEAR": 2,
			},
			expectError: false,
		},
		{
			name: "Multiple products",
			products: []*entity.Product{
				createValidProduct(t, "FG0A-CLEAR-IPHONE16PROMAX", 2),
				createValidProduct(t, "FG05-MATTE-OPPOA3", 1),
				createValidProduct(t, "FG0A-PRIVACY-SAMSUNGS25", 3),
			},
			expectedWipingCloth: 6,
			expectedCleaners: map[string]int{
				"CLEAR":   2,
				"MATTE":   1,
				"PRIVACY": 3,
			},
			expectError: false,
		},
		{
			name:                "Empty products",
			products:            []*entity.Product{},
			expectedWipingCloth: 0,
			expectedCleaners:    map[string]int{},
			expectError:         false,
		},
		{
			name: "Product with invalid texture",
			products: []*entity.Product{
				createInvalidTextureProduct(t),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc, err := entity.CalculateComplementaryItems(tt.products)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, calc)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, calc)

				if tt.expectedWipingCloth > 0 {
					assert.NotNil(t, calc.WipingCloth)
					assert.Equal(t, tt.expectedWipingCloth, calc.WipingCloth.Quantity)
				} else {
					assert.Nil(t, calc.WipingCloth)
				}

				assert.Equal(t, len(tt.expectedCleaners), len(calc.Cleaners))
				for texture, expectedQty := range tt.expectedCleaners {
					assert.Contains(t, calc.Cleaners, texture)
					assert.Equal(t, expectedQty, calc.Cleaners[texture].Quantity)
				}
			}
		})
	}
}

func TestComplementaryCalculation_EdgeCases(t *testing.T) {
	t.Run("Very large quantities", func(t *testing.T) {
		calc := entity.NewComplementaryCalculation()
		product := createValidProduct(t, "FG0A-CLEAR-IPHONE16PROMAX", 1000000)

		err := calc.AddProduct(product)
		assert.NoError(t, err)

		assert.Equal(t, 1000000, calc.WipingCloth.Quantity)
		assert.Equal(t, 1000000, calc.Cleaners["CLEAR"].Quantity)
	})

	t.Run("Zero quantity product", func(t *testing.T) {
		calc := entity.NewComplementaryCalculation()
		product := createValidProduct(t, "FG0A-CLEAR-IPHONE16PROMAX", 0)

		err := calc.AddProduct(product)
		assert.NoError(t, err)

		assert.Equal(t, 0, calc.WipingCloth.Quantity)
		assert.Equal(t, 0, calc.Cleaners["CLEAR"].Quantity)
	})

	t.Run("Sequential additions", func(t *testing.T) {
		calc := entity.NewComplementaryCalculation()

		// Add first product
		product1 := createValidProduct(t, "FG0A-CLEAR-IPHONE16PROMAX", 1)
		err := calc.AddProduct(product1)
		assert.NoError(t, err)
		assert.Equal(t, 1, calc.WipingCloth.Quantity)
		assert.Equal(t, 1, calc.Cleaners["CLEAR"].Quantity)

		// Add second product with same texture
		product2 := createValidProduct(t, "FG0A-CLEAR-OPPOA3", 2)
		err = calc.AddProduct(product2)
		assert.NoError(t, err)
		assert.Equal(t, 3, calc.WipingCloth.Quantity)
		assert.Equal(t, 3, calc.Cleaners["CLEAR"].Quantity)

		// Add third product with different texture
		product3 := createValidProduct(t, "FG05-MATTE-OPPOA3", 1)
		err = calc.AddProduct(product3)
		assert.NoError(t, err)
		assert.Equal(t, 4, calc.WipingCloth.Quantity)
		assert.Equal(t, 3, calc.Cleaners["CLEAR"].Quantity)
		assert.Equal(t, 1, calc.Cleaners["MATTE"].Quantity)
	})
}

// func createValidProduct(t *testing.T, productId string, quantity int) *entity.Product {
// 	t.Helper()

// 	unitPrice := value_object.MustNewPrice(50.0)
// 	totalPrice := value_object.MustNewPrice(float64(quantity) * 50.0)

// 	product, err := entity.NewProduct(productId, quantity, unitPrice, totalPrice)
// 	require.NoError(t, err)

// 	return product
// }

func createValidProductWithTexture(t *testing.T, texture string, quantity int) *entity.Product {
	t.Helper()

	productId := "FG0A-" + strings.ToUpper(texture) + "-IPHONE16PROMAX"
	return createValidProduct(t, productId, quantity)
}

func createInvalidTextureProduct(t *testing.T) *entity.Product {
	t.Helper()

	// Create a product with manual material ID to bypass validation
	unitPrice := value_object.MustNewPrice(50.0)
	totalPrice := value_object.MustNewPrice(50.0)

	return &entity.Product{
		ProductId:  "FG0A-INVALID-IPHONE16PROMAX",
		MaterialId: "FG0A-INVALID",
		ModelId:    "IPHONE16PROMAX",
		Quantity:   1,
		UnitPrice:  unitPrice,
		TotalPrice: totalPrice,
	}
}

func createProductWithEmptyMaterialId(t *testing.T) *entity.Product {
	t.Helper()

	unitPrice := value_object.MustNewPrice(50.0)
	totalPrice := value_object.MustNewPrice(50.0)

	return &entity.Product{
		ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
		MaterialId: "",
		ModelId:    "IPHONE16PROMAX",
		Quantity:   1,
		UnitPrice:  unitPrice,
		TotalPrice: totalPrice,
	}
}

func BenchmarkComplementaryCalculation_ToCleanedOrders(b *testing.B) {
	calc := entity.NewComplementaryCalculation()
	products := []*entity.Product{
		createValidProduct(b, "FG0A-CLEAR-IPHONE16PROMAX", 2),
		createValidProduct(b, "FG05-MATTE-OPPOA3", 1),
		createValidProduct(b, "FG0A-PRIVACY-SAMSUNGS25", 3),
	}

	for _, product := range products {
		_ = calc.AddProduct(product)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = calc.ToCleanedOrders(1)
	}
}

func BenchmarkCalculateComplementaryItems(b *testing.B) {
	products := []*entity.Product{
		createValidProduct(b, "FG0A-CLEAR-IPHONE16PROMAX", 2),
		createValidProduct(b, "FG05-MATTE-OPPOA3", 1),
		createValidProduct(b, "FG0A-PRIVACY-SAMSUNGS25", 3),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = entity.CalculateComplementaryItems(products)
	}
}

// Helper function for benchmark tests
func createValidProduct(tb testing.TB, productId string, quantity int) *entity.Product {
	tb.Helper()

	unitPrice := value_object.MustNewPrice(50.0)
	totalPrice := value_object.MustNewPrice(float64(quantity) * 50.0)

	product, err := entity.NewProduct(productId, quantity, unitPrice, totalPrice)
	if err != nil {
		tb.Fatal(err)
	}

	return product
}
