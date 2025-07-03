package parser_test

import (
	"order-placement-system/internal/domain/value_object"
	"order-placement-system/pkg/log"
	"order-placement-system/pkg/utils/parser"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	log.Init("dev")
}

func TestProductParser_CleanPrefix(t *testing.T) {

	parser := parser.NewProductParser()

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Clean x2-3& prefix",
			input:    "x2-3&FG0A-CLEAR-IPHONE16PROMAX",
			expected: "FG0A-CLEAR-IPHONE16PROMAX",
		},
		{
			name:     "Clean -- prefix",
			input:    "--FG0A-CLEAR-IPHONE16PROMAX",
			expected: "FG0A-CLEAR-IPHONE16PROMAX",
		},
		{
			name:     "Clean %20- prefix",
			input:    "%20-FG0A-CLEAR-IPHONE16PROMAX",
			expected: "FG0A-CLEAR-IPHONE16PROMAX",
		},
		{
			name:     "Clean %20x prefix",
			input:    "%20xFG0A-CLEAR-IPHONE16PROMAX",
			expected: "FG0A-CLEAR-IPHONE16PROMAX",
		},
		{
			name:     "No prefix to clean",
			input:    "FG0A-CLEAR-IPHONE16PROMAX",
			expected: "FG0A-CLEAR-IPHONE16PROMAX",
		},
		{
			name:     "Nested prefixes (-- then x2-3&)",
			input:    "--x2-3&FG0A-CLEAR-IPHONE16PROMAX",
			expected: "FG0A-CLEAR-IPHONE16PROMAX",
		},
		{
			name:     "Multiple same prefixes",
			input:    "----FG0A-CLEAR-IPHONE16PROMAX",
			expected: "FG0A-CLEAR-IPHONE16PROMAX",
		},
		{
			name:     "Complex nested prefixes",
			input:    "%20--%20x--x2-3&FG0A-CLEAR-IPHONE16PROMAX",
			expected: "FG0A-CLEAR-IPHONE16PROMAX",
		},
		{
			name:     "Mixed order prefixes",
			input:    "x2-3&--FG0A-CLEAR-IPHONE16PROMAX",
			expected: "FG0A-CLEAR-IPHONE16PROMAX",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := parser.CleanPrefix(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestProductParser_ExtractQuantity(t *testing.T) {

	parser := parser.NewProductParser()

	testCases := []struct {
		name        string
		input       string
		cleanId     string
		quantity    int
		hasQuantity bool
	}{
		{
			name:        "Extract *2",
			input:       "FG0A-CLEAR-IPHONE16PROMAX*2",
			cleanId:     "FG0A-CLEAR-IPHONE16PROMAX",
			quantity:    2,
			hasQuantity: true,
		},
		{
			name:        "Extract *3",
			input:       "FG0A-MATTE-IPHONE16PROMAX*3",
			cleanId:     "FG0A-MATTE-IPHONE16PROMAX",
			quantity:    3,
			hasQuantity: true,
		},
		{
			name:        "Extract *10",
			input:       "FG0A-PRIVACY-SAMSUNGS25*10",
			cleanId:     "FG0A-PRIVACY-SAMSUNGS25",
			quantity:    10,
			hasQuantity: true,
		},
		{
			name:        "No quantity indicator",
			input:       "FG0A-CLEAR-IPHONE16PROMAX",
			cleanId:     "FG0A-CLEAR-IPHONE16PROMAX",
			quantity:    1,
			hasQuantity: false,
		},
		{
			name:        "Empty string",
			input:       "",
			cleanId:     "",
			quantity:    1,
			hasQuantity: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cleanId, quantity, hasQuantity := parser.ExtractQuantity(tc.input)
			assert.Equal(t, tc.cleanId, cleanId)
			assert.Equal(t, tc.quantity, quantity)
			assert.Equal(t, tc.hasQuantity, hasQuantity)
		})
	}
}

func TestProductParser_SplitBundle(t *testing.T) {

	parser := parser.NewProductParser()

	testCases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:  "Split two products",
			input: "FG0A-CLEAR-OPPOA3/%20xFG0A-CLEAR-OPPOA3-B",
			expected: []string{
				"FG0A-CLEAR-OPPOA3",
				"FG0A-CLEAR-OPPOA3-B",
			},
		},
		{
			name:  "Split three products",
			input: "FG0A-CLEAR-OPPOA3/%20xFG0A-CLEAR-OPPOA3-B/FG0A-MATTE-OPPOA3",
			expected: []string{
				"FG0A-CLEAR-OPPOA3",
				"FG0A-CLEAR-OPPOA3-B",
				"FG0A-MATTE-OPPOA3",
			},
		},
		{
			name:  "Single product (no bundle)",
			input: "FG0A-CLEAR-IPHONE16PROMAX",
			expected: []string{
				"FG0A-CLEAR-IPHONE16PROMAX",
			},
		},
		{
			name:     "Empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:  "Bundle with spaces",
			input: "FG0A-CLEAR-OPPOA3 / FG0A-MATTE-OPPOA3",
			expected: []string{
				"FG0A-CLEAR-OPPOA3",
				"FG0A-MATTE-OPPOA3",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := parser.SplitBundle(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestProductParser_ParseProductCode(t *testing.T) {

	parser := parser.NewProductParser()

	testCases := []struct {
		name       string
		input      string
		materialId string
		modelId    string
		expectErr  bool
	}{
		{
			name:       "Parse standard product code",
			input:      "FG0A-CLEAR-IPHONE16PROMAX",
			materialId: "FG0A-CLEAR",
			modelId:    "IPHONE16PROMAX",
			expectErr:  false,
		},
		{
			name:       "Parse product with multiple dashes in model",
			input:      "FG0A-MATTE-OPPOA3-B",
			materialId: "FG0A-MATTE",
			modelId:    "OPPOA3-B",
			expectErr:  false,
		},
		{
			name:       "Parse privacy product",
			input:      "FG0A-PRIVACY-SAMSUNGS25",
			materialId: "FG0A-PRIVACY",
			modelId:    "SAMSUNGS25",
			expectErr:  false,
		},
		{
			name:       "Parse with different film type",
			input:      "FG05-MATTE-OPPOA3",
			materialId: "FG05-MATTE",
			modelId:    "OPPOA3",
			expectErr:  false,
		},
		{
			name:      "Invalid format - too few parts",
			input:     "FG0A-CLEAR",
			expectErr: true,
		},
		{
			name:      "Invalid format - empty string",
			input:     "",
			expectErr: true,
		},
		{
			name:      "Invalid format - single word",
			input:     "INVALID",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			materialId, modelId, err := parser.ParseProductCode(tc.input)

			if tc.expectErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.materialId, materialId)
				assert.Equal(t, tc.modelId, modelId)
			}
		})
	}
}

func TestProductParser_ParseFromFloat64(t *testing.T) {

	parser := parser.NewProductParser()

	testCases := []struct {
		name              string
		platformProductId string
		originalQty       int
		totalPrice        float64
		expectedCount     int
		expectErr         bool

		expectedFirstProduct struct {
			cleanId    string
			quantity   int
			unitPrice  float64
			totalPrice float64
		}
	}{
		{
			name:              "Simple product",
			platformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
			originalQty:       2,
			totalPrice:        100.0,
			expectedCount:     1,
			expectErr:         false,
			expectedFirstProduct: struct {
				cleanId    string
				quantity   int
				unitPrice  float64
				totalPrice float64
			}{
				cleanId:    "FG0A-CLEAR-IPHONE16PROMAX",
				quantity:   2,
				unitPrice:  50.0,
				totalPrice: 100.0,
			},
		},
		{
			name:              "Product with prefix",
			platformProductId: "x2-3&FG0A-CLEAR-IPHONE16PROMAX",
			originalQty:       2,
			totalPrice:        100.0,
			expectedCount:     1,
			expectErr:         false,
			expectedFirstProduct: struct {
				cleanId    string
				quantity   int
				unitPrice  float64
				totalPrice float64
			}{
				cleanId:    "FG0A-CLEAR-IPHONE16PROMAX",
				quantity:   2,
				unitPrice:  50.0,
				totalPrice: 100.0,
			},
		},
		{
			name:              "Product with quantity indicator",
			platformProductId: "FG0A-MATTE-IPHONE16PROMAX*3",
			originalQty:       1,
			totalPrice:        90.0,
			expectedCount:     1,
			expectErr:         false,
			expectedFirstProduct: struct {
				cleanId    string
				quantity   int
				unitPrice  float64
				totalPrice float64
			}{
				cleanId:    "FG0A-MATTE-IPHONE16PROMAX",
				quantity:   3,
				unitPrice:  30.0,
				totalPrice: 90.0,
			},
		},
		{
			name:              "Bundle product",
			platformProductId: "FG0A-CLEAR-OPPOA3/%20xFG0A-CLEAR-OPPOA3-B",
			originalQty:       1,
			totalPrice:        80.0,
			expectedCount:     2,
			expectErr:         false,
			expectedFirstProduct: struct {
				cleanId    string
				quantity   int
				unitPrice  float64
				totalPrice float64
			}{
				cleanId:    "FG0A-CLEAR-OPPOA3",
				quantity:   1,
				unitPrice:  40.0,
				totalPrice: 40.0,
			},
		},
		{
			name:              "Complex bundle with quantity",
			platformProductId: "--FG0A-CLEAR-OPPOA3*2/FG0A-MATTE-OPPOA3",
			originalQty:       1,
			totalPrice:        120.0,
			expectedCount:     2,
			expectErr:         false,
			expectedFirstProduct: struct {
				cleanId    string
				quantity   int
				unitPrice  float64
				totalPrice float64
			}{
				cleanId:    "FG0A-CLEAR-OPPOA3",
				quantity:   2,
				unitPrice:  40.0,
				totalPrice: 80.0,
			},
		},
		{
			name:              "Empty product ID",
			platformProductId: "",
			originalQty:       1,
			totalPrice:        50.0,
			expectedCount:     0,
			expectErr:         true,
		},
		{
			name:              "Invalid total price",
			platformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
			originalQty:       1,
			totalPrice:        -50.0,
			expectedCount:     0,
			expectErr:         true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := parser.ParseFromFloat64(tc.platformProductId, tc.originalQty, tc.totalPrice)

			if tc.expectErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expectedCount, len(result))

			if len(result) > 0 {

				firstProduct := result[0]
				assert.Equal(t, tc.expectedFirstProduct.cleanId, firstProduct.CleanProductId)
				assert.Equal(t, tc.expectedFirstProduct.quantity, firstProduct.Quantity)
				assert.InDelta(t, tc.expectedFirstProduct.unitPrice, firstProduct.UnitPrice.Amount(), 0.01)
				assert.InDelta(t, tc.expectedFirstProduct.totalPrice, firstProduct.TotalPrice.Amount(), 0.01)
			}

			totalCalculated := 0.0
			for _, product := range result {
				totalCalculated += product.TotalPrice.Amount()
			}
			assert.InDelta(t, tc.totalPrice, totalCalculated, 0.01)
		})
	}
}

