package model_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"order-placement-system/internal/adapter/handler/model"
	"order-placement-system/internal/domain/entity"
	"order-placement-system/internal/domain/value_object"
	"order-placement-system/pkg/errors"
	"order-placement-system/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	log.Init("dev")
}

func TestInputOrder_Parse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		requestBody string
		expectError bool
		expected    []*model.InputOrder
	}{
		{
			name: "Valid single order",
			requestBody: `[
				{
					"no": 1,
					"platformProductId": "FG0A-CLEAR-IPHONE16PROMAX",
					"qty": 2,
					"unitPrice": 50.0,
					"totalPrice": 100.0
				}
			]`,
			expectError: false,
			expected: []*model.InputOrder{
				{
					No:                1,
					PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
					Qty:               2,
					UnitPrice:         50.0,
					TotalPrice:        100.0,
				},
			},
		},
		{
			name: "Valid multiple orders",
			requestBody: `[
				{
					"no": 1,
					"platformProductId": "FG0A-CLEAR-IPHONE16PROMAX",
					"qty": 2,
					"unitPrice": 50.0,
					"totalPrice": 100.0
				},
				{
					"no": 2,
					"platformProductId": "FG0A-MATTE-OPPOA3",
					"qty": 1,
					"unitPrice": 40.0,
					"totalPrice": 40.0
				}
			]`,
			expectError: false,
			expected: []*model.InputOrder{
				{
					No:                1,
					PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
					Qty:               2,
					UnitPrice:         50.0,
					TotalPrice:        100.0,
				},
				{
					No:                2,
					PlatformProductId: "FG0A-MATTE-OPPOA3",
					Qty:               1,
					UnitPrice:         40.0,
					TotalPrice:        40.0,
				},
			},
		},
		{
			name:        "Empty array",
			requestBody: `[]`,
			expectError: true,
			expected:    nil,
		},
		{
			name:        "Invalid JSON",
			requestBody: `invalid json`,
			expectError: true,
			expected:    nil,
		},
		{
			name: "Missing required fields",
			requestBody: `[
				{
					"no": 1,
					"platformProductId": "",
					"qty": 2
				}
			]`,
			expectError: true,
			expected:    nil,
		},
		{
			name: "Invalid no field (zero)",
			requestBody: `[
				{
					"no": 0,
					"platformProductId": "FG0A-CLEAR-IPHONE16PROMAX",
					"qty": 2,
					"unitPrice": 50.0,
					"totalPrice": 100.0
				}
			]`,
			expectError: true,
			expected:    nil,
		},
		{
			name: "Invalid qty field (zero)",
			requestBody: `[
				{
					"no": 1,
					"platformProductId": "FG0A-CLEAR-IPHONE16PROMAX",
					"qty": 0,
					"unitPrice": 50.0,
					"totalPrice": 100.0
				}
			]`,
			expectError: true,
			expected:    nil,
		},
		{
			name: "Invalid unitPrice field (negative)",
			requestBody: `[
				{
					"no": 1,
					"platformProductId": "FG0A-CLEAR-IPHONE16PROMAX",
					"qty": 2,
					"unitPrice": -50.0,
					"totalPrice": 100.0
				}
			]`,
			expectError: true,
			expected:    nil,
		},
		{
			name: "Invalid totalPrice field (negative)",
			requestBody: `[
				{
					"no": 1,
					"platformProductId": "FG0A-CLEAR-IPHONE16PROMAX",
					"qty": 2,
					"unitPrice": 50.0,
					"totalPrice": -100.0
				}
			]`,
			expectError: true,
			expected:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			inputOrder := &model.InputOrder{}
			result, err := inputOrder.Parse(c)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, len(tt.expected), len(result))

				for i, expected := range tt.expected {
					assert.Equal(t, expected.No, result[i].No)
					assert.Equal(t, expected.PlatformProductId, result[i].PlatformProductId)
					assert.Equal(t, expected.Qty, result[i].Qty)
					assert.Equal(t, expected.UnitPrice, result[i].UnitPrice)
					assert.Equal(t, expected.TotalPrice, result[i].TotalPrice)
				}
			}
		})
	}
}

