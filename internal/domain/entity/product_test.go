package entity_test

import (
	"encoding/json"
	"testing"
	"time"

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

func TestNewTexture(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    value_object.Texture
		expectError bool
	}{
		{
			name:        "Valid CLEAR texture",
			input:       "CLEAR",
			expected:    value_object.TextureClear,
			expectError: false,
		},
		{
			name:        "Valid MATTE texture",
			input:       "MATTE",
			expected:    value_object.TextureMatte,
			expectError: false,
		},
		{
			name:        "Valid PRIVACY texture",
			input:       "PRIVACY",
			expected:    value_object.TexturePrivacy,
			expectError: false,
		},
		{
			name:        "Valid lowercase clear",
			input:       "clear",
			expected:    value_object.TextureClear,
			expectError: false,
		},
		{
			name:        "Valid mixed case matte",
			input:       "MaTtE",
			expected:    value_object.TextureMatte,
			expectError: false,
		},
		{
			name:        "Valid with whitespace",
			input:       " CLEAR ",
			expected:    value_object.TextureClear,
			expectError: false,
		},
		{
			name:        "Invalid texture",
			input:       "INVALID",
			expectError: true,
		},
		{
			name:        "Empty string",
			input:       "",
			expectError: true,
		},
		{
			name:        "Only whitespace",
			input:       "   ",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := value_object.NewTexture(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, errors.ErrInvalidInput, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestTexture_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		texture  value_object.Texture
		expected bool
	}{
		{
			name:     "Valid CLEAR texture",
			texture:  value_object.TextureClear,
			expected: true,
		},
		{
			name:     "Valid MATTE texture",
			texture:  value_object.TextureMatte,
			expected: true,
		},
		{
			name:     "Valid PRIVACY texture",
			texture:  value_object.TexturePrivacy,
			expected: true,
		},
		{
			name:     "Invalid texture",
			texture:  value_object.Texture("INVALID"),
			expected: false,
		},
		{
			name:     "Empty texture",
			texture:  value_object.Texture(""),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.texture.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTexture_String(t *testing.T) {
	tests := []struct {
		name     string
		texture  value_object.Texture
		expected string
	}{
		{
			name:     "CLEAR texture to string",
			texture:  value_object.TextureClear,
			expected: "CLEAR",
		},
		{
			name:     "MATTE texture to string",
			texture:  value_object.TextureMatte,
			expected: "MATTE",
		},
		{
			name:     "PRIVACY texture to string",
			texture:  value_object.TexturePrivacy,
			expected: "PRIVACY",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.texture.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTexture_GetCleanerProductId(t *testing.T) {
	tests := []struct {
		name     string
		texture  value_object.Texture
		expected string
	}{
		{
			name:     "CLEAR cleaner product ID",
			texture:  value_object.TextureClear,
			expected: "CLEAR-CLEANNER",
		},
		{
			name:     "MATTE cleaner product ID",
			texture:  value_object.TextureMatte,
			expected: "MATTE-CLEANNER",
		},
		{
			name:     "PRIVACY cleaner product ID",
			texture:  value_object.TexturePrivacy,
			expected: "PRIVACY-CLEANNER",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.texture.GetCleanerProductId()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTexture_Equals(t *testing.T) {
	tests := []struct {
		name     string
		texture1 value_object.Texture
		texture2 value_object.Texture
		expected bool
	}{
		{
			name:     "Same CLEAR textures",
			texture1: value_object.TextureClear,
			texture2: value_object.TextureClear,
			expected: true,
		},
		{
			name:     "Different textures",
			texture1: value_object.TextureClear,
			texture2: value_object.TextureMatte,
			expected: false,
		},
		{
			name:     "Same MATTE textures",
			texture1: value_object.TextureMatte,
			texture2: value_object.TextureMatte,
			expected: true,
		},
		{
			name:     "Same PRIVACY textures",
			texture1: value_object.TexturePrivacy,
			texture2: value_object.TexturePrivacy,
			expected: true,
		},
		{
			name:     "CLEAR vs PRIVACY",
			texture1: value_object.TextureClear,
			texture2: value_object.TexturePrivacy,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.texture1.Equals(tt.texture2)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTexture_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		texture  value_object.Texture
		expected string
	}{
		{
			name:     "Marshal CLEAR texture",
			texture:  value_object.TextureClear,
			expected: `"CLEAR"`,
		},
		{
			name:     "Marshal MATTE texture",
			texture:  value_object.TextureMatte,
			expected: `"MATTE"`,
		},
		{
			name:     "Marshal PRIVACY texture",
			texture:  value_object.TexturePrivacy,
			expected: `"PRIVACY"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.texture.MarshalJSON()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, string(result))
		})
	}
}

func TestTexture_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    value_object.Texture
		expectError bool
	}{
		{
			name:        "Unmarshal CLEAR texture",
			input:       `"CLEAR"`,
			expected:    value_object.TextureClear,
			expectError: false,
		},
		{
			name:        "Unmarshal MATTE texture",
			input:       `"MATTE"`,
			expected:    value_object.TextureMatte,
			expectError: false,
		},
		{
			name:        "Unmarshal PRIVACY texture",
			input:       `"PRIVACY"`,
			expected:    value_object.TexturePrivacy,
			expectError: false,
		},
		{
			name:        "Unmarshal lowercase texture",
			input:       `"clear"`,
			expected:    value_object.TextureClear,
			expectError: false,
		},
		{
			name:        "Unmarshal invalid texture",
			input:       `"INVALID"`,
			expectError: true,
		},
		{
			name:        "Unmarshal empty string",
			input:       `""`,
			expectError: true,
		},
		{
			name:        "Unmarshal invalid JSON",
			input:       `invalid json`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var texture value_object.Texture
			err := texture.UnmarshalJSON([]byte(tt.input))

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, texture)
			}
		})
	}
}

func TestTexture_IsCompatibleWithFilmType(t *testing.T) {
	tests := []struct {
		name     string
		texture  value_object.Texture
		filmType string
		expected bool
	}{
		{
			name:     "Valid CLEAR texture with any film type",
			texture:  value_object.TextureClear,
			filmType: "FG0A",
			expected: true,
		},
		{
			name:     "Valid MATTE texture with any film type",
			texture:  value_object.TextureMatte,
			filmType: "FG05",
			expected: true,
		},
		{
			name:     "Valid PRIVACY texture with any film type",
			texture:  value_object.TexturePrivacy,
			filmType: "FG1A",
			expected: true,
		},
		{
			name:     "Invalid texture with any film type",
			texture:  value_object.Texture("INVALID"),
			filmType: "FG0A",
			expected: false,
		},
		{
			name:     "Empty texture with any film type",
			texture:  value_object.Texture(""),
			filmType: "FG0A",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.texture.IsCompatibleWithFilmType(tt.filmType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTexture_GetDisplayName(t *testing.T) {
	tests := []struct {
		name     string
		texture  value_object.Texture
		expected string
	}{
		{
			name:     "CLEAR display name",
			texture:  value_object.TextureClear,
			expected: "Clear",
		},
		{
			name:     "MATTE display name",
			texture:  value_object.TextureMatte,
			expected: "Matte",
		},
		{
			name:     "PRIVACY display name",
			texture:  value_object.TexturePrivacy,
			expected: "Privacy",
		},
		{
			name:     "Unknown texture display name",
			texture:  value_object.Texture("UNKNOWN"),
			expected: "UNKNOWN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.texture.GetDisplayName()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTexture_GetPriority(t *testing.T) {
	tests := []struct {
		name     string
		texture  value_object.Texture
		expected int
	}{
		{
			name:     "CLEAR priority",
			texture:  value_object.TextureClear,
			expected: 1,
		},
		{
			name:     "MATTE priority",
			texture:  value_object.TextureMatte,
			expected: 2,
		},
		{
			name:     "PRIVACY priority",
			texture:  value_object.TexturePrivacy,
			expected: 3,
		},
		{
			name:     "Unknown texture priority",
			texture:  value_object.Texture("UNKNOWN"),
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.texture.GetPriority()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseTextureFromMaterialId(t *testing.T) {
	tests := []struct {
		name        string
		materialId  string
		expected    value_object.Texture
		expectError bool
	}{
		{
			name:        "Valid material ID with CLEAR texture",
			materialId:  "FG0A-CLEAR",
			expected:    value_object.TextureClear,
			expectError: false,
		},
		{
			name:        "Valid material ID with MATTE texture",
			materialId:  "FG05-MATTE",
			expected:    value_object.TextureMatte,
			expectError: false,
		},
		{
			name:        "Valid material ID with PRIVACY texture",
			materialId:  "FG1A-PRIVACY",
			expected:    value_object.TexturePrivacy,
			expectError: false,
		},
		{
			name:        "Material ID with lowercase texture",
			materialId:  "FG0A-clear",
			expected:    value_object.TextureClear,
			expectError: false,
		},
		{
			name:        "Material ID with multiple dashes",
			materialId:  "FG0A-CLEAR-EXTRA",
			expected:    value_object.TextureClear,
			expectError: false,
		},
		{
			name:        "Invalid material ID format",
			materialId:  "FG0A",
			expectError: true,
		},
		{
			name:        "Empty material ID",
			materialId:  "",
			expectError: true,
		},
		{
			name:        "Material ID with invalid texture",
			materialId:  "FG0A-INVALID",
			expectError: true,
		},
		{
			name:        "Material ID with no texture part",
			materialId:  "FG0A-",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := value_object.ParseTextureFromMaterialId(tt.materialId)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, errors.ErrInvalidInput, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestTexture_JSONRoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		texture value_object.Texture
	}{
		{
			name:    "CLEAR texture JSON round trip",
			texture: value_object.TextureClear,
		},
		{
			name:    "MATTE texture JSON round trip",
			texture: value_object.TextureMatte,
		},
		{
			name:    "PRIVACY texture JSON round trip",
			texture: value_object.TexturePrivacy,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.texture)
			require.NoError(t, err)

			var unmarshaledTexture value_object.Texture
			err = json.Unmarshal(jsonData, &unmarshaledTexture)
			require.NoError(t, err)

			assert.Equal(t, tt.texture, unmarshaledTexture)
		})
	}
}

func TestAllTextures_Constant(t *testing.T) {
	t.Run("All textures should be valid", func(t *testing.T) {
		for _, texture := range value_object.AllTextures {
			assert.True(t, texture.IsValid(), "Texture %s should be valid", texture.String())
		}
	})

	t.Run("All textures should have unique priorities", func(t *testing.T) {
		priorities := make(map[int]bool)
		for _, texture := range value_object.AllTextures {
			priority := texture.GetPriority()
			assert.False(t, priorities[priority], "Priority %d should be unique", priority)
			priorities[priority] = true
		}
	})

	t.Run("All textures should have different cleaner product IDs", func(t *testing.T) {
		cleanerIds := make(map[string]bool)
		for _, texture := range value_object.AllTextures {
			cleanerId := texture.GetCleanerProductId()
			assert.False(t, cleanerIds[cleanerId], "Cleaner ID %s should be unique", cleanerId)
			cleanerIds[cleanerId] = true
		}
	})
}

func TestTexture_EdgeCases(t *testing.T) {
	t.Run("Case sensitivity", func(t *testing.T) {
		testCases := []string{"clear", "CLEAR", "Clear", "cLeAr"}
		for _, testCase := range testCases {
			texture, err := value_object.NewTexture(testCase)
			assert.NoError(t, err)
			assert.Equal(t, value_object.TextureClear, texture)
		}
	})

	t.Run("Whitespace handling", func(t *testing.T) {
		testCases := []string{" CLEAR ", "\tCLEAR\t", "\nCLEAR\n", "  CLEAR  "}
		for _, testCase := range testCases {
			texture, err := value_object.NewTexture(testCase)
			assert.NoError(t, err)
			assert.Equal(t, value_object.TextureClear, texture)
		}
	})

	t.Run("Empty and whitespace-only strings", func(t *testing.T) {
		testCases := []string{"", " ", "\t", "\n", "   "}
		for _, testCase := range testCases {
			_, err := value_object.NewTexture(testCase)
			assert.Error(t, err)
			assert.Equal(t, errors.ErrInvalidInput, err)
		}
	})
}

func TestProduct_ToCleanedOrder(t *testing.T) {
	tests := []struct {
		name               string
		product            *entity.Product
		orderNo            int
		expectedNo         int
		expectedProductId  string
		expectedMaterialId string
		expectedModelId    string
		expectedQty        int
		expectedUnitPrice  float64
		expectedTotalPrice float64
	}{
		{
			name: "Convert simple product to cleaned order",
			product: &entity.Product{
				ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
				MaterialId: "FG0A-CLEAR",
				ModelId:    "IPHONE16PROMAX",
				Quantity:   2,
				UnitPrice:  value_object.MustNewPrice(50.0),
				TotalPrice: value_object.MustNewPrice(100.0),
			},
			orderNo:            1,
			expectedNo:         1,
			expectedProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
			expectedMaterialId: "FG0A-CLEAR",
			expectedModelId:    "IPHONE16PROMAX",
			expectedQty:        2,
			expectedUnitPrice:  50.0,
			expectedTotalPrice: 100.0,
		},
		{
			name: "Convert product with complex model ID",
			product: &entity.Product{
				ProductId:  "FG05-MATTE-OPPOA3-B",
				MaterialId: "FG05-MATTE",
				ModelId:    "OPPOA3-B",
				Quantity:   1,
				UnitPrice:  value_object.MustNewPrice(40.0),
				TotalPrice: value_object.MustNewPrice(40.0),
			},
			orderNo:            5,
			expectedNo:         5,
			expectedProductId:  "FG05-MATTE-OPPOA3-B",
			expectedMaterialId: "FG05-MATTE",
			expectedModelId:    "OPPOA3-B",
			expectedQty:        1,
			expectedUnitPrice:  40.0,
			expectedTotalPrice: 40.0,
		},
		{
			name: "Convert product with zero prices",
			product: &entity.Product{
				ProductId:  "FG0A-PRIVACY-SAMSUNGS25",
				MaterialId: "FG0A-PRIVACY",
				ModelId:    "SAMSUNGS25",
				Quantity:   3,
				UnitPrice:  value_object.ZeroPrice(),
				TotalPrice: value_object.ZeroPrice(),
			},
			orderNo:            10,
			expectedNo:         10,
			expectedProductId:  "FG0A-PRIVACY-SAMSUNGS25",
			expectedMaterialId: "FG0A-PRIVACY",
			expectedModelId:    "SAMSUNGS25",
			expectedQty:        3,
			expectedUnitPrice:  0.0,
			expectedTotalPrice: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanedOrder := tt.product.ToCleanedOrder(tt.orderNo)

			assert.Equal(t, tt.expectedNo, cleanedOrder.No)
			assert.Equal(t, tt.expectedProductId, cleanedOrder.ProductId)
			assert.Equal(t, tt.expectedMaterialId, cleanedOrder.MaterialId)
			assert.Equal(t, tt.expectedModelId, cleanedOrder.ModelId)
			assert.Equal(t, tt.expectedQty, cleanedOrder.Qty)
			assert.Equal(t, tt.expectedUnitPrice, cleanedOrder.UnitPrice.Amount())
			assert.Equal(t, tt.expectedTotalPrice, cleanedOrder.TotalPrice.Amount())
		})
	}
}

func TestProduct_IsValid(t *testing.T) {
	validUnitPrice := value_object.MustNewPrice(50.0)
	validTotalPrice := value_object.MustNewPrice(100.0)
	zeroPrice := value_object.ZeroPrice()

	tests := []struct {
		name        string
		product     *entity.Product
		expectError bool
		description string
	}{
		{
			name: "Valid product - all fields correct",
			product: &entity.Product{
				ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
				MaterialId: "FG0A-CLEAR",
				ModelId:    "IPHONE16PROMAX",
				Quantity:   2,
				UnitPrice:  validUnitPrice,
				TotalPrice: validTotalPrice,
			},
			expectError: false,
			description: "Product with all valid fields should pass validation",
		},
		{
			name: "Invalid - empty product ID",
			product: &entity.Product{
				ProductId:  "",
				MaterialId: "FG0A-CLEAR",
				ModelId:    "IPHONE16PROMAX",
				Quantity:   2,
				UnitPrice:  validUnitPrice,
				TotalPrice: validTotalPrice,
			},
			expectError: true,
			description: "Product with empty product ID should fail validation",
		},
		{
			name: "Invalid - empty material ID",
			product: &entity.Product{
				ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
				MaterialId: "",
				ModelId:    "IPHONE16PROMAX",
				Quantity:   2,
				UnitPrice:  validUnitPrice,
				TotalPrice: validTotalPrice,
			},
			expectError: true,
			description: "Product with empty material ID should fail validation",
		},
		{
			name: "Invalid - empty model ID",
			product: &entity.Product{
				ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
				MaterialId: "FG0A-CLEAR",
				ModelId:    "",
				Quantity:   2,
				UnitPrice:  validUnitPrice,
				TotalPrice: validTotalPrice,
			},
			expectError: true,
			description: "Product with empty model ID should fail validation",
		},
		{
			name: "Invalid - zero quantity",
			product: &entity.Product{
				ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
				MaterialId: "FG0A-CLEAR",
				ModelId:    "IPHONE16PROMAX",
				Quantity:   0,
				UnitPrice:  validUnitPrice,
				TotalPrice: validTotalPrice,
			},
			expectError: true,
			description: "Product with zero quantity should fail validation",
		},
		{
			name: "Invalid - negative quantity",
			product: &entity.Product{
				ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
				MaterialId: "FG0A-CLEAR",
				ModelId:    "IPHONE16PROMAX",
				Quantity:   -1,
				UnitPrice:  validUnitPrice,
				TotalPrice: validTotalPrice,
			},
			expectError: true,
			description: "Product with negative quantity should fail validation",
		},
		{
			name: "Invalid - nil unit price",
			product: &entity.Product{
				ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
				MaterialId: "FG0A-CLEAR",
				ModelId:    "IPHONE16PROMAX",
				Quantity:   2,
				UnitPrice:  nil,
				TotalPrice: validTotalPrice,
			},
			expectError: true,
			description: "Product with nil unit price should fail validation",
		},
		{
			name: "Invalid - nil total price",
			product: &entity.Product{
				ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
				MaterialId: "FG0A-CLEAR",
				ModelId:    "IPHONE16PROMAX",
				Quantity:   2,
				UnitPrice:  validUnitPrice,
				TotalPrice: nil,
			},
			expectError: true,
			description: "Product with nil total price should fail validation",
		},
		{
			name: "Valid - zero prices are acceptable",
			product: &entity.Product{
				ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
				MaterialId: "FG0A-CLEAR",
				ModelId:    "IPHONE16PROMAX",
				Quantity:   2,
				UnitPrice:  zeroPrice,
				TotalPrice: zeroPrice,
			},
			expectError: false,
			description: "Product with zero prices should pass validation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.product.IsValid()

			if tt.expectError {
				assert.Error(t, err, tt.description)
				assert.Equal(t, errors.ErrInvalidInput, err)
			} else {
				assert.NoError(t, err, tt.description)
			}
		})
	}
}

func TestCleanedOrder_IsMainProduct(t *testing.T) {
	tests := []struct {
		name        string
		order       *entity.CleanedOrder
		expected    bool
		description string
	}{
		{
			name: "Main product with material and model ID",
			order: &entity.CleanedOrder{
				ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
				MaterialId: "FG0A-CLEAR",
				ModelId:    "IPHONE16PROMAX",
			},
			expected:    true,
			description: "Order with both material and model ID should be main product",
		},
		{
			name: "Main product with complex model ID",
			order: &entity.CleanedOrder{
				ProductId:  "FG05-MATTE-OPPOA3-B-SPECIAL",
				MaterialId: "FG05-MATTE",
				ModelId:    "OPPOA3-B-SPECIAL",
			},
			expected:    true,
			description: "Order with complex model ID should be main product",
		},
		{
			name: "Complementary product - wiping cloth",
			order: &entity.CleanedOrder{
				ProductId:  "WIPING-CLOTH",
				MaterialId: "",
				ModelId:    "",
			},
			expected:    false,
			description: "Wiping cloth should not be main product",
		},
		{
			name: "Complementary product - clear cleaner",
			order: &entity.CleanedOrder{
				ProductId:  "CLEAR-CLEANNER",
				MaterialId: "",
				ModelId:    "",
			},
			expected:    false,
			description: "Clear cleaner should not be main product",
		},
		{
			name: "Complementary product - matte cleaner",
			order: &entity.CleanedOrder{
				ProductId:  "MATTE-CLEANNER",
				MaterialId: "",
				ModelId:    "",
			},
			expected:    false,
			description: "Matte cleaner should not be main product",
		},
		{
			name: "Complementary product - privacy cleaner",
			order: &entity.CleanedOrder{
				ProductId:  "PRIVACY-CLEANNER",
				MaterialId: "",
				ModelId:    "",
			},
			expected:    false,
			description: "Privacy cleaner should not be main product",
		},
		{
			name: "Product with empty material ID",
			order: &entity.CleanedOrder{
				ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
				MaterialId: "",
				ModelId:    "IPHONE16PROMAX",
			},
			expected:    false,
			description: "Product with empty material ID should not be main product",
		},
		{
			name: "Product with empty model ID",
			order: &entity.CleanedOrder{
				ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
				MaterialId: "FG0A-CLEAR",
				ModelId:    "",
			},
			expected:    false,
			description: "Product with empty model ID should not be main product",
		},
		{
			name: "Product with both empty material and model ID",
			order: &entity.CleanedOrder{
				ProductId:  "SOME-PRODUCT",
				MaterialId: "",
				ModelId:    "",
			},
			expected:    false,
			description: "Product with both empty material and model ID should not be main product",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.order.IsMainProduct()
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

func TestCleanedOrder_IsComplementaryProduct(t *testing.T) {
	tests := []struct {
		name        string
		order       *entity.CleanedOrder
		expected    bool
		description string
	}{
		{
			name: "Main product with material and model ID",
			order: &entity.CleanedOrder{
				ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
				MaterialId: "FG0A-CLEAR",
				ModelId:    "IPHONE16PROMAX",
			},
			expected:    false,
			description: "Main product should not be complementary",
		},
		{
			name: "Complementary product - wiping cloth",
			order: &entity.CleanedOrder{
				ProductId:  "WIPING-CLOTH",
				MaterialId: "",
				ModelId:    "",
			},
			expected:    true,
			description: "Wiping cloth should be complementary product",
		},
		{
			name: "Complementary product - clear cleaner",
			order: &entity.CleanedOrder{
				ProductId:  "CLEAR-CLEANNER",
				MaterialId: "",
				ModelId:    "",
			},
			expected:    true,
			description: "Clear cleaner should be complementary product",
		},
		{
			name: "Complementary product - matte cleaner",
			order: &entity.CleanedOrder{
				ProductId:  "MATTE-CLEANNER",
				MaterialId: "",
				ModelId:    "",
			},
			expected:    true,
			description: "Matte cleaner should be complementary product",
		},
		{
			name: "Complementary product - privacy cleaner",
			order: &entity.CleanedOrder{
				ProductId:  "PRIVACY-CLEANNER",
				MaterialId: "",
				ModelId:    "",
			},
			expected:    true,
			description: "Privacy cleaner should be complementary product",
		},
		{
			name: "Product with empty material ID only",
			order: &entity.CleanedOrder{
				ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
				MaterialId: "",
				ModelId:    "IPHONE16PROMAX",
			},
			expected:    true,
			description: "Product with empty material ID should be complementary",
		},
		{
			name: "Product with empty model ID only",
			order: &entity.CleanedOrder{
				ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
				MaterialId: "FG0A-CLEAR",
				ModelId:    "",
			},
			expected:    true,
			description: "Product with empty model ID should be complementary",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.order.IsComplementaryProduct()
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

func TestProduct_Clone(t *testing.T) {
	t.Run("Clone creates exact copy", func(t *testing.T) {
		originalUnitPrice := value_object.MustNewPrice(50.0)
		originalTotalPrice := value_object.MustNewPrice(100.0)

		original := &entity.Product{
			ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
			MaterialId: "FG0A-CLEAR",
			ModelId:    "IPHONE16PROMAX",
			Quantity:   2,
			UnitPrice:  originalUnitPrice,
			TotalPrice: originalTotalPrice,
		}

		cloned := original.Clone()

		assert.Equal(t, original.ProductId, cloned.ProductId)
		assert.Equal(t, original.MaterialId, cloned.MaterialId)
		assert.Equal(t, original.ModelId, cloned.ModelId)
		assert.Equal(t, original.Quantity, cloned.Quantity)
		assert.Equal(t, original.UnitPrice.Amount(), cloned.UnitPrice.Amount())
		assert.Equal(t, original.TotalPrice.Amount(), cloned.TotalPrice.Amount())
	})

	t.Run("Clone is independent of original", func(t *testing.T) {
		originalUnitPrice := value_object.MustNewPrice(50.0)
		originalTotalPrice := value_object.MustNewPrice(100.0)

		original := &entity.Product{
			ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
			MaterialId: "FG0A-CLEAR",
			ModelId:    "IPHONE16PROMAX",
			Quantity:   2,
			UnitPrice:  originalUnitPrice,
			TotalPrice: originalTotalPrice,
		}

		cloned := original.Clone()

		cloned.ProductId = "MODIFIED-PRODUCT"
		cloned.MaterialId = "MODIFIED-MATERIAL"
		cloned.ModelId = "MODIFIED-MODEL"
		cloned.Quantity = 5

		assert.Equal(t, "FG0A-CLEAR-IPHONE16PROMAX", original.ProductId)
		assert.Equal(t, "FG0A-CLEAR", original.MaterialId)
		assert.Equal(t, "IPHONE16PROMAX", original.ModelId)
		assert.Equal(t, 2, original.Quantity)
		assert.Equal(t, 50.0, original.UnitPrice.Amount())
		assert.Equal(t, 100.0, original.TotalPrice.Amount())

		assert.Equal(t, "MODIFIED-PRODUCT", cloned.ProductId)
		assert.Equal(t, "MODIFIED-MATERIAL", cloned.MaterialId)
		assert.Equal(t, "MODIFIED-MODEL", cloned.ModelId)
		assert.Equal(t, 5, cloned.Quantity)
	})

	t.Run("Clone with zero prices", func(t *testing.T) {
		original := &entity.Product{
			ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
			MaterialId: "FG0A-CLEAR",
			ModelId:    "IPHONE16PROMAX",
			Quantity:   1,
			UnitPrice:  value_object.ZeroPrice(),
			TotalPrice: value_object.ZeroPrice(),
		}

		cloned := original.Clone()

		assert.Equal(t, original.ProductId, cloned.ProductId)
		assert.Equal(t, original.MaterialId, cloned.MaterialId)
		assert.Equal(t, original.ModelId, cloned.ModelId)
		assert.Equal(t, original.Quantity, cloned.Quantity)
		assert.Equal(t, 0.0, cloned.UnitPrice.Amount())
		assert.Equal(t, 0.0, cloned.TotalPrice.Amount())
	})

	t.Run("Clone with complex model ID", func(t *testing.T) {
		original := &entity.Product{
			ProductId:  "FG05-MATTE-OPPOA3-B-SPECIAL-EDITION",
			MaterialId: "FG05-MATTE",
			ModelId:    "OPPOA3-B-SPECIAL-EDITION",
			Quantity:   3,
			UnitPrice:  value_object.MustNewPrice(75.0),
			TotalPrice: value_object.MustNewPrice(225.0),
		}

		cloned := original.Clone()

		assert.Equal(t, original.ProductId, cloned.ProductId)
		assert.Equal(t, original.MaterialId, cloned.MaterialId)
		assert.Equal(t, original.ModelId, cloned.ModelId)
		assert.Equal(t, original.Quantity, cloned.Quantity)
		assert.Equal(t, original.UnitPrice.Amount(), cloned.UnitPrice.Amount())
		assert.Equal(t, original.TotalPrice.Amount(), cloned.TotalPrice.Amount())
	})
}

func BenchmarkProduct_ToCleanedOrder(b *testing.B) {
	unitPrice := value_object.MustNewPrice(50.0)
	totalPrice := value_object.MustNewPrice(100.0)

	product := &entity.Product{
		ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
		MaterialId: "FG0A-CLEAR",
		ModelId:    "IPHONE16PROMAX",
		Quantity:   2,
		UnitPrice:  unitPrice,
		TotalPrice: totalPrice,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = product.ToCleanedOrder(1)
	}
}

func BenchmarkProduct_IsValid(b *testing.B) {
	unitPrice := value_object.MustNewPrice(50.0)
	totalPrice := value_object.MustNewPrice(100.0)

	product := &entity.Product{
		ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
		MaterialId: "FG0A-CLEAR",
		ModelId:    "IPHONE16PROMAX",
		Quantity:   2,
		UnitPrice:  unitPrice,
		TotalPrice: totalPrice,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = product.IsValid()
	}
}

func BenchmarkCleanedOrder_IsMainProduct(b *testing.B) {
	order := &entity.CleanedOrder{
		ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
		MaterialId: "FG0A-CLEAR",
		ModelId:    "IPHONE16PROMAX",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = order.IsMainProduct()
	}
}

func BenchmarkCleanedOrder_IsComplementaryProduct(b *testing.B) {
	order := &entity.CleanedOrder{
		ProductId:  "WIPING-CLOTH",
		MaterialId: "",
		ModelId:    "",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = order.IsComplementaryProduct()
	}
}

func BenchmarkProduct_Clone(b *testing.B) {
	unitPrice := value_object.MustNewPrice(50.0)
	totalPrice := value_object.MustNewPrice(100.0)

	product := &entity.Product{
		ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
		MaterialId: "FG0A-CLEAR",
		ModelId:    "IPHONE16PROMAX",
		Quantity:   2,
		UnitPrice:  unitPrice,
		TotalPrice: totalPrice,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = product.Clone()
	}
}

func TestProduct_EdgeCases(t *testing.T) {
	t.Run("Product with very large quantity", func(t *testing.T) {
		unitPrice := value_object.MustNewPrice(0.01)
		totalPrice := value_object.MustNewPrice(1000000.0)

		product, err := entity.NewProduct("FG0A-CLEAR-IPHONE16PROMAX", 100000000, unitPrice, totalPrice)
		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, 100000000, product.Quantity)

		assert.NoError(t, product.IsValid())

		cloned := product.Clone()
		assert.Equal(t, product.Quantity, cloned.Quantity)

		cleanedOrder := product.ToCleanedOrder(1)
		assert.Equal(t, product.Quantity, cleanedOrder.Qty)
	})

	t.Run("Product with very small price", func(t *testing.T) {
		unitPrice := value_object.MustNewPrice(0.01)
		totalPrice := value_object.MustNewPrice(0.02)

		product, err := entity.NewProduct("FG0A-CLEAR-IPHONE16PROMAX", 2, unitPrice, totalPrice)
		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, 0.01, product.UnitPrice.Amount())

		assert.NoError(t, product.IsValid())

		cloned := product.Clone()
		assert.Equal(t, product.UnitPrice.Amount(), cloned.UnitPrice.Amount())
	})

	t.Run("Product with complex model ID containing multiple hyphens", func(t *testing.T) {
		unitPrice := value_object.MustNewPrice(50.0)
		totalPrice := value_object.MustNewPrice(100.0)

		product, err := entity.NewProduct("FG0A-CLEAR-OPPOA3-B-SPECIAL-EDITION", 2, unitPrice, totalPrice)
		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, "FG0A-CLEAR", product.MaterialId)
		assert.Equal(t, "OPPOA3-B-SPECIAL-EDITION", product.ModelId)

		assert.NoError(t, product.IsValid())

		cloned := product.Clone()
		assert.Equal(t, product.MaterialId, cloned.MaterialId)
		assert.Equal(t, product.ModelId, cloned.ModelId)

		cleanedOrder := product.ToCleanedOrder(1)
		assert.Equal(t, product.MaterialId, cleanedOrder.MaterialId)
		assert.Equal(t, product.ModelId, cleanedOrder.ModelId)
		assert.True(t, cleanedOrder.IsMainProduct())
		assert.False(t, cleanedOrder.IsComplementaryProduct())
	})

	t.Run("CleanedOrder with various product IDs", func(t *testing.T) {
		testCases := []struct {
			productId             string
			materialId            string
			modelId               string
			expectedMain          bool
			expectedComplementary bool
		}{
			{"FG0A-CLEAR-IPHONE16PROMAX", "FG0A-CLEAR", "IPHONE16PROMAX", true, false},
			{"WIPING-CLOTH", "", "", false, true},
			{"CLEAR-CLEANNER", "", "", false, true},
			{"MATTE-CLEANNER", "", "", false, true},
			{"PRIVACY-CLEANNER", "", "", false, true},
			{"CUSTOM-PRODUCT", "", "SOME-MODEL", false, true},
			{"ANOTHER-PRODUCT", "SOME-MATERIAL", "", false, true},
		}

		for _, tc := range testCases {
			order := &entity.CleanedOrder{
				ProductId:  tc.productId,
				MaterialId: tc.materialId,
				ModelId:    tc.modelId,
			}

			assert.Equal(t, tc.expectedMain, order.IsMainProduct(),
				"IsMainProduct for %s should be %v", tc.productId, tc.expectedMain)
			assert.Equal(t, tc.expectedComplementary, order.IsComplementaryProduct(),
				"IsComplementaryProduct for %s should be %v", tc.productId, tc.expectedComplementary)
		}
	})
}

func TestProduct_ComprehensiveScenarios(t *testing.T) {
	scenarios := []struct {
		name            string
		productId       string
		materialId      string
		modelId         string
		quantity        int
		unitPrice       float64
		totalPrice      float64
		shouldBeValid   bool
		expectedTexture string
	}{
		{
			name:            "Standard CLEAR product",
			productId:       "FG0A-CLEAR-IPHONE16PROMAX",
			materialId:      "FG0A-CLEAR",
			modelId:         "IPHONE16PROMAX",
			quantity:        2,
			unitPrice:       50.0,
			totalPrice:      100.0,
			shouldBeValid:   true,
			expectedTexture: "CLEAR",
		},
		{
			name:            "Standard MATTE product",
			productId:       "FG05-MATTE-OPPOA3",
			materialId:      "FG05-MATTE",
			modelId:         "OPPOA3",
			quantity:        1,
			unitPrice:       40.0,
			totalPrice:      40.0,
			shouldBeValid:   true,
			expectedTexture: "MATTE",
		},
		{
			name:            "Standard PRIVACY product",
			productId:       "FG1A-PRIVACY-SAMSUNGS25",
			materialId:      "FG1A-PRIVACY",
			modelId:         "SAMSUNGS25",
			quantity:        3,
			unitPrice:       60.0,
			totalPrice:      180.0,
			shouldBeValid:   true,
			expectedTexture: "PRIVACY",
		},
		{
			name:            "Product with complex model ID",
			productId:       "FG0A-CLEAR-OPPOA3-B-SPECIAL-EDITION",
			materialId:      "FG0A-CLEAR",
			modelId:         "OPPOA3-B-SPECIAL-EDITION",
			quantity:        1,
			unitPrice:       55.0,
			totalPrice:      55.0,
			shouldBeValid:   true,
			expectedTexture: "CLEAR",
		},
		{
			name:            "Product with zero prices",
			productId:       "FG0A-MATTE-TESTPHONE",
			materialId:      "FG0A-MATTE",
			modelId:         "TESTPHONE",
			quantity:        1,
			unitPrice:       0.0,
			totalPrice:      0.0,
			shouldBeValid:   true,
			expectedTexture: "MATTE",
		},
		{
			name:            "Invalid product with empty material",
			productId:       "FG0A-CLEAR-IPHONE16PROMAX",
			materialId:      "",
			modelId:         "IPHONE16PROMAX",
			quantity:        1,
			unitPrice:       50.0,
			totalPrice:      50.0,
			shouldBeValid:   false,
			expectedTexture: "",
		},
		{
			name:            "Invalid product with negative quantity",
			productId:       "FG0A-CLEAR-IPHONE16PROMAX",
			materialId:      "FG0A-CLEAR",
			modelId:         "IPHONE16PROMAX",
			quantity:        -1,
			unitPrice:       50.0,
			totalPrice:      50.0,
			shouldBeValid:   false,
			expectedTexture: "CLEAR",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			unitPrice := value_object.MustNewPrice(scenario.unitPrice)
			totalPrice := value_object.MustNewPrice(scenario.totalPrice)

			product := &entity.Product{
				ProductId:  scenario.productId,
				MaterialId: scenario.materialId,
				ModelId:    scenario.modelId,
				Quantity:   scenario.quantity,
				UnitPrice:  unitPrice,
				TotalPrice: totalPrice,
			}

			err := product.IsValid()
			if scenario.shouldBeValid {
				assert.NoError(t, err, "Product should be valid")
			} else {
				assert.Error(t, err, "Product should be invalid")
				assert.Equal(t, errors.ErrInvalidInput, err)
				return
			}

			texture := product.GetTexture()
			assert.Equal(t, scenario.expectedTexture, texture, "Texture should match expected")

			cleanedOrder := product.ToCleanedOrder(1)
			assert.Equal(t, scenario.productId, cleanedOrder.ProductId)
			assert.Equal(t, scenario.materialId, cleanedOrder.MaterialId)
			assert.Equal(t, scenario.modelId, cleanedOrder.ModelId)
			assert.Equal(t, scenario.quantity, cleanedOrder.Qty)
			assert.Equal(t, scenario.unitPrice, cleanedOrder.UnitPrice.Amount())
			assert.Equal(t, scenario.totalPrice, cleanedOrder.TotalPrice.Amount())

			cloned := product.Clone()
			assert.Equal(t, product.ProductId, cloned.ProductId)
			assert.Equal(t, product.MaterialId, cloned.MaterialId)
			assert.Equal(t, product.ModelId, cloned.ModelId)
			assert.Equal(t, product.Quantity, cloned.Quantity)
			assert.Equal(t, product.UnitPrice.Amount(), cloned.UnitPrice.Amount())
			assert.Equal(t, product.TotalPrice.Amount(), cloned.TotalPrice.Amount())

			cloned.Quantity = 999
			assert.NotEqual(t, product.Quantity, cloned.Quantity, "Clone should be independent")
		})
	}
}

func TestProduct_IntegrationWithCleanedOrder(t *testing.T) {
	t.Run("Main product flow", func(t *testing.T) {
		unitPrice := value_object.MustNewPrice(50.0)
		totalPrice := value_object.MustNewPrice(100.0)

		product, err := entity.NewProduct("FG0A-CLEAR-IPHONE16PROMAX", 2, unitPrice, totalPrice)
		require.NoError(t, err)

		assert.NoError(t, product.IsValid())
		assert.Equal(t, "CLEAR", product.GetTexture())

		cleanedOrder := product.ToCleanedOrder(1)
		assert.True(t, cleanedOrder.IsMainProduct())
		assert.False(t, cleanedOrder.IsComplementaryProduct())

		assert.NoError(t, cleanedOrder.IsValid())

		cloned := product.Clone()
		clonedOrder := cloned.ToCleanedOrder(2)
		assert.True(t, clonedOrder.IsMainProduct())
		assert.Equal(t, 2, clonedOrder.No)
	})

	t.Run("Complementary product flow", func(t *testing.T) {
		wipingCloth := &entity.CleanedOrder{
			No:         2,
			ProductId:  "WIPING-CLOTH",
			MaterialId: "",
			ModelId:    "",
			Qty:        2,
			UnitPrice:  value_object.ZeroPrice(),
			TotalPrice: value_object.ZeroPrice(),
		}

		assert.False(t, wipingCloth.IsMainProduct())
		assert.True(t, wipingCloth.IsComplementaryProduct())
		assert.NoError(t, wipingCloth.IsValid())

		cleaner := &entity.CleanedOrder{
			No:         3,
			ProductId:  "CLEAR-CLEANNER",
			MaterialId: "",
			ModelId:    "",
			Qty:        2,
			UnitPrice:  value_object.ZeroPrice(),
			TotalPrice: value_object.ZeroPrice(),
		}

		assert.False(t, cleaner.IsMainProduct())
		assert.True(t, cleaner.IsComplementaryProduct())
		assert.NoError(t, cleaner.IsValid())
	})
}

func TestProduct_ErrorHandling(t *testing.T) {
	t.Run("Nil price handling", func(t *testing.T) {
		product := &entity.Product{
			ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
			MaterialId: "FG0A-CLEAR",
			ModelId:    "IPHONE16PROMAX",
			Quantity:   2,
			UnitPrice:  nil,
			TotalPrice: value_object.MustNewPrice(100.0),
		}

		err := product.IsValid()
		assert.Error(t, err)
		assert.Equal(t, errors.ErrInvalidInput, err)
	})

	t.Run("Invalid quantity scenarios", func(t *testing.T) {
		testCases := []int{0, -1, -100}

		for _, qty := range testCases {
			product := &entity.Product{
				ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
				MaterialId: "FG0A-CLEAR",
				ModelId:    "IPHONE16PROMAX",
				Quantity:   qty,
				UnitPrice:  value_object.MustNewPrice(50.0),
				TotalPrice: value_object.MustNewPrice(100.0),
			}

			err := product.IsValid()
			assert.Error(t, err, "Quantity %d should be invalid", qty)
			assert.Equal(t, errors.ErrInvalidInput, err)
		}
	})
}

func TestProduct_Performance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance tests in short mode")
	}

	t.Run("Large batch processing", func(t *testing.T) {
		const batchSize = 1000

		products := make([]*entity.Product, batchSize)
		for i := 0; i < batchSize; i++ {
			unitPrice := value_object.MustNewPrice(50.0)
			totalPrice := value_object.MustNewPrice(100.0)

			product, err := entity.NewProduct("FG0A-CLEAR-IPHONE16PROMAX", 2, unitPrice, totalPrice)
			require.NoError(t, err)
			products[i] = product
		}

		start := time.Now()
		for _, product := range products {
			assert.NoError(t, product.IsValid())
		}
		validationTime := time.Since(start)
		t.Logf("Validation of %d products took: %v", batchSize, validationTime)

		start = time.Now()
		for _, product := range products {
			_ = product.Clone()
		}
		cloneTime := time.Since(start)
		t.Logf("Cloning of %d products took: %v", batchSize, cloneTime)

		start = time.Now()
		for i, product := range products {
			_ = product.ToCleanedOrder(i + 1)
		}
		conversionTime := time.Since(start)
		t.Logf("Conversion of %d products took: %v", batchSize, conversionTime)
	})
}

func TestProduct_Memory(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping memory tests in short mode")
	}

	t.Run("Memory usage validation", func(t *testing.T) {
		unitPrice := value_object.MustNewPrice(50.0)
		totalPrice := value_object.MustNewPrice(100.0)

		original := &entity.Product{
			ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
			MaterialId: "FG0A-CLEAR",
			ModelId:    "IPHONE16PROMAX",
			Quantity:   2,
			UnitPrice:  unitPrice,
			TotalPrice: totalPrice,
		}

		clones := make([]*entity.Product, 100)
		for i := 0; i < 100; i++ {
			clones[i] = original.Clone()
		}

		clones[0].ProductId = "MODIFIED"

		assert.Equal(t, "FG0A-CLEAR-IPHONE16PROMAX", original.ProductId)
		for i := 1; i < 100; i++ {
			assert.Equal(t, "FG0A-CLEAR-IPHONE16PROMAX", clones[i].ProductId)
		}
		assert.Equal(t, "MODIFIED", clones[0].ProductId)
	})
}
