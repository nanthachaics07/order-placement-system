// pkg/utils/parser/product_parser.go
package parser

import (
	"regexp"
	"strconv"
	"strings"

	"order-placement-system/internal/domain/entity"
	"order-placement-system/internal/usecases/interfaces"
	"order-placement-system/pkg/errors"
	"order-placement-system/pkg/log"
)

type ProductParserImpl struct{}

func NewProductParser() interfaces.ProductParser {
	return &ProductParserImpl{}
}

func (p *ProductParserImpl) Parse(platformProductId string, originalQty int, totalPrice float64) ([]*entity.ParsedProduct, error) {
	if platformProductId == "" {
		log.Error("platform product id cannot be empty")
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

	// ? mod
	pricePerUnit := totalPrice / float64(totalQuantityUnits) // calculate price per unit

	for i, bundleProduct := range bundleProducts {

		cleanProduct, _, _ := p.ExtractQuantity(bundleProduct)
		quantity := productQuantities[i]

		productTotalPrice := pricePerUnit * float64(quantity) // calculate total price for this product
		unitPrice := pricePerUnit

		parsedProduct := &entity.ParsedProduct{
			CleanProductId: cleanProduct,
			Quantity:       quantity,
			OriginalQty:    originalQty,
			UnitPrice:      unitPrice,
			TotalPrice:     productTotalPrice,
		}

		parsedProducts = append(parsedProducts, parsedProduct)
	}

	return parsedProducts, nil
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
		// if matches[1] is a valid integer, extract it
		if qty, err := strconv.Atoi(matches[1]); err == nil {
			cleanId = re.ReplaceAllString(productId, "")
			quantity = qty
			hasQuantity = true
			return
		}
	}

	// none or invalid quantity found
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

		// if strings.HasPrefix(part, "%20x") {
		// 	part = strings.TrimPrefix(part, "%20x")
		// }
		part = strings.TrimPrefix(part, "%20x")

		if part != "" {
			cleanParts = append(cleanParts, part)
		}
	}

	return cleanParts
}

func (p *ProductParserImpl) ParseProductCode(productId string) (materialId, modelId string, err error) {
	if productId == "" {
		log.Error("product id cannot be empty")
		return "", "", errors.ErrInvalidInput
	}

	parts := strings.Split(productId, "-")
	if len(parts) < 3 {
		log.Errorf("invalid product format", log.S("productId", productId))
		return "", "", errors.ErrInvalidInput
	}

	// MaterialId to {film type ID}-{texture ID}
	materialId = strings.Join(parts[:2], "-")

	// ModelId to {phone model ID}
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
	// Product code should start with "FG"
	return strings.HasPrefix(s, "FG")
}