func TestInputOrder_ToEntity(t *testing.T) {
	tests := []struct {
		name        string
		inputOrder  *model.InputOrder
		expectError bool
		expected    *entity.InputOrder
	}{
		{
			name: "Valid input order",
			inputOrder: &model.InputOrder{
				No:                1,
				PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
				Qty:               2,
				UnitPrice:         50.0,
				TotalPrice:        100.0,
			},
			expectError: false,
			expected: &entity.InputOrder{
				No:                1,
				PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
				Qty:               2,
			},
		},
		{
			name: "Valid input order with zero prices",
			inputOrder: &model.InputOrder{
				No:                1,
				PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
				Qty:               2,
				UnitPrice:         0.0,
				TotalPrice:        0.0,
			},
			expectError: false,
			expected: &entity.InputOrder{
				No:                1,
				PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
				Qty:               2,
			},
		},
		{
			name: "Invalid unit price (negative)",
			inputOrder: &model.InputOrder{
				No:                1,
				PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
				Qty:               2,
				UnitPrice:         -50.0,
				TotalPrice:        100.0,
			},
			expectError: true,
			expected:    nil,
		},
		{
			name: "Invalid total price (negative)",
			inputOrder: &model.InputOrder{
				No:                1,
				PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
				Qty:               2,
				UnitPrice:         50.0,
				TotalPrice:        -100.0,
			},
			expectError: true,
			expected:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity, err := tt.inputOrder.ToEntity()

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, entity)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, entity)
				assert.Equal(t, tt.expected.No, entity.No)
				assert.Equal(t, tt.expected.PlatformProductId, entity.PlatformProductId)
				assert.Equal(t, tt.expected.Qty, entity.Qty)
				assert.Equal(t, tt.inputOrder.UnitPrice, entity.UnitPrice.Amount())
				assert.Equal(t, tt.inputOrder.TotalPrice, entity.TotalPrice.Amount())
			}
		})
	}
}

func TestToEntity(t *testing.T) {
	tests := []struct {
		name        string
		models      []*model.InputOrder
		expectError bool
		expected    []*entity.InputOrder
	}{
		{
			name: "Valid models",
			models: []*model.InputOrder{
				{
					No:                1,
					PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
					Qty:               2,
					UnitPrice:         50.0,
					TotalPrice:        100.0,
				},
				{
					No:                2,
					PlatformProductId: "FG0A-MATTE-OPPOA3",
					Qty:               1,
					UnitPrice:         40.0,
					TotalPrice:        40.0,
				},
			},
			expectError: false,
			expected: []*entity.InputOrder{
				{
					No:                1,
					PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
					Qty:               2,
				},
				{
					No:                2,
					PlatformProductId: "FG0A-MATTE-OPPOA3",
					Qty:               1,
				},
			},
		},
		{
			name:        "Empty models",
			models:      []*model.InputOrder{},
			expectError: false,
			expected:    []*entity.InputOrder{},
		},
		{
			name: "Invalid model in array",
			models: []*model.InputOrder{
				{
					No:                1,
					PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
					Qty:               2,
					UnitPrice:         50.0,
					TotalPrice:        100.0,
				},
				{
					No:                2,
					PlatformProductId: "FG0A-MATTE-OPPOA3",
					Qty:               1,
					UnitPrice:         -40.0, // Invalid negative price
					TotalPrice:        40.0,
				},
			},
			expectError: true,
			expected:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entities, err := model.ToEntity(tt.models)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, entities)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, entities)
				assert.Equal(t, len(tt.expected), len(entities))

				for i, expected := range tt.expected {
					assert.Equal(t, expected.No, entities[i].No)
					assert.Equal(t, expected.PlatformProductId, entities[i].PlatformProductId)
					assert.Equal(t, expected.Qty, entities[i].Qty)
					assert.Equal(t, tt.models[i].UnitPrice, entities[i].UnitPrice.Amount())
					assert.Equal(t, tt.models[i].TotalPrice, entities[i].TotalPrice.Amount())
				}
			}
		})
	}
}

