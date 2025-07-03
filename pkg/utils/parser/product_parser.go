package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"order-placement-system/internal/domain/entity"
	"order-placement-system/internal/domain/service"
	"order-placement-system/internal/domain/value_object"
	"order-placement-system/pkg/errors"
	"order-placement-system/pkg/log"
)

type ProductParserImpl struct {
	priceCalculator service.PriceCalculator
}

func NewProductParser() service.ProductParser {
	return &ProductParserImpl{
		priceCalculator: NewPriceCalculator(),
	}
}

func (p *ProductParserImpl) Parse(platformProductId string, originalQty int, totalPrice *value_object.Price) ([]*entity.ParsedProduct, error) {
	if platformProductId == "" {
		log.Error("platform product id cannot be empty")
		return nil, errors.ErrInvalidInput
	}

	if totalPrice == nil {
		log.Error("total price cannot be nil")
		return nil, errors.ErrInvalidInput
	}

	cleanedId := p.CleanPrefix(platformProductId)
	bundleProducts := p.SplitBundle(cleanedId)

	var parsedProducts []*entity.ParsedProduct
	totalQuantityUnits := 0
	productQuantities := make([]int, len(bundleProducts))

	for i, bundleProduct := range bundleProducts {
		_, quantity, hasQuantity := p.ExtractQuantity(bundleProduct)
		if !hasQuantity {
			quantity = originalQty
		}
		productQuantities[i] = quantity
		totalQuantityUnits += quantity
	}

	pricePerUnit, err := p.priceCalculator.CalculateUnitPrice(totalPrice, totalQuantityUnits)
	if err != nil {
		log.Errorf("failed to calculate unit price", log.E(err))
		return nil, err
	}

	for i, bundleProduct := range bundleProducts {
		cleanProduct, _, _ := p.ExtractQuantity(bundleProduct)
		quantity := productQuantities[i]

		productTotalPrice, err := p.priceCalculator.CalculateTotalPrice(pricePerUnit, quantity)
		if err != nil {
			log.Errorf("failed to calculate product total price", log.E(err))
			return nil, err
		}

		parsedProduct := &entity.ParsedProduct{
			CleanProductId: cleanProduct,
			Quantity:       quantity,
			OriginalQty:    originalQty,
			UnitPrice:      pricePerUnit,
			TotalPrice:     productTotalPrice,
		}

		parsedProducts = append(parsedProducts, parsedProduct)
	}

	return parsedProducts, nil
}

func (p *ProductParserImpl) ParseFromFloat64(platformProductId string, originalQty int, totalPrice float64) ([]*entity.ParsedProduct, error) {
	totalPriceVO, err := value_object.NewPrice(totalPrice)
	if err != nil {
		log.Errorf("invalid total price", log.S("price", strconv.FormatFloat(totalPrice, 'f', 2, 64)), log.E(err))
		return nil, errors.ErrInvalidInput
	}

	return p.Parse(platformProductId, originalQty, totalPriceVO)
}

func (p *ProductParserImpl) CleanPrefix(productId string) string {
	if productId == "" {
		return ""
	}

	cleaned := productId

	prefixes := []string{
		"%20--%20x",
		"%20--",
		"--%20x",
		"x2-3&",
		"%20x",
		"%20-",
		"--",
	}

	for {
		before := cleaned

		for _, prefix := range prefixes {
			if strings.HasPrefix(cleaned, prefix) {
				cleaned = cleaned[len(prefix):]
				goto next
			}
		}

		if strings.HasPrefix(cleaned, "-") {
			if !p.isValidProductStart(cleaned[1:]) {
				cleaned = cleaned[1:]
				goto next
			}
		}

		if cleaned == before {
			break
		}

	next:
		continue
	}

	return cleaned
}

func (p *ProductParserImpl) ExtractQuantity(productId string) (cleanId string, quantity int, hasQuantity bool) {
	re := regexp.MustCompile(`\*(\d+)$`)
	matches := re.FindStringSubmatch(productId)

	if len(matches) == 2 {
		if qty, err := strconv.Atoi(matches[1]); err == nil {
			cleanId = re.ReplaceAllString(productId, "")
			quantity = qty
			hasQuantity = true
			return
		}
	}

	cleanId = productId
	quantity = 1
	hasQuantity = false
	return
}

func (p *ProductParserImpl) SplitBundle(productId string) []string {
	parts := strings.Split(productId, "/")
	cleanParts := make([]string, 0)

	for _, part := range parts {
		part = strings.TrimSpace(part)
		part = strings.TrimPrefix(part, "%20x")

		if part != "" {
			fixedPart := p.fixIncompleteProductId(part)
			cleanParts = append(cleanParts, fixedPart)
		}
	}

	return cleanParts
}

func (p *ProductParserImpl) fixIncompleteProductId(productId string) string {
	parts := strings.Split(productId, "-")

	if len(parts) == 2 {
		filmType := parts[0]
		texture := parts[1]

		if texture == "MAT" {
			texture = "MATTE"
		}

		modelId := p.inferModelId(filmType, texture, productId)

		if modelId != "" {
			fixedId := fmt.Sprintf("%s-%s-%s", filmType, texture, modelId)
			log.Debugf("Fixed incomplete product ID",
				log.S("original", productId),
				log.S("fixed", fixedId))
			return fixedId
		}
	}

	return productId
}

