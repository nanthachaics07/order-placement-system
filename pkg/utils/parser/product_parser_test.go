package parser_test

import (
	"order-placement-system/pkg/log"
	"order-placement-system/pkg/utils/parser"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductParser_CleanPrefix(t *testing.T) {
	log.Init("dev")
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
	log.Init("dev")
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
	log.Init("dev")
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
	log.Init("dev")
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

func TestProductParser_Parse(t *testing.T) {
	log.Init("dev")
	parser := parser.NewProductParser()

	testCases := []struct {
		name              string
		platformProductId string
		originalQty       int
		totalPrice        float64
		expectedCount     int
		expectErr         bool
	}{
		{
			name:              "Simple product",
			platformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
			originalQty:       2,
			totalPrice:        100.0,
			expectedCount:     1,
			expectErr:         false,
		},
		{
			name:              "Product with prefix",
			platformProductId: "x2-3&FG0A-CLEAR-IPHONE16PROMAX",
			originalQty:       2,
			totalPrice:        100.0,
			expectedCount:     1,
			expectErr:         false,
		},
		{
			name:              "Product with quantity indicator",
			platformProductId: "FG0A-MATTE-IPHONE16PROMAX*3",
			originalQty:       1,
			totalPrice:        90.0,
			expectedCount:     1,
			expectErr:         false,
		},
		{
			name:              "Bundle product",
			platformProductId: "FG0A-CLEAR-OPPOA3/%20xFG0A-CLEAR-OPPOA3-B",
			originalQty:       1,
			totalPrice:        80.0,
			expectedCount:     2,
			expectErr:         false,
		},
		{
			name:              "Complex bundle with quantity",
			platformProductId: "--FG0A-CLEAR-OPPOA3*2/FG0A-MATTE-OPPOA3",
			originalQty:       1,
			totalPrice:        120.0,
			expectedCount:     2,
			expectErr:         false,
		},
		{
			name:              "Empty product ID",
			platformProductId: "",
			originalQty:       1,
			totalPrice:        50.0,
			expectedCount:     0,
			expectErr:         true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := parser.Parse(tc.platformProductId, tc.originalQty, tc.totalPrice)

			if tc.expectErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedCount, len(result))

				totalCalculated := 0.0
				for _, product := range result {
					totalCalculated += product.TotalPrice
				}
				assert.InDelta(t, tc.totalPrice, totalCalculated, 0.01)
			}
		})
	}
}

func TestProductParser_Validate(t *testing.T) {
	log.Init("dev")
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