func TestFromEntity(t *testing.T) {
	tests := []struct {
		name     string
		entity   *entity.CleanedOrder
		expected *model.CleanedOrder
	}{
		{
			name: "Valid main product entity",
			entity: &entity.CleanedOrder{
				No:         1,
				ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
				MaterialId: "FG0A-CLEAR",
				ModelId:    "IPHONE16PROMAX",
				Qty:        2,
				UnitPrice:  value_object.MustNewPrice(50.0),
				TotalPrice: value_object.MustNewPrice(100.0),
			},
			expected: &model.CleanedOrder{
				No:         1,
				ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
				MaterialId: "FG0A-CLEAR",
				ModelId:    "IPHONE16PROMAX",
				Qty:        2,
				UnitPrice:  value_object.MustNewPrice(50.0),
				TotalPrice: value_object.MustNewPrice(100.0),
			},
		},
		{
			name: "Valid complementary product entity",
			entity: &entity.CleanedOrder{
				No:         2,
				ProductId:  "WIPING-CLOTH",
				MaterialId: "",
				ModelId:    "",
				Qty:        2,
				UnitPrice:  value_object.ZeroPrice(),
				TotalPrice: value_object.ZeroPrice(),
			},
			expected: &model.CleanedOrder{
				No:         2,
				ProductId:  "WIPING-CLOTH",
				MaterialId: "",
				ModelId:    "",
				Qty:        2,
				UnitPrice:  value_object.ZeroPrice(),
				TotalPrice: value_object.ZeroPrice(),
			},
		},
		{
			name: "Valid cleaner product entity",
			entity: &entity.CleanedOrder{
				No:         3,
				ProductId:  "CLEAR-CLEANNER",
				MaterialId: "",
				ModelId:    "",
				Qty:        2,
				UnitPrice:  value_object.ZeroPrice(),
				TotalPrice: value_object.ZeroPrice(),
			},
			expected: &model.CleanedOrder{
				No:         3,
				ProductId:  "CLEAR-CLEANNER",
				MaterialId: "",
				ModelId:    "",
				Qty:        2,
				UnitPrice:  value_object.ZeroPrice(),
				TotalPrice: value_object.ZeroPrice(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := model.FromEntity(tt.entity)

			assert.NotNil(t, result)
			assert.Equal(t, tt.expected.No, result.No)
			assert.Equal(t, tt.expected.ProductId, result.ProductId)
			assert.Equal(t, tt.expected.MaterialId, result.MaterialId)
			assert.Equal(t, tt.expected.ModelId, result.ModelId)
			assert.Equal(t, tt.expected.Qty, result.Qty)
			assert.Equal(t, tt.expected.UnitPrice.Amount(), result.UnitPrice.Amount())
			assert.Equal(t, tt.expected.TotalPrice.Amount(), result.TotalPrice.Amount())
		})
	}
}

func TestFromEntities(t *testing.T) {
	tests := []struct {
		name     string
		entities []*entity.CleanedOrder
		expected []*model.CleanedOrder
	}{
		{
			name: "Valid entities",
			entities: []*entity.CleanedOrder{
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
					UnitPrice:  value_object.ZeroPrice(),
					TotalPrice: value_object.ZeroPrice(),
				},
			},
			expected: []*model.CleanedOrder{
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
					UnitPrice:  value_object.ZeroPrice(),
					TotalPrice: value_object.ZeroPrice(),
				},
			},
		},
		{
			name:     "Empty entities",
			entities: []*entity.CleanedOrder{},
			expected: []*model.CleanedOrder{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := model.FromEntities(tt.entities)

			assert.NotNil(t, result)
			assert.Equal(t, len(tt.expected), len(result))

			for i, expected := range tt.expected {
				assert.Equal(t, expected.No, result[i].No)
				assert.Equal(t, expected.ProductId, result[i].ProductId)
				assert.Equal(t, expected.MaterialId, result[i].MaterialId)
				assert.Equal(t, expected.ModelId, result[i].ModelId)
				assert.Equal(t, expected.Qty, result[i].Qty)
				assert.Equal(t, expected.UnitPrice.Amount(), result[i].UnitPrice.Amount())
				assert.Equal(t, expected.TotalPrice.Amount(), result[i].TotalPrice.Amount())
			}
		})
	}
}