func (p *ProductParserImpl) inferModelId(filmType, texture, originalId string) string {
	knownPatterns := map[string]string{
		"FG0A-MATTE": "OPPOA3",
		"FG0A-CLEAR": "OPPOA3",
		"FG05-MATTE": "OPPOA3",
	}

	key := fmt.Sprintf("%s-%s", filmType, texture)
	if modelId, exists := knownPatterns[key]; exists {
		return modelId
	}

	return "OPPOA3"
}

func (p *ProductParserImpl) ParseProductCode(productId string) (materialId, modelId string, err error) {
	if productId == "" {
		log.Error("product id cannot be empty")
		return "", "", errors.ErrInvalidInput
	}

	parts := strings.Split(productId, "-")
	if len(parts) < 3 {
		log.Errorf("invalid product format - expected at least 3 parts separated by '-', got %d parts",
			log.S("productId", productId),
			log.S("parts", fmt.Sprintf("%v", parts)))
		return "", "", errors.ErrInvalidInput
	}

	filmType := parts[0]
	texture := p.normalizeTexture(parts[1])

	if !p.isValidFilmType(filmType) {
		log.Errorf("invalid film type", log.S("filmType", filmType))
		return "", "", errors.ErrInvalidInput
	}

	if !p.isValidTexture(texture) {
		log.Errorf("invalid texture", log.S("texture", texture))
		return "", "", errors.ErrInvalidInput
	}

	if parts[2] == "" {
		log.Errorf("model id cannot be empty", log.S("productId", productId))
		return "", "", errors.ErrInvalidInput
	}

	materialId = fmt.Sprintf("%s-%s", filmType, texture)
	modelId = strings.Join(parts[2:], "-")

	return materialId, modelId, nil
}

func (p *ProductParserImpl) Validate(productId string) error {
	if productId == "" {
		log.Error("product id cannot be empty")
		return errors.ErrInvalidInput
	}

	if strings.Count(productId, "-") < 2 {
		log.Errorf("invalid product code format", log.S("productId", productId))
		return errors.ErrInvalidInput
	}

	return nil
}

func (p *ProductParserImpl) isValidProductStart(s string) bool {
	if len(s) < 2 {
		return false
	}
	return strings.HasPrefix(s, "FG")
}

func (p *ProductParserImpl) isValidFilmType(filmType string) bool {
	validFilmTypes := []string{"FG0A", "FG05", "FG1A", "FG1B"}
	for _, valid := range validFilmTypes {
		if filmType == valid {
			return true
		}
	}
	return strings.HasPrefix(filmType, "FG") && len(filmType) >= 3
}

func (p *ProductParserImpl) isValidTexture(texture string) bool {
	validTextures := []string{"CLEAR", "MATTE", "PRIVACY"}
	for _, valid := range validTextures {
		if texture == valid {
			return true
		}
	}
	return false
}

func (p *ProductParserImpl) normalizeTexture(texture string) string {
	switch strings.ToUpper(texture) {
	case "MAT":
		return "MATTE"
	case "CLEAR":
		return "CLEAR"
	case "MATTE":
		return "MATTE"
	case "PRIVACY":
		return "PRIVACY"
	default:
		log.Debugf("unknown texture, normalizing to uppercase", log.S("texture", texture))
		return strings.ToUpper(texture)
	}
}

type PriceCalculatorImpl struct{}

func NewPriceCalculator() service.PriceCalculator {
	return &PriceCalculatorImpl{}
}

func (c *PriceCalculatorImpl) CalculateUnitPrice(totalPrice *value_object.Price, quantity int) (*value_object.Price, error) {
	if totalPrice == nil {
		return nil, errors.ErrInvalidInput
	}

	if quantity <= 0 {
		return nil, errors.ErrInvalidInput
	}

	return totalPrice.DivideByInt(quantity)
}

func (c *PriceCalculatorImpl) CalculateTotalPrice(unitPrice *value_object.Price, quantity int) (*value_object.Price, error) {
	if unitPrice == nil {
		return nil, errors.ErrInvalidInput
	}

	if quantity <= 0 {
		return nil, errors.ErrInvalidInput
	}

	return unitPrice.MultiplyByInt(quantity)
}

func (c *PriceCalculatorImpl) DividePriceEqually(totalPrice *value_object.Price, parts int) (*value_object.Price, error) {
	if totalPrice == nil {
		return nil, errors.ErrInvalidInput
	}

	if parts <= 0 {
		return nil, errors.ErrInvalidInput
	}

	return totalPrice.DivideByInt(parts)
}

func (c *PriceCalculatorImpl) SumPrices(prices ...*value_object.Price) (*value_object.Price, error) {
	total := value_object.ZeroPrice()

	for _, price := range prices {
		if price != nil {
			var err error
			total, err = total.Add(price)
			if err != nil {
				return nil, err
			}
		}
	}

	return total, nil
}
