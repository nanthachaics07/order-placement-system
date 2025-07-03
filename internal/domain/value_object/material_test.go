package value_object_test

import (
	"encoding/json"
	"testing"

	"order-placement-system/internal/domain/value_object"
	"order-placement-system/pkg/errors"
	"order-placement-system/pkg/log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	log.Init("dev")
}

func TestNewMaterial(t *testing.T) {
	testCases := []struct {
		name        string
		filmTypeID  string
		texture     value_object.Texture
		expected    *value_object.Material
		expectError bool
	}{
		{
			name:       "Valid material with CLEAR texture",
			filmTypeID: "FG0A",
			texture:    value_object.TextureClear,
			expected: &value_object.Material{
				FilmTypeID: "FG0A",
				Texture:    value_object.TextureClear,
			},
			expectError: false,
		},
		{
			name:       "Valid material with MATTE texture",
			filmTypeID: "FG05",
			texture:    value_object.TextureMatte,
			expected: &value_object.Material{
				FilmTypeID: "FG05",
				Texture:    value_object.TextureMatte,
			},
			expectError: false,
		},
		{
			name:       "Valid material with PRIVACY texture",
			filmTypeID: "FG1A",
			texture:    value_object.TexturePrivacy,
			expected: &value_object.Material{
				FilmTypeID: "FG1A",
				Texture:    value_object.TexturePrivacy,
			},
			expectError: false,
		},
		{
			name:       "Trim spaces in film type ID",
			filmTypeID: "  FG0A  ",
			texture:    value_object.TextureClear,
			expected: &value_object.Material{
				FilmTypeID: "FG0A",
				Texture:    value_object.TextureClear,
			},
			expectError: false,
		},
		{
			name:       "Lowercase film type ID should be converted to uppercase",
			filmTypeID: "fg0a",
			texture:    value_object.TextureClear,
			expected: &value_object.Material{
				FilmTypeID: "FG0A",
				Texture:    value_object.TextureClear,
			},
			expectError: false,
		},
		{
			name:        "Empty film type ID should return error",
			filmTypeID:  "",
			texture:     value_object.TextureClear,
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Invalid texture should return error",
			filmTypeID:  "FG0A",
			texture:     value_object.Texture("INVALID"),
			expected:    nil,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := value_object.NewMaterial(tc.filmTypeID, tc.texture)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, errors.ErrInvalidInput, err)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tc.expected.FilmTypeID, result.FilmTypeID)
				assert.Equal(t, tc.expected.Texture, result.Texture)
			}
		})
	}
}

func TestNewMaterialFromString(t *testing.T) {
	testCases := []struct {
		name        string
		materialId  string
		expected    *value_object.Material
		expectError bool
	}{
		{
			name:       "Valid material ID with CLEAR texture",
			materialId: "FG0A-CLEAR",
			expected: &value_object.Material{
				FilmTypeID: "FG0A",
				Texture:    value_object.TextureClear,
			},
			expectError: false,
		},
		{
			name:       "Valid material ID with MATTE texture",
			materialId: "FG05-MATTE",
			expected: &value_object.Material{
				FilmTypeID: "FG05",
				Texture:    value_object.TextureMatte,
			},
			expectError: false,
		},
		{
			name:       "Valid material ID with PRIVACY texture",
			materialId: "FG1A-PRIVACY",
			expected: &value_object.Material{
				FilmTypeID: "FG1A",
				Texture:    value_object.TexturePrivacy,
			},
			expectError: false,
		},
		{
			name:        "Empty material ID should return error",
			materialId:  "",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Material ID without separator should return error",
			materialId:  "FG0ACLEAR",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Material ID with only one part should return error",
			materialId:  "FG0A",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Material ID with invalid texture should return error",
			materialId:  "FG0A-INVALID",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Material ID with empty texture should return error",
			materialId:  "FG0A-",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Material ID with empty film type should return error",
			materialId:  "-CLEAR",
			expected:    nil,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := value_object.NewMaterialFromString(tc.materialId)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, errors.ErrInvalidInput, err)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tc.expected.FilmTypeID, result.FilmTypeID)
				assert.Equal(t, tc.expected.Texture, result.Texture)
			}
		})
	}
}

