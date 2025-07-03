package value_object

import (
	"encoding/json"
	"fmt"
	"math"
	"order-placement-system/pkg/errors"
	"order-placement-system/pkg/log"
)

type Price struct {
	amount float64
}

func NewPrice(amount float64) (*Price, error) {
	if amount < 0 {
		log.Errorf("price cannot be negative", amount)
		return nil, errors.ErrInvalidInput
	}

	if math.IsNaN(amount) || math.IsInf(amount, 0) {
		log.Errorf("price must be a valid number", amount)
		return nil, errors.ErrInvalidInput
	}

	return &Price{amount: amount}, nil
}

func MustNewPrice(amount float64) *Price {
	if amount < 0 {
		panic(fmt.Sprintf("price cannot be negative: %f", amount))
	}
	price, err := NewPrice(amount)
	if err != nil {
		panic(fmt.Sprintf("invalid price: %v", err))
	}
	return price
}

func ZeroPrice() *Price {
	return &Price{amount: 0}
}

func (p *Price) Amount() float64 {
	if p == nil {
		return 0
	}
	return p.amount
}

func (p *Price) IsZero() bool {
	return p == nil || p.amount == 0
}

func (p *Price) IsPositive() bool {
	return p != nil && p.amount > 0
}

func (p *Price) Add(other *Price) (*Price, error) {
	if p == nil {
		p = ZeroPrice()
	}
	if other == nil {
		other = ZeroPrice()
	}

	return NewPrice(p.amount + other.amount)
}

func (p *Price) Subtract(other *Price) (*Price, error) {
	if p == nil {
		p = ZeroPrice()
	}
	if other == nil {
		other = ZeroPrice()
	}

	return NewPrice(p.amount - other.amount)
}

func (p *Price) Multiply(multiplier float64) (*Price, error) {
	if p == nil {
		return ZeroPrice(), nil
	}

	return NewPrice(p.amount * multiplier)
}

func (p *Price) MultiplyByInt(quantity int) (*Price, error) {
	if quantity < 0 {
		log.Errorf("quantity cannot be negative", quantity)
		return nil, errors.ErrInvalidInput
	}

	return p.Multiply(float64(quantity))
}

func (p *Price) Divide(divisor float64) (*Price, error) {
	if divisor == 0 {
		log.Error("cannot divide by zero")
		return nil, errors.ErrInvalidInput
	}

	if p == nil {
		return ZeroPrice(), nil
	}

	return NewPrice(p.amount / divisor)
}

func (p *Price) DivideByInt(divisor int) (*Price, error) {
	if divisor == 0 {
		log.Error("cannot divide by zero")
		return nil, errors.ErrInvalidInput
	}

	return p.Divide(float64(divisor))
}

func (p *Price) Equals(other *Price) bool {
	if p == nil && other == nil {
		return true
	}

	if p == nil || other == nil {
		return false
	}

	const epsilon = 1e-9
	return math.Abs(p.amount-other.amount) < epsilon
}

func (p *Price) GreaterThan(other *Price) bool {
	if p == nil {
		return false
	}
	if other == nil {
		return p.amount > 0
	}

	return p.amount > other.amount
}

func (p *Price) LessThan(other *Price) bool {
	return other.GreaterThan(p)
}

func (p *Price) String() string {
	if p == nil {
		return "0.00"
	}
	return fmt.Sprintf("%.2f", p.amount)
}

func (p *Price) MarshalJSON() ([]byte, error) {
	if p == nil {
		return []byte("0.00"), nil
	}
	return []byte(fmt.Sprintf("%.2f", p.amount)), nil
}

func (p *Price) UnmarshalJSON(data []byte) error {
	var amount float64
	if err := json.Unmarshal(data, &amount); err != nil {
		return err
	}

	price, err := NewPrice(amount)
	if err != nil {
		return err
	}

	*p = *price
	return nil
}

func (p *Price) Clone() *Price {
	if p == nil {
		return nil
	}

	return &Price{amount: p.amount}
}

func (p *Price) Round(precision int) *Price {
	if p == nil {
		return ZeroPrice()
	}

	multiplier := math.Pow(10, float64(precision))
	rounded := math.Round(p.amount*multiplier) / multiplier

	return MustNewPrice(rounded)
}

func (p *Price) ToDisplayString(currency string) string {
	if currency == "" {
		currency = "THB"
	}

	return fmt.Sprintf("%s %.2f", currency, p.Amount())
}