func TestProductParser_Validate(t *testing.T) {

	parser := parser.NewProductParser()

	testCases := []struct {
		name      string
		input     string
		expectErr bool
	}{
		{
			name:      "Valid product code",
			input:     "FG0A-CLEAR-IPHONE16PROMAX",
			expectErr: false,
		},
		{
			name:      "Valid product code with multiple dashes",
			input:     "FG0A-MATTE-OPPOA3-B",
			expectErr: false,
		},
		{
			name:      "Invalid - empty string",
			input:     "",
			expectErr: true,
		},
		{
			name:      "Invalid - too few dashes",
			input:     "FG0A-CLEAR",
			expectErr: true,
		},
		{
			name:      "Invalid - single word",
			input:     "INVALID",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := parser.Validate(tc.input)

			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
func TestProductParser_EdgeCases(t *testing.T) {

	parser := parser.NewProductParser()

	t.Run("ParseFromFloat64 with zero quantity", func(t *testing.T) {
		result, err := parser.ParseFromFloat64("FG0A-CLEAR-IPHONE16PROMAX", 0, 100.0)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("ParseFromFloat64 with negative quantity", func(t *testing.T) {
		result, err := parser.ParseFromFloat64("FG0A-CLEAR-IPHONE16PROMAX", -1, 100.0)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("CleanPrefix with only prefixes", func(t *testing.T) {
		result := parser.CleanPrefix("x2-3&--")
		assert.Equal(t, "", result)
	})

	t.Run("SplitBundle with empty parts", func(t *testing.T) {
		result := parser.SplitBundle("FG0A-CLEAR-OPPOA3//FG0A-MATTE-OPPOA3")
		expected := []string{"FG0A-CLEAR-OPPOA3", "FG0A-MATTE-OPPOA3"}
		assert.Equal(t, expected, result)
	})

	t.Run("ExtractQuantity with invalid quantity", func(t *testing.T) {
		cleanId, quantity, hasQuantity := parser.ExtractQuantity("FG0A-CLEAR-IPHONE16PROMAX*abc")
		assert.Equal(t, "FG0A-CLEAR-IPHONE16PROMAX*abc", cleanId)
		assert.Equal(t, 1, quantity)
		assert.False(t, hasQuantity)
	})
}

func TestProductParser_ComprehensiveScenarios(t *testing.T) {

	parser := parser.NewProductParser()

	t.Run("Complex nested bundle with mixed quantities", func(t *testing.T) {
		result, err := parser.ParseFromFloat64(
			"x2-3&FG0A-CLEAR-OPPOA3*2/%20xFG0A-MATTE-OPPOA3/FG0A-PRIVACY-SAMSUNGS25*3",
			1,
			300.0,
		)

		require.NoError(t, err)
		assert.Equal(t, 3, len(result))

		totalQuantity := 0
		totalPrice := 0.0
		for _, product := range result {
			totalQuantity += product.Quantity
			totalPrice += product.TotalPrice.Amount()
		}

		assert.Equal(t, 6, totalQuantity)
		assert.InDelta(t, 300.0, totalPrice, 0.01)
	})

	t.Run("Multiple nested prefixes with bundle", func(t *testing.T) {
		result, err := parser.ParseFromFloat64(
			"%20--%20x--x2-3&FG0A-CLEAR-OPPOA3/%20xFG0A-MATTE-OPPOA3",
			1,
			100.0,
		)

		require.NoError(t, err)
		assert.Equal(t, 2, len(result))

		for _, product := range result {
			assert.InDelta(t, 50.0, product.TotalPrice.Amount(), 0.01)
		}
	})
}

func TestPriceCalculator_DividePriceEqually(t *testing.T) {

	calculator := parser.NewPriceCalculator()

	testCases := []struct {
		name        string
		totalPrice  float64
		parts       int
		expected    float64
		expectError bool
	}{
		{
			name:        "Divide 100 by 2 parts",
			totalPrice:  100.0,
			parts:       2,
			expected:    50.0,
			expectError: false,
		},
		{
			name:        "Divide 120 by 3 parts",
			totalPrice:  120.0,
			parts:       3,
			expected:    40.0,
			expectError: false,
		},
		{
			name:        "Divide 100 by 7 parts (with decimal)",
			totalPrice:  100.0,
			parts:       7,
			expected:    14.285714285714286,
			expectError: false,
		},
		{
			name:        "Divide 0 by 5 parts",
			totalPrice:  0.0,
			parts:       5,
			expected:    0.0,
			expectError: false,
		},
		{
			name:        "Divide 50 by 1 part",
			totalPrice:  50.0,
			parts:       1,
			expected:    50.0,
			expectError: false,
		},
		{
			name:        "Divide by 0 parts (should error)",
			totalPrice:  100.0,
			parts:       0,
			expected:    0.0,
			expectError: true,
		},
		{
			name:        "Divide by negative parts (should error)",
			totalPrice:  100.0,
			parts:       -1,
			expected:    0.0,
			expectError: true,
		},
		{
			name:        "Divide large amount by many parts",
			totalPrice:  1000000.0,
			parts:       1000,
			expected:    1000.0,
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			totalPrice, err := value_object.NewPrice(tc.totalPrice)
			require.NoError(t, err)

			result, err := calculator.DividePriceEqually(totalPrice, tc.parts)

			if tc.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.InDelta(t, tc.expected, result.Amount(), 0.0001)
		})
	}
}

func TestPriceCalculator_DividePriceEqually_EdgeCases(t *testing.T) {

	calculator := parser.NewPriceCalculator()

	t.Run("Nil price input", func(t *testing.T) {
		result, err := calculator.DividePriceEqually(nil, 2)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("Very small price division", func(t *testing.T) {
		smallPrice, err := value_object.NewPrice(0.01)
		require.NoError(t, err)

		result, err := calculator.DividePriceEqually(smallPrice, 10)
		require.NoError(t, err)
		assert.InDelta(t, 0.001, result.Amount(), 0.0001)
	})

	t.Run("Large parts number", func(t *testing.T) {
		price, err := value_object.NewPrice(100.0)
		require.NoError(t, err)

		result, err := calculator.DividePriceEqually(price, 10000)
		require.NoError(t, err)
		assert.InDelta(t, 0.01, result.Amount(), 0.0001)
	})
}

func TestPriceCalculator_SumPrices(t *testing.T) {

	calculator := parser.NewPriceCalculator()

	testCases := []struct {
		name        string
		prices      []float64
		expected    float64
		expectError bool
	}{
		{
			name:        "Sum two prices",
			prices:      []float64{50.0, 30.0},
			expected:    80.0,
			expectError: false,
		},
		{
			name:        "Sum three prices",
			prices:      []float64{10.0, 20.0, 30.0},
			expected:    60.0,
			expectError: false,
		},
		{
			name:        "Sum single price",
			prices:      []float64{100.0},
			expected:    100.0,
			expectError: false,
		},
		{
			name:        "Sum with zero values",
			prices:      []float64{50.0, 0.0, 25.0},
			expected:    75.0,
			expectError: false,
		},
		{
			name:        "Sum all zeros",
			prices:      []float64{0.0, 0.0, 0.0},
			expected:    0.0,
			expectError: false,
		},
		{
			name:        "Sum many prices",
			prices:      []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0},
			expected:    55.0,
			expectError: false,
		},
		{
			name:        "Sum decimal prices",
			prices:      []float64{12.50, 25.75, 33.25},
			expected:    71.50,
			expectError: false,
		},
		{
			name:        "Sum very small prices",
			prices:      []float64{0.01, 0.02, 0.03},
			expected:    0.06,
			expectError: false,
		},
		{
			name:        "Sum large prices",
			prices:      []float64{1000000.0, 2000000.0, 3000000.0},
			expected:    6000000.0,
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			var priceObjects []*value_object.Price
			for _, p := range tc.prices {
				price, err := value_object.NewPrice(p)
				require.NoError(t, err)
				priceObjects = append(priceObjects, price)
			}

			result, err := calculator.SumPrices(priceObjects...)

			if tc.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.InDelta(t, tc.expected, result.Amount(), 0.0001)
		})
	}
}

func TestPriceCalculator_SumPrices_EdgeCases(t *testing.T) {

	calculator := parser.NewPriceCalculator()

	t.Run("Sum with no prices (empty varargs)", func(t *testing.T) {
		result, err := calculator.SumPrices()
		require.NoError(t, err)
		assert.Equal(t, 0.0, result.Amount())
	})

	t.Run("Sum with nil prices", func(t *testing.T) {
		price1, err := value_object.NewPrice(50.0)
		require.NoError(t, err)

		price2, err := value_object.NewPrice(30.0)
		require.NoError(t, err)

		result, err := calculator.SumPrices(price1, nil, price2)
		require.NoError(t, err)
		assert.InDelta(t, 80.0, result.Amount(), 0.0001)
	})

	t.Run("Sum with all nil prices", func(t *testing.T) {
		result, err := calculator.SumPrices(nil, nil, nil)
		require.NoError(t, err)
		assert.Equal(t, 0.0, result.Amount())
	})

	t.Run("Sum mixed nil and valid prices", func(t *testing.T) {
		price1, err := value_object.NewPrice(100.0)
		require.NoError(t, err)

		price2, err := value_object.NewPrice(200.0)
		require.NoError(t, err)

		result, err := calculator.SumPrices(nil, price1, nil, price2, nil)
		require.NoError(t, err)
		assert.InDelta(t, 300.0, result.Amount(), 0.0001)
	})

	t.Run("Sum with precision test", func(t *testing.T) {

		price1, err := value_object.NewPrice(0.1)
		require.NoError(t, err)

		price2, err := value_object.NewPrice(0.2)
		require.NoError(t, err)

		price3, err := value_object.NewPrice(0.3)
		require.NoError(t, err)

		result, err := calculator.SumPrices(price1, price2, price3)
		require.NoError(t, err)
		assert.InDelta(t, 0.6, result.Amount(), 0.0001)
	})
}
func TestPriceCalculator_Integration(t *testing.T) {

	calculator := parser.NewPriceCalculator()

	t.Run("Divide then sum back", func(t *testing.T) {
		originalPrice, err := value_object.NewPrice(120.0)
		require.NoError(t, err)

		dividedPrice, err := calculator.DividePriceEqually(originalPrice, 3)
		require.NoError(t, err)

		result, err := calculator.SumPrices(dividedPrice, dividedPrice, dividedPrice)
		require.NoError(t, err)

		assert.InDelta(t, 120.0, result.Amount(), 0.0001)
	})

	t.Run("Complex calculation scenario", func(t *testing.T) {
		totalBundlePrice, err := value_object.NewPrice(300.0)
		require.NoError(t, err)

		pricePerProduct, err := calculator.DividePriceEqually(totalBundlePrice, 4)
		require.NoError(t, err)

		firstTypeTotal, err := calculator.CalculateTotalPrice(pricePerProduct, 2)
		require.NoError(t, err)

		secondTypeTotal, err := calculator.CalculateTotalPrice(pricePerProduct, 2)
		require.NoError(t, err)

		grandTotal, err := calculator.SumPrices(firstTypeTotal, secondTypeTotal)
		require.NoError(t, err)

		assert.InDelta(t, 300.0, grandTotal.Amount(), 0.0001)
	})

	t.Run("Real-world bundle scenario", func(t *testing.T) {
		totalPrice, err := value_object.NewPrice(120.0)
		require.NoError(t, err)

		unitPrice, err := calculator.CalculateUnitPrice(totalPrice, 6)
		require.NoError(t, err)

		product1Total, err := calculator.CalculateTotalPrice(unitPrice, 2)
		require.NoError(t, err)

		product2Total, err := calculator.CalculateTotalPrice(unitPrice, 1)
		require.NoError(t, err)

		product3Total, err := calculator.CalculateTotalPrice(unitPrice, 3)
		require.NoError(t, err)

		calculatedTotal, err := calculator.SumPrices(product1Total, product2Total, product3Total)
		require.NoError(t, err)

		assert.InDelta(t, 20.0, unitPrice.Amount(), 0.0001)
		assert.InDelta(t, 40.0, product1Total.Amount(), 0.0001)
		assert.InDelta(t, 20.0, product2Total.Amount(), 0.0001)
		assert.InDelta(t, 60.0, product3Total.Amount(), 0.0001)
		assert.InDelta(t, 120.0, calculatedTotal.Amount(), 0.0001)
	})
}

func TestPriceCalculator_Performance(t *testing.T) {

	calculator := parser.NewPriceCalculator()

	t.Run("Sum many prices performance", func(t *testing.T) {
		var prices []*value_object.Price
		for i := 0; i < 1000; i++ {
			price, err := value_object.NewPrice(float64(i + 1))
			require.NoError(t, err)
			prices = append(prices, price)
		}

		result, err := calculator.SumPrices(prices...)
		require.NoError(t, err)

		expected := float64(1000 * 1001 / 2)
		assert.InDelta(t, expected, result.Amount(), 0.0001)
	})

	t.Run("Division with large numbers", func(t *testing.T) {
		largePrice, err := value_object.NewPrice(1000000000.0)
		require.NoError(t, err)

		result, err := calculator.DividePriceEqually(largePrice, 1000000)
		require.NoError(t, err)

		assert.InDelta(t, 1000.0, result.Amount(), 0.0001)
	})
}