func TestMaterial_String(t *testing.T) {
	testCases := []struct {
		name     string
		material *value_object.Material
		expected string
	}{
		{
			name: "Material with CLEAR texture",
			material: &value_object.Material{
				FilmTypeID: "FG0A",
				Texture:    value_object.TextureClear,
			},
			expected: "FG0A-CLEAR",
		},
		{
			name: "Material with MATTE texture",
			material: &value_object.Material{
				FilmTypeID: "FG05",
				Texture:    value_object.TextureMatte,
			},
			expected: "FG05-MATTE",
		},
		{
			name: "Material with PRIVACY texture",
			material: &value_object.Material{
				FilmTypeID: "FG1A",
				Texture:    value_object.TexturePrivacy,
			},
			expected: "FG1A-PRIVACY",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.material.String()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestMaterial_IsValid(t *testing.T) {
	testCases := []struct {
		name      string
		material  *value_object.Material
		expectErr bool
	}{
		{
			name: "Valid material",
			material: &value_object.Material{
				FilmTypeID: "FG0A",
				Texture:    value_object.TextureClear,
			},
			expectErr: false,
		},
		{
			name: "Material with empty film type ID",
			material: &value_object.Material{
				FilmTypeID: "",
				Texture:    value_object.TextureClear,
			},
			expectErr: true,
		},
		{
			name: "Material with invalid texture",
			material: &value_object.Material{
				FilmTypeID: "FG0A",
				Texture:    value_object.Texture("INVALID"),
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.material.IsValid()

			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, errors.ErrInvalidInput, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMaterial_Equals(t *testing.T) {
	material1 := &value_object.Material{
		FilmTypeID: "FG0A",
		Texture:    value_object.TextureClear,
	}

	material2 := &value_object.Material{
		FilmTypeID: "FG0A",
		Texture:    value_object.TextureClear,
	}

	material3 := &value_object.Material{
		FilmTypeID: "FG05",
		Texture:    value_object.TextureClear,
	}

	material4 := &value_object.Material{
		FilmTypeID: "FG0A",
		Texture:    value_object.TextureMatte,
	}

	testCases := []struct {
		name     string
		material *value_object.Material
		other    *value_object.Material
		expected bool
	}{
		{
			name:     "Same materials should be equal",
			material: material1,
			other:    material2,
			expected: true,
		},
		{
			name:     "Different film types should not be equal",
			material: material1,
			other:    material3,
			expected: false,
		},
		{
			name:     "Different textures should not be equal",
			material: material1,
			other:    material4,
			expected: false,
		},
		{
			name:     "Comparing with nil should return false",
			material: material1,
			other:    nil,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.material.Equals(tc.other)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestMaterial_GetCleanerProductId(t *testing.T) {
	testCases := []struct {
		name     string
		material *value_object.Material
		expected string
	}{
		{
			name: "CLEAR texture should return CLEAR-CLEANNER",
			material: &value_object.Material{
				FilmTypeID: "FG0A",
				Texture:    value_object.TextureClear,
			},
			expected: "CLEAR-CLEANNER",
		},
		{
			name: "MATTE texture should return MATTE-CLEANNER",
			material: &value_object.Material{
				FilmTypeID: "FG05",
				Texture:    value_object.TextureMatte,
			},
			expected: "MATTE-CLEANNER",
		},
		{
			name: "PRIVACY texture should return PRIVACY-CLEANNER",
			material: &value_object.Material{
				FilmTypeID: "FG1A",
				Texture:    value_object.TexturePrivacy,
			},
			expected: "PRIVACY-CLEANNER",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.material.GetCleanerProductId()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestMaterial_IsCompatibleWith(t *testing.T) {
	validMaterial := &value_object.Material{
		FilmTypeID: "FG0A",
		Texture:    value_object.TextureClear,
	}

	anotherValidMaterial := &value_object.Material{
		FilmTypeID: "FG05",
		Texture:    value_object.TextureMatte,
	}

	invalidMaterial := &value_object.Material{
		FilmTypeID: "",
		Texture:    value_object.TextureClear,
	}

	testCases := []struct {
		name     string
		material *value_object.Material
		other    *value_object.Material
		expected bool
	}{
		{
			name:     "Both valid materials should be compatible",
			material: validMaterial,
			other:    anotherValidMaterial,
			expected: true,
		},
		{
			name:     "Valid material with invalid material should not be compatible",
			material: validMaterial,
			other:    invalidMaterial,
			expected: false,
		},
		{
			name:     "Invalid material with valid material should not be compatible",
			material: invalidMaterial,
			other:    validMaterial,
			expected: false,
		},
		{
			name:     "Valid material with nil should not be compatible",
			material: validMaterial,
			other:    nil,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.material.IsCompatibleWith(tc.other)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestMaterial_GetDisplayName(t *testing.T) {
	testCases := []struct {
		name     string
		material *value_object.Material
		expected string
	}{
		{
			name: "CLEAR texture display name",
			material: &value_object.Material{
				FilmTypeID: "FG0A",
				Texture:    value_object.TextureClear,
			},
			expected: "FG0A Clear",
		},
		{
			name: "MATTE texture display name",
			material: &value_object.Material{
				FilmTypeID: "FG05",
				Texture:    value_object.TextureMatte,
			},
			expected: "FG05 Matte",
		},
		{
			name: "PRIVACY texture display name",
			material: &value_object.Material{
				FilmTypeID: "FG1A",
				Texture:    value_object.TexturePrivacy,
			},
			expected: "FG1A Privacy",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.material.GetDisplayName()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestMaterial_HasTexture(t *testing.T) {
	material := &value_object.Material{
		FilmTypeID: "FG0A",
		Texture:    value_object.TextureClear,
	}

	testCases := []struct {
		name     string
		material *value_object.Material
		texture  value_object.Texture
		expected bool
	}{
		{
			name:     "Material has CLEAR texture",
			material: material,
			texture:  value_object.TextureClear,
			expected: true,
		},
		{
			name:     "Material does not have MATTE texture",
			material: material,
			texture:  value_object.TextureMatte,
			expected: false,
		},
		{
			name:     "Material does not have PRIVACY texture",
			material: material,
			texture:  value_object.TexturePrivacy,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.material.HasTexture(tc.texture)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestMaterial_MarshalJSON(t *testing.T) {
	testCases := []struct {
		name     string
		material *value_object.Material
		expected string
	}{
		{
			name: "Marshal CLEAR material",
			material: &value_object.Material{
				FilmTypeID: "FG0A",
				Texture:    value_object.TextureClear,
			},
			expected: `"FG0A-CLEAR"`,
		},
		{
			name: "Marshal MATTE material",
			material: &value_object.Material{
				FilmTypeID: "FG05",
				Texture:    value_object.TextureMatte,
			},
			expected: `"FG05-MATTE"`,
		},
		{
			name: "Marshal PRIVACY material",
			material: &value_object.Material{
				FilmTypeID: "FG1A",
				Texture:    value_object.TexturePrivacy,
			},
			expected: `"FG1A-PRIVACY"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.material.MarshalJSON()
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, string(result))
		})
	}
}

func TestMaterial_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name        string
		jsonData    string
		expected    *value_object.Material
		expectError bool
	}{
		{
			name:     "Unmarshal CLEAR material",
			jsonData: `"FG0A-CLEAR"`,
			expected: &value_object.Material{
				FilmTypeID: "FG0A",
				Texture:    value_object.TextureClear,
			},
			expectError: false,
		},
		{
			name:     "Unmarshal MATTE material",
			jsonData: `"FG05-MATTE"`,
			expected: &value_object.Material{
				FilmTypeID: "FG05",
				Texture:    value_object.TextureMatte,
			},
			expectError: false,
		},
		{
			name:     "Unmarshal PRIVACY material",
			jsonData: `"FG1A-PRIVACY"`,
			expected: &value_object.Material{
				FilmTypeID: "FG1A",
				Texture:    value_object.TexturePrivacy,
			},
			expectError: false,
		},
		{
			name:        "Unmarshal invalid format should return error",
			jsonData:    `"FG0A"`,
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Unmarshal empty string should return error",
			jsonData:    `""`,
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Unmarshal invalid texture should return error",
			jsonData:    `"FG0A-INVALID"`,
			expected:    nil,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var material value_object.Material
			err := material.UnmarshalJSON([]byte(tc.jsonData))

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected.FilmTypeID, material.FilmTypeID)
				assert.Equal(t, tc.expected.Texture, material.Texture)
			}
		})
	}
}

func TestMaterial_Clone(t *testing.T) {
	original := &value_object.Material{
		FilmTypeID: "FG0A",
		Texture:    value_object.TextureClear,
	}

	cloned := original.Clone()

	// Should be equal but not the same instance
	assert.True(t, original.Equals(cloned))
	assert.NotSame(t, original, cloned)

	// Modifying cloned should not affect original
	cloned.FilmTypeID = "FG05"
	assert.Equal(t, "FG0A", original.FilmTypeID)
	assert.Equal(t, "FG05", cloned.FilmTypeID)
}

func TestValidateFilmTypeFormat(t *testing.T) {
	testCases := []struct {
		name        string
		filmTypeID  string
		expectError bool
	}{
		{
			name:        "Valid film type FG0A",
			filmTypeID:  "FG0A",
			expectError: false,
		},
		{
			name:        "Valid film type FG05",
			filmTypeID:  "FG05",
			expectError: false,
		},
		{
			name:        "Valid film type FG123",
			filmTypeID:  "FG123",
			expectError: false,
		},
		{
			name:        "Empty film type should return error",
			filmTypeID:  "",
			expectError: true,
		},
		{
			name:        "Film type not starting with FG should return error",
			filmTypeID:  "AB0A",
			expectError: true,
		},
		{
			name:        "Film type too short should return error",
			filmTypeID:  "FG",
			expectError: true,
		},
		{
			name:        "Film type with only F should return error",
			filmTypeID:  "F",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := value_object.ValidateFilmTypeFormat(tc.filmTypeID)

			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, errors.ErrInvalidInput, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMaterial_JSONRoundTrip(t *testing.T) {
	// Test complete JSON marshal/unmarshal cycle
	original := &value_object.Material{
		FilmTypeID: "FG0A",
		Texture:    value_object.TextureClear,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(original)
	require.NoError(t, err)

	// Unmarshal from JSON
	var restored value_object.Material
	err = json.Unmarshal(jsonData, &restored)
	require.NoError(t, err)

	// Should be equal
	assert.True(t, original.Equals(&restored))
}

func TestMaterial_EdgeCases(t *testing.T) {
	t.Run("Material with whitespace in film type ID should be trimmed", func(t *testing.T) {
		material, err := value_object.NewMaterial("  FG0A  ", value_object.TextureClear)
		require.NoError(t, err)
		assert.Equal(t, "FG0A", material.FilmTypeID)
	})

	t.Run("Material string should handle special characters", func(t *testing.T) {
		material := &value_object.Material{
			FilmTypeID: "FG0A",
			Texture:    value_object.TextureClear,
		}
		result := material.String()
		assert.Equal(t, "FG0A-CLEAR", result)
		assert.NotContains(t, result, " ")
	})

	t.Run("IsCompatibleWith should handle nil pointer", func(t *testing.T) {
		material := &value_object.Material{
			FilmTypeID: "FG0A",
			Texture:    value_object.TextureClear,
		}
		result := material.IsCompatibleWith(nil)
		assert.False(t, result)
	})
}
