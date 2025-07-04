package implementation_test

import (
	"testing"

	"order-placement-system/internal/domain/entity"
	"order-placement-system/internal/domain/value_object"
	"order-placement-system/internal/usecases/implementation"
	"order-placement-system/pkg/log"
	"order-placement-system/pkg/utils/parser"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	log.Init("dev")
}

func TestOrderProcessor_ProcessOrders_SevenCases(t *testing.T) {

	processor := implementation.NewOrderProcessor(
		parser.NewProductParser(),
		implementation.NewComplementaryCalculator(),
	)

	testCases := []struct {
		name     string
		input    []*entity.InputOrder
		expected []*entity.CleanedOrder
	}{
		{
			name: "Case 1: Only one product",
			input: []*entity.InputOrder{
				{
					No:                1,
					PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
					Qty:               2,
					UnitPrice:         value_object.MustNewPrice(50),
					TotalPrice:        value_object.MustNewPrice(100),
				},
			},
			expected: []*entity.CleanedOrder{
				{
					No:         1,
					ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
					MaterialId: "FG0A-CLEAR",
					ModelId:    "IPHONE16PROMAX",
					Qty:        2,
					UnitPrice:  value_object.MustNewPrice(50.00),
					TotalPrice: value_object.MustNewPrice(100.00),
				},
				{
					No:         2,
					ProductId:  "WIPING-CLOTH",
					Qty:        2,
					UnitPrice:  value_object.ZeroPrice(),
					TotalPrice: value_object.ZeroPrice(),
				},
				{
					No:         3,
					ProductId:  "CLEAR-CLEANNER",
					Qty:        2,
					UnitPrice:  value_object.ZeroPrice(),
					TotalPrice: value_object.ZeroPrice(),
				},
			},
		},
		{
			name: "Case 2: One product with wrong prefix",
			input: []*entity.InputOrder{
				{
					No:                1,
					PlatformProductId: "x2-3&FG0A-CLEAR-IPHONE16PROMAX",
					Qty:               2,
					UnitPrice:         value_object.MustNewPrice(50),
					TotalPrice:        value_object.MustNewPrice(100),
				},
			},
			expected: []*entity.CleanedOrder{
				{
					No:         1,
					ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
					MaterialId: "FG0A-CLEAR",
					ModelId:    "IPHONE16PROMAX",
					Qty:        2,
					UnitPrice:  value_object.MustNewPrice(50.00),
					TotalPrice: value_object.MustNewPrice(100.00),
				},
				{
					No:         2,
					ProductId:  "WIPING-CLOTH",
					Qty:        2,
					UnitPrice:  value_object.ZeroPrice(),
					TotalPrice: value_object.ZeroPrice(),
				},
				{
					No:         3,
					ProductId:  "CLEAR-CLEANNER",
					Qty:        2,
					UnitPrice:  value_object.ZeroPrice(),
					TotalPrice: value_object.ZeroPrice(),
				},
			},
		},
		{
			name: "Case 3: One product with wrong prefix and * symbol",
			input: []*entity.InputOrder{
				{
					No:                1,
					PlatformProductId: "x2-3&FG0A-MATTE-IPHONE16PROMAX*3",
					Qty:               1,
					UnitPrice:         value_object.MustNewPrice(90),
					TotalPrice:        value_object.MustNewPrice(90),
				},
			},
			expected: []*entity.CleanedOrder{
				{
					No:         1,
					ProductId:  "FG0A-MATTE-IPHONE16PROMAX",
					MaterialId: "FG0A-MATTE",
					ModelId:    "IPHONE16PROMAX",
					Qty:        3,
					UnitPrice:  value_object.MustNewPrice(30.00),
					TotalPrice: value_object.MustNewPrice(90.00),
				},
				{
					No:         2,
					ProductId:  "WIPING-CLOTH",
					Qty:        3,
					UnitPrice:  value_object.ZeroPrice(),
					TotalPrice: value_object.ZeroPrice(),
				},
				{
					No:         3,
					ProductId:  "MATTE-CLEANNER",
					Qty:        3,
					UnitPrice:  value_object.ZeroPrice(),
					TotalPrice: value_object.ZeroPrice(),
				},
			},
		},
		{
			name: "Case 4: Bundle product with two items",
			input: []*entity.InputOrder{
				{
					No:                1,
					PlatformProductId: "FG0A-CLEAR-OPPOA3/%20xFG0A-CLEAR-OPPOA3-B",
					Qty:               1,
					UnitPrice:         value_object.MustNewPrice(80),
					TotalPrice:        value_object.MustNewPrice(80),
				},
			},
			expected: []*entity.CleanedOrder{
				{
					No:         1,
					ProductId:  "FG0A-CLEAR-OPPOA3",
					MaterialId: "FG0A-CLEAR",
					ModelId:    "OPPOA3",
					Qty:        1,
					UnitPrice:  value_object.MustNewPrice(40.00),
					TotalPrice: value_object.MustNewPrice(40.00),
				},
				{
					No:         2,
					ProductId:  "FG0A-CLEAR-OPPOA3-B",
					MaterialId: "FG0A-CLEAR",
					ModelId:    "OPPOA3-B",
					Qty:        1,
					UnitPrice:  value_object.MustNewPrice(40.00),
					TotalPrice: value_object.MustNewPrice(40.00),
				},
				{
					No:         3,
					ProductId:  "WIPING-CLOTH",
					Qty:        2,
					UnitPrice:  value_object.ZeroPrice(),
					TotalPrice: value_object.ZeroPrice(),
				},
				{
					No:         4,
					ProductId:  "CLEAR-CLEANNER",
					Qty:        2,
					UnitPrice:  value_object.ZeroPrice(),
					TotalPrice: value_object.ZeroPrice(),
				},
			},
		},
		{
			name: "Case 5: Bundle product with three items",
			input: []*entity.InputOrder{
				{
					No:                1,
					PlatformProductId: "FG0A-CLEAR-OPPOA3/%20xFG0A-CLEAR-OPPOA3-B/FG0A-MATTE-OPPOA3",
					Qty:               1,
					UnitPrice:         value_object.MustNewPrice(120),
					TotalPrice:        value_object.MustNewPrice(120),
				},
			},
			expected: []*entity.CleanedOrder{
				{
					No:         1,
					ProductId:  "FG0A-CLEAR-OPPOA3",
					MaterialId: "FG0A-CLEAR",
					ModelId:    "OPPOA3",
					Qty:        1,
					UnitPrice:  value_object.MustNewPrice(40.00),
					TotalPrice: value_object.MustNewPrice(40.00),
				},
				{
					No:         2,
					ProductId:  "FG0A-CLEAR-OPPOA3-B",
					MaterialId: "FG0A-CLEAR",
					ModelId:    "OPPOA3-B",
					Qty:        1,
					UnitPrice:  value_object.MustNewPrice(40.00),
					TotalPrice: value_object.MustNewPrice(40.00),
				},
				{
					No:         3,
					ProductId:  "FG0A-MATTE-OPPOA3",
					MaterialId: "FG0A-MATTE",
					ModelId:    "OPPOA3",
					Qty:        1,
					UnitPrice:  value_object.MustNewPrice(40.00),
					TotalPrice: value_object.MustNewPrice(40.00),
				},
				{
					No:         4,
					ProductId:  "WIPING-CLOTH",
					Qty:        3,
					UnitPrice:  value_object.ZeroPrice(),
					TotalPrice: value_object.ZeroPrice(),
				},
				{
					No:         5,
					ProductId:  "CLEAR-CLEANNER",
					Qty:        2,
					UnitPrice:  value_object.ZeroPrice(),
					TotalPrice: value_object.ZeroPrice(),
				},
				{
					No:         6,
					ProductId:  "MATTE-CLEANNER",
					Qty:        1,
					UnitPrice:  value_object.ZeroPrice(),
					TotalPrice: value_object.ZeroPrice(),
				},
			},
		},
		{
			name: "Case 6: Bundle with * symbol",
			input: []*entity.InputOrder{
				{
					No:                1,
					PlatformProductId: "--FG0A-CLEAR-OPPOA3*2/FG0A-MATTE-OPPOA3",
					Qty:               1,
					UnitPrice:         value_object.MustNewPrice(120),
					TotalPrice:        value_object.MustNewPrice(120),
				},
			},
			expected: []*entity.CleanedOrder{
				{
					No:         1,
					ProductId:  "FG0A-CLEAR-OPPOA3",
					MaterialId: "FG0A-CLEAR",
					ModelId:    "OPPOA3",
					Qty:        2,
					UnitPrice:  value_object.MustNewPrice(40.00),
					TotalPrice: value_object.MustNewPrice(80.00),
				},
				{
					No:         2,
					ProductId:  "FG0A-MATTE-OPPOA3",
					MaterialId: "FG0A-MATTE",
					ModelId:    "OPPOA3",
					Qty:        1,
					UnitPrice:  value_object.MustNewPrice(40.00),
					TotalPrice: value_object.MustNewPrice(40.00),
				},
				{
					No:         3,
					ProductId:  "WIPING-CLOTH",
					Qty:        3,
					UnitPrice:  value_object.ZeroPrice(),
					TotalPrice: value_object.ZeroPrice(),
				},
				{
					No:         4,
					ProductId:  "CLEAR-CLEANNER",
					Qty:        2,
					UnitPrice:  value_object.ZeroPrice(),
					TotalPrice: value_object.ZeroPrice(),
				},
				{
					No:         5,
					ProductId:  "MATTE-CLEANNER",
					Qty:        1,
					UnitPrice:  value_object.ZeroPrice(),
					TotalPrice: value_object.ZeroPrice(),
				},
			},
		},
		{
			name: "Case 7: Multiple products",
			input: []*entity.InputOrder{
				{
					No:                1,
					PlatformProductId: "--FG0A-CLEAR-OPPOA3*2/FG0A-MATTE-OPPOA3*2",
					Qty:               1,
					UnitPrice:         value_object.MustNewPrice(160),
					TotalPrice:        value_object.MustNewPrice(160),
				},
				{
					No:                2,
					PlatformProductId: "FG0A-PRIVACY-IPHONE16PROMAX",
					Qty:               1,
					UnitPrice:         value_object.MustNewPrice(50),
					TotalPrice:        value_object.MustNewPrice(50),
				},
			},
			expected: []*entity.CleanedOrder{
				{
					No:         1,
					ProductId:  "FG0A-CLEAR-OPPOA3",
					MaterialId: "FG0A-CLEAR",
					ModelId:    "OPPOA3",
					Qty:        2,
					UnitPrice:  value_object.MustNewPrice(40.00),
					TotalPrice: value_object.MustNewPrice(80.00),
				},
				{
					No:         2,
					ProductId:  "FG0A-MATTE-OPPOA3",
					MaterialId: "FG0A-MATTE",
					ModelId:    "OPPOA3",
					Qty:        2,
					UnitPrice:  value_object.MustNewPrice(40.00),
					TotalPrice: value_object.MustNewPrice(80.00),
				},
				{
					No:         3,
					ProductId:  "FG0A-PRIVACY-IPHONE16PROMAX",
					MaterialId: "FG0A-PRIVACY",
					ModelId:    "IPHONE16PROMAX",
					Qty:        1,
					UnitPrice:  value_object.MustNewPrice(50.00),
					TotalPrice: value_object.MustNewPrice(50.00),
				},
				{
					No:         4,
					ProductId:  "WIPING-CLOTH",
					Qty:        5,
					UnitPrice:  value_object.ZeroPrice(),
					TotalPrice: value_object.ZeroPrice(),
				},
				{
					No:         5,
					ProductId:  "CLEAR-CLEANNER",
					Qty:        2,
					UnitPrice:  value_object.ZeroPrice(),
					TotalPrice: value_object.ZeroPrice(),
				},
				{
					No:         6,
					ProductId:  "MATTE-CLEANNER",
					Qty:        2,
					UnitPrice:  value_object.ZeroPrice(),
					TotalPrice: value_object.ZeroPrice(),
				},
				{
					No:         7,
					ProductId:  "PRIVACY-CLEANNER",
					Qty:        1,
					UnitPrice:  value_object.ZeroPrice(),
					TotalPrice: value_object.ZeroPrice(),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result, err := processor.ProcessOrders(tc.input)

			// Assert
			require.NoError(t, err, "ProcessOrders should not return error")
			require.Equal(t, len(tc.expected), len(result), "Result length should match expected")

			for i, expectedOrder := range tc.expected {
				assert.Equal(t, expectedOrder.No, result[i].No, "Order number should match")
				assert.Equal(t, expectedOrder.ProductId, result[i].ProductId, "Product ID should match")
				assert.Equal(t, expectedOrder.MaterialId, result[i].MaterialId, "Material ID should match")
				assert.Equal(t, expectedOrder.ModelId, result[i].ModelId, "Model ID should match")
				assert.Equal(t, expectedOrder.Qty, result[i].Qty, "Quantity should match")

				// เปรียบเทียบ Price objects
				if expectedOrder.UnitPrice != nil && result[i].UnitPrice != nil {
					assert.InDelta(t, expectedOrder.UnitPrice.Amount(), result[i].UnitPrice.Amount(), 0.01,
						"Unit price should match")
				} else {
					assert.Equal(t, expectedOrder.UnitPrice, result[i].UnitPrice, "Unit price should match")
				}

				if expectedOrder.TotalPrice != nil && result[i].TotalPrice != nil {
					assert.InDelta(t, expectedOrder.TotalPrice.Amount(), result[i].TotalPrice.Amount(), 0.01,
						"Total price should match")
				} else {
					assert.Equal(t, expectedOrder.TotalPrice, result[i].TotalPrice, "Total price should match")
				}
			}
		})
	}
}

func TestOrderProcessor_EdgeCases(t *testing.T) {

	processor := implementation.NewOrderProcessor(
		parser.NewProductParser(),
		implementation.NewComplementaryCalculator(),
	)

	t.Run("Empty input", func(t *testing.T) {
		result, err := processor.ProcessOrders([]*entity.InputOrder{})
		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("Invalid product ID", func(t *testing.T) {
		input := []*entity.InputOrder{
			{
				No:                1,
				PlatformProductId: "INVALID-ID",
				Qty:               1,
				UnitPrice:         value_object.MustNewPrice(50),
				TotalPrice:        value_object.MustNewPrice(50),
			},
		}

		_, err := processor.ProcessOrders(input)
		assert.Error(t, err)
	})

	t.Run("Nil input order", func(t *testing.T) {
		input := []*entity.InputOrder{nil}

		_, err := processor.ProcessOrders(input)
		assert.Error(t, err)
	})

	t.Run("Zero quantity", func(t *testing.T) {
		input := []*entity.InputOrder{
			{
				No:                1,
				PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
				Qty:               0,
				UnitPrice:         value_object.MustNewPrice(50),
				TotalPrice:        value_object.MustNewPrice(0),
			},
		}

		_, err := processor.ProcessOrders(input)
		assert.Error(t, err)
	})

	// t.Run("Negative price", func(t *testing.T) {
	// 	input := []*entity.InputOrder{
	// 		{
	// 			No:                1,
	// 			PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
	// 			Qty:               1,
	// 			UnitPrice:         value_object.MustNewPrice(-10),
	// 			TotalPrice:        value_object.MustNewPrice(-10),
	// 		},
	// 	}

	// 	_, err := processor.ProcessOrders(input)
	// 	assert.Error(t, err)
	// })
}
