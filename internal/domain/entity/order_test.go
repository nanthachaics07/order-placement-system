package entity_test

import (
	"testing"

	"order-placement-system/internal/domain/entity"
	"order-placement-system/internal/domain/value_object"
	"order-placement-system/pkg/errors"
	"order-placement-system/pkg/log"

	"github.com/stretchr/testify/assert"
)

func init() {
	log.Init("dev")
}

func TestInputOrder_IsValid(t *testing.T) {
	tests := []struct {
		name        string
		order       *entity.InputOrder
		expectError bool
		expectedErr error
	}{
		{
			name: "Valid input order",
			order: &entity.InputOrder{
				No:                1,
				PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
				Qty:               2,
				UnitPrice:         value_object.MustNewPrice(50.0),
				TotalPrice:        value_object.MustNewPrice(100.0),
			},
			expectError: false,
		},
		{
			name: "Valid input order with zero price",
			order: &entity.InputOrder{
				No:                1,
				PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
				Qty:               1,
				UnitPrice:         value_object.MustNewPrice(0.0),
				TotalPrice:        value_object.MustNewPrice(0.0),
			},
			expectError: false,
		},
		{
			name: "Invalid order - zero order number",
			order: &entity.InputOrder{
				No:                0,
				PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
				Qty:               2,
				UnitPrice:         value_object.MustNewPrice(50.0),
				TotalPrice:        value_object.MustNewPrice(100.0),
			},
			expectError: true,
			expectedErr: errors.ErrInvalidInput,
		},
		{
			name: "Invalid order - negative order number",
			order: &entity.InputOrder{
				No:                -1,
				PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
				Qty:               2,
				UnitPrice:         value_object.MustNewPrice(50.0),
				TotalPrice:        value_object.MustNewPrice(100.0),
			},
			expectError: true,
			expectedErr: errors.ErrInvalidInput,
		},
		{
			name: "Invalid order - empty platform product id",
			order: &entity.InputOrder{
				No:                1,
				PlatformProductId: "",
				Qty:               2,
				UnitPrice:         value_object.MustNewPrice(50.0),
				TotalPrice:        value_object.MustNewPrice(100.0),
			},
			expectError: true,
			expectedErr: errors.ErrInvalidInput,
		},
		{
			name: "Invalid order - zero quantity",
			order: &entity.InputOrder{
				No:                1,
				PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
				Qty:               0,
				UnitPrice:         value_object.MustNewPrice(50.0),
				TotalPrice:        value_object.MustNewPrice(100.0),
			},
			expectError: true,
			expectedErr: errors.ErrInvalidInput,
		},
		{
			name: "Invalid order - negative quantity",
			order: &entity.InputOrder{
				No:                1,
				PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
				Qty:               -1,
				UnitPrice:         value_object.MustNewPrice(50.0),
				TotalPrice:        value_object.MustNewPrice(100.0),
			},
			expectError: true,
			expectedErr: errors.ErrInvalidInput,
		},
		{
			name: "Invalid order - nil unit price",
			order: &entity.InputOrder{
				No:                1,
				PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
				Qty:               2,
				UnitPrice:         nil,
				TotalPrice:        value_object.MustNewPrice(100.0),
			},
			expectError: true,
			expectedErr: errors.ErrInvalidInput,
		},
		{
			name: "Invalid order - nil total price",
			order: &entity.InputOrder{
				No:                1,
				PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
				Qty:               2,
				UnitPrice:         value_object.MustNewPrice(50.0),
				TotalPrice:        nil,
			},
			expectError: true,
			expectedErr: errors.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.order.IsValid()

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.Equal(t, tt.expectedErr, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCleanedOrder_IsValid(t *testing.T) {
	tests := []struct {
		name        string
		order       *entity.CleanedOrder
		expectError bool
		expectedErr error
	}{
		{
			name: "Valid cleaned order - main product",
			order: &entity.CleanedOrder{
				No:         1,
				ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
				MaterialId: "FG0A-CLEAR",
				ModelId:    "IPHONE16PROMAX",
				Qty:        2,
				UnitPrice:  value_object.MustNewPrice(50.0),
				TotalPrice: value_object.MustNewPrice(100.0),
			},
			expectError: false,
		},
		{
			name: "Valid cleaned order - complementary product",
			order: &entity.CleanedOrder{
				No:         2,
				ProductId:  "WIPING-CLOTH",
				MaterialId: "",
				ModelId:    "",
				Qty:        2,
				UnitPrice:  value_object.MustNewPrice(0.0),
				TotalPrice: value_object.MustNewPrice(0.0),
			},
			expectError: false,
		},
		{
			name: "Valid cleaned order - cleaner product",
			order: &entity.CleanedOrder{
				No:         3,
				ProductId:  "CLEAR-CLEANNER",
				MaterialId: "",
				ModelId:    "",
				Qty:        2,
				UnitPrice:  value_object.MustNewPrice(0.0),
				TotalPrice: value_object.MustNewPrice(0.0),
			},
			expectError: false,
		},
		{
			name: "Invalid order - zero order number",
			order: &entity.CleanedOrder{
				No:         0,
				ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
				MaterialId: "FG0A-CLEAR",
				ModelId:    "IPHONE16PROMAX",
				Qty:        2,
				UnitPrice:  value_object.MustNewPrice(50.0),
				TotalPrice: value_object.MustNewPrice(100.0),
			},
			expectError: true,
			expectedErr: errors.ErrInvalidInput,
		},
		{
			name: "Invalid order - negative order number",
			order: &entity.CleanedOrder{
				No:         -1,
				ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
				MaterialId: "FG0A-CLEAR",
				ModelId:    "IPHONE16PROMAX",
				Qty:        2,
				UnitPrice:  value_object.MustNewPrice(50.0),
				TotalPrice: value_object.MustNewPrice(100.0),
			},
			expectError: true,
			expectedErr: errors.ErrInvalidInput,
		},
		{
			name: "Invalid order - empty product id",
			order: &entity.CleanedOrder{
				No:         1,
				ProductId:  "",
				MaterialId: "FG0A-CLEAR",
				ModelId:    "IPHONE16PROMAX",
				Qty:        2,
				UnitPrice:  value_object.MustNewPrice(50.0),
				TotalPrice: value_object.MustNewPrice(100.0),
			},
			expectError: true,
			expectedErr: errors.ErrInvalidInput,
		},
		{
			name: "Invalid order - zero quantity",
			order: &entity.CleanedOrder{
				No:         1,
				ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
				MaterialId: "FG0A-CLEAR",
				ModelId:    "IPHONE16PROMAX",
				Qty:        0,
				UnitPrice:  value_object.MustNewPrice(50.0),
				TotalPrice: value_object.MustNewPrice(100.0),
			},
			expectError: true,
			expectedErr: errors.ErrInvalidInput,
		},
		{
			name: "Invalid order - negative quantity",
			order: &entity.CleanedOrder{
				No:         1,
				ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
				MaterialId: "FG0A-CLEAR",
				ModelId:    "IPHONE16PROMAX",
				Qty:        -1,
				UnitPrice:  value_object.MustNewPrice(50.0),
				TotalPrice: value_object.MustNewPrice(100.0),
			},
			expectError: true,
			expectedErr: errors.ErrInvalidInput,
		},
		{
			name: "Invalid order - nil unit price",
			order: &entity.CleanedOrder{
				No:         1,
				ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
				MaterialId: "FG0A-CLEAR",
				ModelId:    "IPHONE16PROMAX",
				Qty:        2,
				UnitPrice:  nil,
				TotalPrice: value_object.MustNewPrice(100.0),
			},
			expectError: true,
			expectedErr: errors.ErrInvalidInput,
		},
		{
			name: "Invalid order - nil total price",
			order: &entity.CleanedOrder{
				No:         1,
				ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
				MaterialId: "FG0A-CLEAR",
				ModelId:    "IPHONE16PROMAX",
				Qty:        2,
				UnitPrice:  value_object.MustNewPrice(50.0),
				TotalPrice: nil,
			},
			expectError: true,
			expectedErr: errors.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.order.IsValid()

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.Equal(t, tt.expectedErr, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewOrderBatch(t *testing.T) {
	tests := []struct {
		name     string
		orders   []entity.InputOrder
		expected *entity.OrderBatch
	}{
		{
			name:   "Empty orders",
			orders: []entity.InputOrder{},
			expected: &entity.OrderBatch{
				Orders: []entity.InputOrder{},
			},
		},
		{
			name: "Single order",
			orders: []entity.InputOrder{
				{
					No:                1,
					PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
					Qty:               2,
					UnitPrice:         value_object.MustNewPrice(50.0),
					TotalPrice:        value_object.MustNewPrice(100.0),
				},
			},
			expected: &entity.OrderBatch{
				Orders: []entity.InputOrder{
					{
						No:                1,
						PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
						Qty:               2,
						UnitPrice:         value_object.MustNewPrice(50.0),
						TotalPrice:        value_object.MustNewPrice(100.0),
					},
				},
			},
		},
		{
			name: "Multiple orders",
			orders: []entity.InputOrder{
				{
					No:                1,
					PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
					Qty:               2,
					UnitPrice:         value_object.MustNewPrice(50.0),
					TotalPrice:        value_object.MustNewPrice(100.0),
				},
				{
					No:                2,
					PlatformProductId: "FG05-MATTE-OPPOA3",
					Qty:               1,
					UnitPrice:         value_object.MustNewPrice(40.0),
					TotalPrice:        value_object.MustNewPrice(40.0),
				},
			},
			expected: &entity.OrderBatch{
				Orders: []entity.InputOrder{
					{
						No:                1,
						PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
						Qty:               2,
						UnitPrice:         value_object.MustNewPrice(50.0),
						TotalPrice:        value_object.MustNewPrice(100.0),
					},
					{
						No:                2,
						PlatformProductId: "FG05-MATTE-OPPOA3",
						Qty:               1,
						UnitPrice:         value_object.MustNewPrice(40.0),
						TotalPrice:        value_object.MustNewPrice(40.0),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			batch := entity.NewOrderBatch(tt.orders)

			assert.NotNil(t, batch)
			assert.Equal(t, len(tt.expected.Orders), len(batch.Orders))

			for i, expectedOrder := range tt.expected.Orders {
				assert.Equal(t, expectedOrder.No, batch.Orders[i].No)
				assert.Equal(t, expectedOrder.PlatformProductId, batch.Orders[i].PlatformProductId)
				assert.Equal(t, expectedOrder.Qty, batch.Orders[i].Qty)
				assert.Equal(t, expectedOrder.UnitPrice.Amount(), batch.Orders[i].UnitPrice.Amount())
				assert.Equal(t, expectedOrder.TotalPrice.Amount(), batch.Orders[i].TotalPrice.Amount())
			}
		})
	}
}

func TestInputOrder_EdgeCases(t *testing.T) {
	t.Run("Large quantity", func(t *testing.T) {
		order := &entity.InputOrder{
			No:                1,
			PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
			Qty:               1000000,
			UnitPrice:         value_object.MustNewPrice(0.001),
			TotalPrice:        value_object.MustNewPrice(1000.0),
		}

		err := order.IsValid()
		assert.NoError(t, err)
	})

	t.Run("Very small price", func(t *testing.T) {
		order := &entity.InputOrder{
			No:                1,
			PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
			Qty:               1,
			UnitPrice:         value_object.MustNewPrice(0.01),
			TotalPrice:        value_object.MustNewPrice(0.01),
		}

		err := order.IsValid()
		assert.NoError(t, err)
	})

	t.Run("Complex platform product ID", func(t *testing.T) {
		order := &entity.InputOrder{
			No:                1,
			PlatformProductId: "x2-3&--FG0A-CLEAR-IPHONE16PROMAX*2/FG05-MATTE-OPPOA3",
			Qty:               1,
			UnitPrice:         value_object.MustNewPrice(100.0),
			TotalPrice:        value_object.MustNewPrice(100.0),
		}

		err := order.IsValid()
		assert.NoError(t, err)
	})
}

func TestCleanedOrder_EdgeCases(t *testing.T) {
	t.Run("Product with complex model ID", func(t *testing.T) {
		order := &entity.CleanedOrder{
			No:         1,
			ProductId:  "FG0A-CLEAR-OPPOA3-B-SPECIAL-EDITION",
			MaterialId: "FG0A-CLEAR",
			ModelId:    "OPPOA3-B-SPECIAL-EDITION",
			Qty:        1,
			UnitPrice:  value_object.MustNewPrice(75.0),
			TotalPrice: value_object.MustNewPrice(75.0),
		}

		err := order.IsValid()
		assert.NoError(t, err)
	})

	t.Run("Cleaner product with specific texture", func(t *testing.T) {
		order := &entity.CleanedOrder{
			No:         2,
			ProductId:  "PRIVACY-CLEANNER",
			MaterialId: "",
			ModelId:    "",
			Qty:        1,
			UnitPrice:  value_object.MustNewPrice(0.0),
			TotalPrice: value_object.MustNewPrice(0.0),
		}

		err := order.IsValid()
		assert.NoError(t, err)
	})
}

// Benchmark tests สำหรับ performance testing
func BenchmarkInputOrder_IsValid(b *testing.B) {
	order := &entity.InputOrder{
		No:                1,
		PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
		Qty:               2,
		UnitPrice:         value_object.MustNewPrice(50.0),
		TotalPrice:        value_object.MustNewPrice(100.0),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = order.IsValid()
	}
}

func BenchmarkCleanedOrder_IsValid(b *testing.B) {
	order := &entity.CleanedOrder{
		No:         1,
		ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
		MaterialId: "FG0A-CLEAR",
		ModelId:    "IPHONE16PROMAX",
		Qty:        2,
		UnitPrice:  value_object.MustNewPrice(50.0),
		TotalPrice: value_object.MustNewPrice(100.0),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = order.IsValid()
	}
}

func BenchmarkNewOrderBatch(b *testing.B) {
	orders := []entity.InputOrder{
		{
			No:                1,
			PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
			Qty:               2,
			UnitPrice:         value_object.MustNewPrice(50.0),
			TotalPrice:        value_object.MustNewPrice(100.0),
		},
		{
			No:                2,
			PlatformProductId: "FG05-MATTE-OPPOA3",
			Qty:               1,
			UnitPrice:         value_object.MustNewPrice(40.0),
			TotalPrice:        value_object.MustNewPrice(40.0),
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = entity.NewOrderBatch(orders)
	}
}

// Test helper functions
func createValidInputOrder(t *testing.T) *entity.InputOrder {
	t.Helper()

	return &entity.InputOrder{
		No:                1,
		PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
		Qty:               2,
		UnitPrice:         value_object.MustNewPrice(50.0),
		TotalPrice:        value_object.MustNewPrice(100.0),
	}
}

func createValidCleanedOrder(t *testing.T) *entity.CleanedOrder {
	t.Helper()

	return &entity.CleanedOrder{
		No:         1,
		ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
		MaterialId: "FG0A-CLEAR",
		ModelId:    "IPHONE16PROMAX",
		Qty:        2,
		UnitPrice:  value_object.MustNewPrice(50.0),
		TotalPrice: value_object.MustNewPrice(100.0),
	}
}

// Integration test with real scenario
func TestOrderValidation_RealScenario(t *testing.T) {
	t.Run("Case 1: Single product scenario", func(t *testing.T) {
		inputOrder := &entity.InputOrder{
			No:                1,
			PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
			Qty:               2,
			UnitPrice:         value_object.MustNewPrice(50.0),
			TotalPrice:        value_object.MustNewPrice(100.0),
		}

		err := inputOrder.IsValid()
		assert.NoError(t, err)

		expectedCleanedOrders := []*entity.CleanedOrder{
			{
				No:         1,
				ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
				MaterialId: "FG0A-CLEAR",
				ModelId:    "IPHONE16PROMAX",
				Qty:        2,
				UnitPrice:  value_object.MustNewPrice(50.0),
				TotalPrice: value_object.MustNewPrice(100.0),
			},
			{
				No:         2,
				ProductId:  "WIPING-CLOTH",
				MaterialId: "",
				ModelId:    "",
				Qty:        2,
				UnitPrice:  value_object.MustNewPrice(0.0),
				TotalPrice: value_object.MustNewPrice(0.0),
			},
			{
				No:         3,
				ProductId:  "CLEAR-CLEANNER",
				MaterialId: "",
				ModelId:    "",
				Qty:        2,
				UnitPrice:  value_object.MustNewPrice(0.0),
				TotalPrice: value_object.MustNewPrice(0.0),
			},
		}

		for _, order := range expectedCleanedOrders {
			err := order.IsValid()
			assert.NoError(t, err, "Order %d should be valid", order.No)
		}
	})

	t.Run("Case 7: Multiple products scenario", func(t *testing.T) {
		inputOrders := []*entity.InputOrder{
			{
				No:                1,
				PlatformProductId: "--FG0A-CLEAR-OPPOA3*2/FG0A-MATTE-OPPOA3*2",
				Qty:               1,
				UnitPrice:         value_object.MustNewPrice(160.0),
				TotalPrice:        value_object.MustNewPrice(160.0),
			},
			{
				No:                2,
				PlatformProductId: "FG0A-PRIVACY-IPHONE16PROMAX",
				Qty:               1,
				UnitPrice:         value_object.MustNewPrice(50.0),
				TotalPrice:        value_object.MustNewPrice(50.0),
			},
		}

		batch := entity.NewOrderBatch([]entity.InputOrder{})
		assert.NotNil(t, batch)

		for _, order := range inputOrders {
			err := order.IsValid()
			assert.NoError(t, err, "Input order %d should be valid", order.No)
		}
	})
}
