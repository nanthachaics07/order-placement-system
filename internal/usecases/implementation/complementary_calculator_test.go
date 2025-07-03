package implementation_test

import (
	"fmt"
	"order-placement-system/internal/domain/entity"
	"order-placement-system/internal/domain/value_object"
	"order-placement-system/internal/usecases/implementation"
	"order-placement-system/pkg/log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	log.Init("dev")
}

func TestComplementaryCalculatorUseCase_CalculateWithStartingOrderNo(t *testing.T) {

	tests := []struct {
		name            string
		mainProducts    []*entity.Product
		startingOrderNo int
		expectedOrders  []*entity.CleanedOrder
		expectedError   bool
		errorMessage    string
	}{
		{
			name:            "Empty products list",
			mainProducts:    []*entity.Product{},
			startingOrderNo: 1,
			expectedOrders:  []*entity.CleanedOrder{},
			expectedError:   false,
		},
		{
			name:            "Nil products list",
			mainProducts:    nil,
			startingOrderNo: 1,
			expectedOrders:  []*entity.CleanedOrder{},
			expectedError:   false,
		},
		{
			name: "Single CLEAR product",
			mainProducts: []*entity.Product{
				{
					ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
					MaterialId: "FG0A-CLEAR",
					ModelId:    "IPHONE16PROMAX",
					Quantity:   2,
					UnitPrice:  value_object.MustNewPrice(50.00),
					TotalPrice: value_object.MustNewPrice(100.00),
				},
			},
			startingOrderNo: 2,
			expectedOrders: []*entity.CleanedOrder{
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
			expectedError: false,
		},
		{
			name: "Single MATTE product",
			mainProducts: []*entity.Product{
				{
					ProductId:  "FG0A-MATTE-IPHONE16PROMAX",
					MaterialId: "FG0A-MATTE",
					ModelId:    "IPHONE16PROMAX",
					Quantity:   3,
					UnitPrice:  value_object.MustNewPrice(30.00),
					TotalPrice: value_object.MustNewPrice(90.00),
				},
			},
			startingOrderNo: 2,
			expectedOrders: []*entity.CleanedOrder{
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
			expectedError: false,
		},
		{
			name: "Single PRIVACY product",
			mainProducts: []*entity.Product{
				{
					ProductId:  "FG0A-PRIVACY-IPHONE16PROMAX",
					MaterialId: "FG0A-PRIVACY",
					ModelId:    "IPHONE16PROMAX",
					Quantity:   1,
					UnitPrice:  value_object.MustNewPrice(50.00),
					TotalPrice: value_object.MustNewPrice(50.00),
				},
			},
			startingOrderNo: 4,
			expectedOrders: []*entity.CleanedOrder{
				{
					No:         4,
					ProductId:  "WIPING-CLOTH",
					Qty:        1,
					UnitPrice:  value_object.ZeroPrice(),
					TotalPrice: value_object.ZeroPrice(),
				},
				{
					No:         5,
					ProductId:  "PRIVACY-CLEANNER",
					Qty:        1,
					UnitPrice:  value_object.ZeroPrice(),
					TotalPrice: value_object.ZeroPrice(),
				},
			},
			expectedError: false,
		},
		{
			name: "Multiple products same texture",
			mainProducts: []*entity.Product{
				{
					ProductId:  "FG0A-CLEAR-OPPOA3",
					MaterialId: "FG0A-CLEAR",
					ModelId:    "OPPOA3",
					Quantity:   1,
					UnitPrice:  value_object.MustNewPrice(40.00),
					TotalPrice: value_object.MustNewPrice(40.00),
				},
				{
					ProductId:  "FG0A-CLEAR-OPPOA3-B",
					MaterialId: "FG0A-CLEAR",
					ModelId:    "OPPOA3-B",
					Quantity:   1,
					UnitPrice:  value_object.MustNewPrice(40.00),
					TotalPrice: value_object.MustNewPrice(40.00),
				},
			},
			startingOrderNo: 3,
			expectedOrders: []*entity.CleanedOrder{
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
			expectedError: false,
		},
		{
			name: "Multiple products different textures",
			mainProducts: []*entity.Product{
				{
					ProductId:  "FG0A-CLEAR-OPPOA3",
					MaterialId: "FG0A-CLEAR",
					ModelId:    "OPPOA3",
					Quantity:   1,
					UnitPrice:  value_object.MustNewPrice(40.00),
					TotalPrice: value_object.MustNewPrice(40.00),
				},
				{
					ProductId:  "FG0A-CLEAR-OPPOA3-B",
					MaterialId: "FG0A-CLEAR",
					ModelId:    "OPPOA3-B",
					Quantity:   1,
					UnitPrice:  value_object.MustNewPrice(40.00),
					TotalPrice: value_object.MustNewPrice(40.00),
				},
				{
					ProductId:  "FG0A-MATTE-OPPOA3",
					MaterialId: "FG0A-MATTE",
					ModelId:    "OPPOA3",
					Quantity:   1,
					UnitPrice:  value_object.MustNewPrice(40.00),
					TotalPrice: value_object.MustNewPrice(40.00),
				},
			},
			startingOrderNo: 4,
			expectedOrders: []*entity.CleanedOrder{
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
			expectedError: false,
		},
		{
			name: "Complex case with multiple textures and quantities",
			mainProducts: []*entity.Product{
				{
					ProductId:  "FG0A-CLEAR-OPPOA3",
					MaterialId: "FG0A-CLEAR",
					ModelId:    "OPPOA3",
					Quantity:   2,
					UnitPrice:  value_object.MustNewPrice(40.00),
					TotalPrice: value_object.MustNewPrice(80.00),
				},
				{
					ProductId:  "FG0A-MATTE-OPPOA3",
					MaterialId: "FG0A-MATTE",
					ModelId:    "OPPOA3",
					Quantity:   2,
					UnitPrice:  value_object.MustNewPrice(40.00),
					TotalPrice: value_object.MustNewPrice(80.00),
				},
				{
					ProductId:  "FG0A-PRIVACY-IPHONE16PROMAX",
					MaterialId: "FG0A-PRIVACY",
					ModelId:    "IPHONE16PROMAX",
					Quantity:   1,
					UnitPrice:  value_object.MustNewPrice(50.00),
					TotalPrice: value_object.MustNewPrice(50.00),
				},
			},
			startingOrderNo: 4,
			expectedOrders: []*entity.CleanedOrder{
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
			expectedError: false,
		},
		{
			name: "Product with invalid material id format",
			mainProducts: []*entity.Product{
				{
					ProductId:  "INVALID-PRODUCT",
					MaterialId: "INVALID",
					ModelId:    "MODEL",
					Quantity:   1,
					UnitPrice:  value_object.MustNewPrice(50.00),
					TotalPrice: value_object.MustNewPrice(50.00),
				},
			},
			startingOrderNo: 2,
			expectedOrders:  nil,
			expectedError:   true,
			errorMessage:    "invalid input",
		},
		{
			name: "Product with invalid texture",
			mainProducts: []*entity.Product{
				{
					ProductId:  "FG0A-INVALID-IPHONE16PROMAX",
					MaterialId: "FG0A-INVALID",
					ModelId:    "IPHONE16PROMAX",
					Quantity:   1,
					UnitPrice:  value_object.MustNewPrice(50.00),
					TotalPrice: value_object.MustNewPrice(50.00),
				},
			},
			startingOrderNo: 2,
			expectedOrders:  nil,
			expectedError:   true,
			errorMessage:    "invalid input",
		},
		{
			name: "Product with zero quantity",
			mainProducts: []*entity.Product{
				{
					ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
					MaterialId: "FG0A-CLEAR",
					ModelId:    "IPHONE16PROMAX",
					Quantity:   0,
					UnitPrice:  value_object.MustNewPrice(50.00),
					TotalPrice: value_object.ZeroPrice(),
				},
			},
			startingOrderNo: 1,
			expectedOrders:  []*entity.CleanedOrder{},
			expectedError:   false,
		},
		{
			name: "Product with empty material id",
			mainProducts: []*entity.Product{
				{
					ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
					MaterialId: "",
					ModelId:    "IPHONE16PROMAX",
					Quantity:   1,
					UnitPrice:  value_object.MustNewPrice(50.00),
					TotalPrice: value_object.MustNewPrice(50.00),
				},
			},
			startingOrderNo: 1,
			expectedOrders:  nil,
			expectedError:   true,
			errorMessage:    "invalid input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calculator := implementation.NewComplementaryCalculator()

			result, err := calculator.CalculateWithStartingOrderNo(tt.mainProducts, tt.startingOrderNo)

			if tt.expectedError {
				assert.Error(t, err, "Expected error but got none")
				if tt.errorMessage != "" {
					assert.Contains(t, err.Error(), tt.errorMessage, "Error message should contain expected text")
				}
				assert.Nil(t, result, "Result should be nil when error occurs")
			} else {
				assert.NoError(t, err, "Expected no error but got: %v", err)

				// Handle both nil and empty slice cases
				if len(tt.expectedOrders) == 0 {
					assert.True(t, len(result) == 0,
						"Result should be nil or empty when no orders expected")
				} else {
					assert.NotNil(t, result, "Result should not be nil when orders are expected")
					assert.Equal(t, len(tt.expectedOrders), len(result), "Number of orders should match")

					// Compare each order
					for i, expected := range tt.expectedOrders {
						require.Less(t, i, len(result), "Result should have at least %d orders", i+1)

						actual := result[i]
						assert.Equal(t, expected.No, actual.No, "Order No mismatch at index %d", i)
						assert.Equal(t, expected.ProductId, actual.ProductId, "ProductId mismatch at index %d", i)
						assert.Equal(t, expected.Qty, actual.Qty, "Qty mismatch at index %d", i)

						// Compare prices using the Equals method
						assert.True(t, expected.UnitPrice.Equals(actual.UnitPrice),
							"UnitPrice mismatch at index %d: expected %s, got %s", i,
							expected.UnitPrice.String(), actual.UnitPrice.String())
						assert.True(t, expected.TotalPrice.Equals(actual.TotalPrice),
							"TotalPrice mismatch at index %d: expected %s, got %s", i,
							expected.TotalPrice.String(), actual.TotalPrice.String())

						// Complementary items should not have MaterialId and ModelId
						assert.Empty(t, actual.MaterialId, "Complementary item should not have MaterialId")
						assert.Empty(t, actual.ModelId, "Complementary item should not have ModelId")
					}
				}
			}
		})
	}
}

func TestComplementaryCalculatorUseCase_StartingOrderNoSequence(t *testing.T) {

	calculator := implementation.NewComplementaryCalculator()

	products := []*entity.Product{
		{
			ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
			MaterialId: "FG0A-CLEAR",
			ModelId:    "IPHONE16PROMAX",
			Quantity:   1,
			UnitPrice:  value_object.MustNewPrice(50.00),
			TotalPrice: value_object.MustNewPrice(50.00),
		},
	}

	// Test with different starting order numbers
	testCases := []struct {
		startingOrderNo int
		expectedNos     []int
	}{
		{1, []int{1, 2}},
		{10, []int{10, 11}},
		{100, []int{100, 101}},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("StartingOrderNo_%d", tc.startingOrderNo), func(t *testing.T) {
			result, err := calculator.CalculateWithStartingOrderNo(products, tc.startingOrderNo)

			require.NoError(t, err)
			require.Equal(t, len(tc.expectedNos), len(result))

			for i, expectedNo := range tc.expectedNos {
				assert.Equal(t, expectedNo, result[i].No)
			}
		})
	}
}

func TestComplementaryCalculatorUseCase_TextureOrdering(t *testing.T) {

	calculator := implementation.NewComplementaryCalculator()

	products := []*entity.Product{
		{
			ProductId:  "FG0A-PRIVACY-IPHONE16PROMAX",
			MaterialId: "FG0A-PRIVACY",
			ModelId:    "IPHONE16PROMAX",
			Quantity:   1,
			UnitPrice:  value_object.MustNewPrice(50.00),
			TotalPrice: value_object.MustNewPrice(50.00),
		},
		{
			ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
			MaterialId: "FG0A-CLEAR",
			ModelId:    "IPHONE16PROMAX",
			Quantity:   1,
			UnitPrice:  value_object.MustNewPrice(50.00),
			TotalPrice: value_object.MustNewPrice(50.00),
		},
		{
			ProductId:  "FG0A-MATTE-IPHONE16PROMAX",
			MaterialId: "FG0A-MATTE",
			ModelId:    "IPHONE16PROMAX",
			Quantity:   1,
			UnitPrice:  value_object.MustNewPrice(50.00),
			TotalPrice: value_object.MustNewPrice(50.00),
		},
	}

	result, err := calculator.CalculateWithStartingOrderNo(products, 1)

	require.NoError(t, err)
	require.Equal(t, 4, len(result))

	// Expected order: WIPING-CLOTH first, then cleaners in texture order
	expectedOrder := []string{
		"WIPING-CLOTH",
		"CLEAR-CLEANNER",
		"MATTE-CLEANNER",
		"PRIVACY-CLEANNER",
	}

	for i, expected := range expectedOrder {
		assert.Equal(t, expected, result[i].ProductId, "Wrong order at index %d", i)
	}
}

func TestComplementaryCalculatorUseCase_EdgeCases(t *testing.T) {

	calculator := implementation.NewComplementaryCalculator()

	t.Run("Negative starting order number", func(t *testing.T) {
		products := []*entity.Product{
			{
				ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
				MaterialId: "FG0A-CLEAR",
				ModelId:    "IPHONE16PROMAX",
				Quantity:   1,
				UnitPrice:  value_object.MustNewPrice(50.00),
				TotalPrice: value_object.MustNewPrice(50.00),
			},
		}

		result, err := calculator.CalculateWithStartingOrderNo(products, -1)

		// Should still work but with negative order numbers
		require.NoError(t, err)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, -1, result[0].No)
		assert.Equal(t, 0, result[1].No)
	})

	t.Run("Very large quantities", func(t *testing.T) {
		products := []*entity.Product{
			{
				ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
				MaterialId: "FG0A-CLEAR",
				ModelId:    "IPHONE16PROMAX",
				Quantity:   1000,
				UnitPrice:  value_object.MustNewPrice(50.00),
				TotalPrice: value_object.MustNewPrice(50000.00),
			},
		}

		result, err := calculator.CalculateWithStartingOrderNo(products, 1)

		require.NoError(t, err)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, 1000, result[0].Qty) // WIPING-CLOTH
		assert.Equal(t, 1000, result[1].Qty) // CLEAR-CLEANNER
	})
}

// func TestComplementaryCalculatorUseCase_NilProduct(t *testing.T) {
//
// 	calculator := implementation.NewComplementaryCalculator()

// 	products := []*entity.Product{
// 		nil,
// 		{
// 			ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
// 			MaterialId: "FG0A-CLEAR",
// 			ModelId:    "IPHONE16PROMAX",
// 			Quantity:   1,
// 			UnitPrice:  value_object.MustNewPrice(50.00),
// 			TotalPrice: value_object.MustNewPrice(50.00),
// 		},
// 	}

// 	defer func() {
// 		if r := recover(); r != nil {
// 			t.Logf("Expected panic occurred: %v", r)
// 		}
// 	}()

// 	result, err := calculator.CalculateWithStartingOrderNo(products, 1)

// 	if err != nil {
// 		assert.Error(t, err)
// 		assert.Nil(t, result)
// 		assert.Contains(t, err.Error(), "invalid input")
// 	} else {

// 		t.Fatal("Expected error or panic when nil product is provided")
// 	}
// }
