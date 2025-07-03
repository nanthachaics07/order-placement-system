package value_object_test

import (
	"encoding/json"
	"testing"

	"order-placement-system/internal/domain/value_object"
	"order-placement-system/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			// Marshal to JSON
			jsonData, err := json.Marshal(tt.texture)
			require.NoError(t, err)

			// Unmarshal back to texture
			var unmarshaledTexture value_object.Texture
			err = json.Unmarshal(jsonData, &unmarshaledTexture)
			require.NoError(t, err)

			// Check that they are equal
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
		// Test that different cases are handled correctly
		testCases := []string{"clear", "CLEAR", "Clear", "cLeAr"}
		for _, testCase := range testCases {
			texture, err := value_object.NewTexture(testCase)
			assert.NoError(t, err)
			assert.Equal(t, value_object.TextureClear, texture)
		}
	})

	t.Run("Whitespace handling", func(t *testing.T) {
		// Test whitespace trimming
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