func TestInputOrder_Validate(t *testing.T) {
	tests := []struct {
		name        string
		inputOrder  *model.InputOrder
		expectError bool
	}{
		{
			name: "Valid input order",
			inputOrder: &model.InputOrder{
				No:                1,
				PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
				Qty:               2,
				UnitPrice:         50.0,
				TotalPrice:        100.0,
			},
			expectError: false,
		},
		{
			name: "Valid input order with zero prices",
			inputOrder: &model.InputOrder{
				No:                1,
				PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
				Qty:               2,
				UnitPrice:         0.0,
				TotalPrice:        0.0,
			},
			expectError: false,
		},
		{
			name: "Invalid No (zero)",
			inputOrder: &model.InputOrder{
				No:                0,
				PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
				Qty:               2,
				UnitPrice:         50.0,
				TotalPrice:        100.0,
			},
			expectError: true,
		},
		{
			name: "Invalid No (negative)",
			inputOrder: &model.InputOrder{
				No:                -1,
				PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
				Qty:               2,
				UnitPrice:         50.0,
				TotalPrice:        100.0,
			},
			expectError: true,
		},
		{
			name: "Invalid PlatformProductId (empty)",
			inputOrder: &model.InputOrder{
				No:                1,
				PlatformProductId: "",
				Qty:               2,
				UnitPrice:         50.0,
				TotalPrice:        100.0,
			},
			expectError: true,
		},
		{
			name: "Invalid Qty (zero)",
			inputOrder: &model.InputOrder{
				No:                1,
				PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
				Qty:               0,
				UnitPrice:         50.0,
				TotalPrice:        100.0,
			},
			expectError: true,
		},
		{
			name: "Invalid Qty (negative)",
			inputOrder: &model.InputOrder{
				No:                1,
				PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
				Qty:               -2,
				UnitPrice:         50.0,
				TotalPrice:        100.0,
			},
			expectError: true,
		},
		{
			name: "Invalid UnitPrice (negative)",
			inputOrder: &model.InputOrder{
				No:                1,
				PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
				Qty:               2,
				UnitPrice:         -50.0,
				TotalPrice:        100.0,
			},
			expectError: true,
		},
		{
			name: "Invalid TotalPrice (negative)",
			inputOrder: &model.InputOrder{
				No:                1,
				PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
				Qty:               2,
				UnitPrice:         50.0,
				TotalPrice:        -100.0,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.inputOrder.Validate()

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, errors.ErrInvalidInput, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Benchmark tests
func BenchmarkInputOrder_ToEntity(b *testing.B) {
	inputOrder := &model.InputOrder{
		No:                1,
		PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
		Qty:               2,
		UnitPrice:         50.0,
		TotalPrice:        100.0,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := inputOrder.ToEntity()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFromEntity(b *testing.B) {
	entity := &entity.CleanedOrder{
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
		_ = model.FromEntity(entity)
	}
}

func BenchmarkInputOrder_Validate(b *testing.B) {
	inputOrder := &model.InputOrder{
		No:                1,
		PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
		Qty:               2,
		UnitPrice:         50.0,
		TotalPrice:        100.0,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = inputOrder.Validate()
	}
}

// Helper functions for testing
func createTestInputOrder() *model.InputOrder {
	return &model.InputOrder{
		No:                1,
		PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
		Qty:               2,
		UnitPrice:         50.0,
		TotalPrice:        100.0,
	}
}

func createTestCleanedOrderEntity() *entity.CleanedOrder {
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

// Table-driven tests for comprehensive coverage
func TestInputOrder_Parse_TableDriven(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		input          string
		expectedLength int
		expectError    bool
	}{
		{
			name:           "Single valid order",
			input:          `[{"no":1,"platformProductId":"FG0A-CLEAR-IPHONE16PROMAX","qty":2,"unitPrice":50.0,"totalPrice":100.0}]`,
			expectedLength: 1,
			expectError:    false,
		},
		{
			name:           "Multiple valid orders",
			input:          `[{"no":1,"platformProductId":"FG0A-CLEAR-IPHONE16PROMAX","qty":2,"unitPrice":50.0,"totalPrice":100.0},{"no":2,"platformProductId":"FG0A-MATTE-OPPOA3","qty":1,"unitPrice":40.0,"totalPrice":40.0}]`,
			expectedLength: 2,
			expectError:    false,
		},
		{
			name:           "Empty array",
			input:          `[]`,
			expectedLength: 0,
			expectError:    true,
		},
		{
			name:           "Invalid JSON structure",
			input:          `{"invalid": "structure"}`,
			expectedLength: 0,
			expectError:    true,
		},
		{
			name:           "Malformed JSON",
			input:          `[{"no":1,"platformProductId":"FG0A-CLEAR-IPHONE16PROMAX","qty":2,"unitPrice":50.0,"totalPrice":100.0`,
			expectedLength: 0,
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tc.input))
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			inputOrder := &model.InputOrder{}
			result, err := inputOrder.Parse(c)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.expectedLength, len(result))
			}
		})
	}
}

// Edge case tests
func TestInputOrder_EdgeCases(t *testing.T) {
	t.Run("Very large numbers", func(t *testing.T) {
		inputOrder := &model.InputOrder{
			No:                1,
			PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
			Qty:               1000000,
			UnitPrice:         9999999.99,
			TotalPrice:        9999999990000.0,
		}

		entity, err := inputOrder.ToEntity()
		assert.NoError(t, err)
		assert.NotNil(t, entity)
		assert.Equal(t, 1000000, entity.Qty)
	})

	t.Run("Very small positive numbers", func(t *testing.T) {
		inputOrder := &model.InputOrder{
			No:                1,
			PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
			Qty:               1,
			UnitPrice:         0.01,
			TotalPrice:        0.01,
		}

		entity, err := inputOrder.ToEntity()
		assert.NoError(t, err)
		assert.NotNil(t, entity)
		assert.Equal(t, 0.01, entity.UnitPrice.Amount())
	})

	t.Run("Complex product ID", func(t *testing.T) {
		inputOrder := &model.InputOrder{
			No:                1,
			PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX-SPECIAL-EDITION-LIMITED",
			Qty:               1,
			UnitPrice:         100.0,
			TotalPrice:        100.0,
		}

		entity, err := inputOrder.ToEntity()
		assert.NoError(t, err)
		assert.NotNil(t, entity)
		assert.Equal(t, "FG0A-CLEAR-IPHONE16PROMAX-SPECIAL-EDITION-LIMITED", entity.PlatformProductId)
	})
}

// Integration-style tests
func TestInputOrder_Integration(t *testing.T) {
	t.Run("Full flow: Parse -> ToEntity -> FromEntity", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		// Step 1: Parse JSON
		requestBody := `[{"no":1,"platformProductId":"FG0A-CLEAR-IPHONE16PROMAX","qty":2,"unitPrice":50.0,"totalPrice":100.0}]`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		inputOrder := &model.InputOrder{}
		parsedOrders, err := inputOrder.Parse(c)
		require.NoError(t, err)
		require.Len(t, parsedOrders, 1)

		// Step 2: Convert to entities
		entities, err := model.ToEntity(parsedOrders)
		require.NoError(t, err)
		require.Len(t, entities, 1)

		// Step 3: Create a cleaned order entity (simulating processing)
		cleanedEntity := &entity.CleanedOrder{
			No:         1,
			ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
			MaterialId: "FG0A-CLEAR",
			ModelId:    "IPHONE16PROMAX",
			Qty:        2,
			UnitPrice:  value_object.MustNewPrice(50.0),
			TotalPrice: value_object.MustNewPrice(100.0),
		}

		// Step 4: Convert back to model
		cleanedModel := model.FromEntity(cleanedEntity)
		require.NotNil(t, cleanedModel)

		// Verify the full cycle
		assert.Equal(t, 1, cleanedModel.No)
		assert.Equal(t, "FG0A-CLEAR-IPHONE16PROMAX", cleanedModel.ProductId)
		assert.Equal(t, "FG0A-CLEAR", cleanedModel.MaterialId)
		assert.Equal(t, "IPHONE16PROMAX", cleanedModel.ModelId)
		assert.Equal(t, 2, cleanedModel.Qty)
		assert.Equal(t, 50.0, cleanedModel.UnitPrice.Amount())
		assert.Equal(t, 100.0, cleanedModel.TotalPrice.Amount())
	})
}
