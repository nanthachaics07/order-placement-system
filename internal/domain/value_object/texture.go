package value_object

import (
	"fmt"
	"order-placement-system/pkg/errors"
	"order-placement-system/pkg/log"
	"strings"
)

type Texture string

const (
	TextureClear   Texture = "CLEAR"
	TextureMatte   Texture = "MATTE"
	TexturePrivacy Texture = "PRIVACY"
)

var AllTextures = []Texture{
	TextureClear,
	TextureMatte,
	TexturePrivacy,
}

func NewTexture(s string) (Texture, error) {
	texture := Texture(strings.ToUpper(strings.TrimSpace(s)))

	if !texture.IsValid() {
		log.Errorf("invalid texture", log.S("texture", s))
		return "", errors.ErrInvalidInput
	}

	return texture, nil
}

func (t Texture) IsValid() bool {
	for _, validTexture := range AllTextures {
		if t == validTexture {
			return true
		}
	}
	return false
}

func (t Texture) String() string {
	return string(t)
}

// get cleaner product id for this texture
func (t Texture) GetCleanerProductId() string {
	return t.String() + "-CLEANNER"
}

func (t Texture) Equals(other Texture) bool {
	return t == other
}

func (t Texture) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.String())), nil
}

func (t *Texture) UnmarshalJSON(data []byte) error {

	s := strings.Trim(string(data), `"`)

	texture, err := NewTexture(s)
	if err != nil {
		log.Errorf("failed to unmarshal texture", log.E(err), log.S("data", string(data)))
		return err
	}

	*t = texture
	return nil
}

func (t Texture) IsCompatibleWithFilmType(filmType string) bool {

	return t.IsValid()
}

func (t Texture) GetDisplayName() string {
	switch t {
	case TextureClear:
		return "Clear"
	case TextureMatte:
		return "Matte"
	case TexturePrivacy:
		return "Privacy"
	default:
		return t.String()
	}
}

func (t Texture) GetPriority() int {
	switch t {
	case TextureClear:
		return 1
	case TextureMatte:
		return 2
	case TexturePrivacy:
		return 3
	default:
		return 0
	}
}

// FG0A-CLEAR to TextureClear get texture from material id
func ParseTextureFromMaterialId(materialId string) (Texture, error) {
	if materialId == "" {
		log.Error("material id cannot be empty")
		return "", errors.ErrInvalidInput
	}

	parts := strings.Split(materialId, "-")
	if len(parts) < 2 {
		log.Errorf("invalid material id format", log.S("materialId", materialId))
		return "", errors.ErrInvalidInput
	}

	textureStr := parts[1] // texture part
	return NewTexture(textureStr)
}
