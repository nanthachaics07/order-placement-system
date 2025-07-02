package value_object

import (
	"fmt"
	"order-placement-system/pkg/errors"
	"order-placement-system/pkg/log"
	"strings"
)

type Material struct {
	FilmTypeID string  `json:"film_type_id"`
	Texture    Texture `json:"texture"`
}

func NewMaterial(filmTypeID string, texture Texture) (*Material, error) {
	if filmTypeID == "" {
		log.Error("film type id cannot be empty")
		return nil, errors.ErrInvalidInput
	}

	if !texture.IsValid() {
		log.Errorf("invalid texture", log.S("texture", texture.String()))
		return nil, errors.ErrInvalidInput
	}

	return &Material{
		FilmTypeID: strings.ToUpper(strings.TrimSpace(filmTypeID)),
		Texture:    texture,
	}, nil
}

// FG0A-CLEAR to Material{FilmTypeID: "FG0A", Texture: TextureClear}
func NewMaterialFromString(materialId string) (*Material, error) {
	if materialId == "" {
		log.Error("material id cannot be empty")
		return nil, errors.ErrInvalidInput
	}

	parts := strings.Split(materialId, "-")
	if len(parts) < 2 {
		log.Errorf("invalid material id format", log.S("materialId", materialId))
		return nil, errors.ErrInvalidInput
	}

	filmTypeID := parts[0]
	textureStr := parts[1]

	texture, err := NewTexture(textureStr)
	if err != nil {
		log.Errorf("failed to create texture from material id", log.E(err), log.S("materialId", materialId))
		return nil, errors.ErrInvalidInput
	}

	return NewMaterial(filmTypeID, texture)
}

// convert to material id
func (m *Material) String() string {
	return fmt.Sprintf("%s-%s", m.FilmTypeID, m.Texture.String())
}

func (m *Material) IsValid() error {
	if m.FilmTypeID == "" {
		log.Error("film type id cannot be empty")
		return errors.ErrInvalidInput
	}

	if !m.Texture.IsValid() {
		log.Errorf("invalid texture", log.S("texture", m.Texture.String()))
		return errors.ErrInvalidInput
	}

	return nil
}

func (m *Material) Equals(other *Material) bool {
	if other == nil {
		return false
	}

	return m.FilmTypeID == other.FilmTypeID && m.Texture.Equals(other.Texture)
}

// get cleaner product id based on texture
func (m *Material) GetCleanerProductId() string {
	return m.Texture.GetCleanerProductId()
}

func (m *Material) IsCompatibleWith(other *Material) bool {

	return m.IsValid() == nil && other != nil && other.IsValid() == nil
}

func (m *Material) GetDisplayName() string {
	return fmt.Sprintf("%s %s", m.FilmTypeID, m.Texture.GetDisplayName())
}

// check if the material has a specific texture
func (m *Material) HasTexture(texture Texture) bool {
	return m.Texture.Equals(texture)
}

func (m *Material) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, m.String())), nil
}

func (m *Material) UnmarshalJSON(data []byte) error {

	s := strings.Trim(string(data), `"`)

	material, err := NewMaterialFromString(s)
	if err != nil {
		return err
	}

	*m = *material
	return nil
}

func (m *Material) Clone() *Material {
	return &Material{
		FilmTypeID: m.FilmTypeID,
		Texture:    m.Texture,
	}
}

func ValidateFilmTypeFormat(filmTypeID string) error {
	if filmTypeID == "" {
		log.Error("film type id cannot be empty")
		return errors.ErrInvalidInput
	}

	// check if it start with "FG"
	if !strings.HasPrefix(filmTypeID, "FG") {
		log.Errorf("film type id must start with 'FG'", log.S("filmTypeID", filmTypeID))
		return errors.ErrInvalidInput
	}

	if len(filmTypeID) < 3 {
		log.Errorf("film type id too short", log.S("filmTypeID", filmTypeID))
		return errors.ErrInvalidInput
	}

	return nil
}
