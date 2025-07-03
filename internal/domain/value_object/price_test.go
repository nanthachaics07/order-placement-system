package value_object_test

import (
	"encoding/json"
	"math"
	"order-placement-system/internal/domain/value_object"
	"order-placement-system/pkg/log"
	"testing"
)

func init() {
	log.Init("dev")
}

func TestNewPrice(t *testing.T) {
	tests := []struct {
		name    string
		amount  float64
		wantErr bool
	}{
		{
			name:    "valid positive price",
			amount:  50.0,
			wantErr: false,
		},
		{
			name:    "valid zero price",
			amount:  0.0,
			wantErr: false,
		},
		{
			name:    "valid decimal price",
			amount:  99.99,
			wantErr: false,
		},
		{
			name:    "negative price should error",
			amount:  -10.0,
			wantErr: true,
		},
		{
			name:    "NaN should error",
			amount:  math.NaN(),
			wantErr: true,
		},
		{
			name:    "positive infinity should error",
			amount:  math.Inf(1),
			wantErr: true,
		},
		{
			name:    "negative infinity should error",
			amount:  math.Inf(-1),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			price, err := value_object.NewPrice(tt.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && price.Amount() != tt.amount {
				t.Errorf("NewPrice() amount = %v, want %v", price.Amount(), tt.amount)
			}
		})
	}
}

func TestMustNewPrice(t *testing.T) {
	t.Run("valid price should not panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("MustNewPrice() panicked unexpectedly: %v", r)
			}
		}()
		price := value_object.MustNewPrice(50.0)
		if price.Amount() != 50.0 {
			t.Errorf("MustNewPrice() amount = %v, want 50.0", price.Amount())
		}
	})

	t.Run("invalid price should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("MustNewPrice() should have panicked")
			}
		}()
		value_object.MustNewPrice(-10.0)
	})
}

func TestZeroPrice(t *testing.T) {
	price := value_object.ZeroPrice()
	if price.Amount() != 0.0 {
		t.Errorf("ZeroPrice() amount = %v, want 0.0", price.Amount())
	}
	if !price.IsZero() {
		t.Error("ZeroPrice() should be zero")
	}
}

func TestPriceAmount(t *testing.T) {
	tests := []struct {
		name  string
		price *value_object.Price
		want  float64
	}{
		{
			name:  "valid price",
			price: value_object.MustNewPrice(50.0),
			want:  50.0,
		},
		{
			name:  "nil price",
			price: nil,
			want:  0.0,
		},
		{
			name:  "zero price",
			price: value_object.ZeroPrice(),
			want:  0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.price.Amount(); got != tt.want {
				t.Errorf("Price.Amount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPriceIsZero(t *testing.T) {
	tests := []struct {
		name  string
		price *value_object.Price
		want  bool
	}{
		{
			name:  "nil price is zero",
			price: nil,
			want:  true,
		},
		{
			name:  "zero price is zero",
			price: value_object.ZeroPrice(),
			want:  true,
		},
		{
			name:  "positive price is not zero",
			price: value_object.MustNewPrice(50.0),
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.price.IsZero(); got != tt.want {
				t.Errorf("Price.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPriceIsPositive(t *testing.T) {
	tests := []struct {
		name  string
		price *value_object.Price
		want  bool
	}{
		{
			name:  "nil price is not positive",
			price: nil,
			want:  false,
		},
		{
			name:  "zero price is not positive",
			price: value_object.ZeroPrice(),
			want:  false,
		},
		{
			name:  "positive price is positive",
			price: value_object.MustNewPrice(50.0),
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.price.IsPositive(); got != tt.want {
				t.Errorf("Price.IsPositive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPriceAdd(t *testing.T) {
	tests := []struct {
		name    string
		price1  *value_object.Price
		price2  *value_object.Price
		want    float64
		wantErr bool
	}{
		{
			name:    "add two positive prices",
			price1:  value_object.MustNewPrice(30.0),
			price2:  value_object.MustNewPrice(20.0),
			want:    50.0,
			wantErr: false,
		},
		{
			name:    "add price to zero",
			price1:  value_object.ZeroPrice(),
			price2:  value_object.MustNewPrice(25.0),
			want:    25.0,
			wantErr: false,
		},
		{
			name:    "add with nil first price",
			price1:  nil,
			price2:  value_object.MustNewPrice(25.0),
			want:    25.0,
			wantErr: false,
		},
		{
			name:    "add with nil second price",
			price1:  value_object.MustNewPrice(25.0),
			price2:  nil,
			want:    25.0,
			wantErr: false,
		},
		{
			name:    "add two nil prices",
			price1:  nil,
			price2:  nil,
			want:    0.0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.price1.Add(tt.price2)
			if (err != nil) != tt.wantErr {
				t.Errorf("Price.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result.Amount() != tt.want {
				t.Errorf("Price.Add() = %v, want %v", result.Amount(), tt.want)
			}
		})
	}
}

func TestPriceSubtract(t *testing.T) {
	tests := []struct {
		name    string
		price1  *value_object.Price
		price2  *value_object.Price
		want    float64
		wantErr bool
	}{
		{
			name:    "subtract smaller from larger",
			price1:  value_object.MustNewPrice(50.0),
			price2:  value_object.MustNewPrice(20.0),
			want:    30.0,
			wantErr: false,
		},
		{
			name:    "subtract from zero",
			price1:  value_object.ZeroPrice(),
			price2:  value_object.MustNewPrice(25.0),
			want:    -25.0,
			wantErr: true, // negative result should error
		},
		{
			name:    "subtract same values",
			price1:  value_object.MustNewPrice(25.0),
			price2:  value_object.MustNewPrice(25.0),
			want:    0.0,
			wantErr: false,
		},
		{
			name:    "subtract with nil second price",
			price1:  value_object.MustNewPrice(25.0),
			price2:  nil,
			want:    25.0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.price1.Subtract(tt.price2)
			if (err != nil) != tt.wantErr {
				t.Errorf("Price.Subtract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result.Amount() != tt.want {
				t.Errorf("Price.Subtract() = %v, want %v", result.Amount(), tt.want)
			}
		})
	}
}

func TestPriceMultiply(t *testing.T) {
	tests := []struct {
		name       string
		price      *value_object.Price
		multiplier float64
		want       float64
		wantErr    bool
	}{
		{
			name:       "multiply by positive number",
			price:      value_object.MustNewPrice(25.0),
			multiplier: 2.0,
			want:       50.0,
			wantErr:    false,
		},
		{
			name:       "multiply by zero",
			price:      value_object.MustNewPrice(25.0),
			multiplier: 0.0,
			want:       0.0,
			wantErr:    false,
		},
		{
			name:       "multiply by decimal",
			price:      value_object.MustNewPrice(100.0),
			multiplier: 0.5,
			want:       50.0,
			wantErr:    false,
		},
		{
			name:       "multiply nil price",
			price:      nil,
			multiplier: 2.0,
			want:       0.0,
			wantErr:    false,
		},
		{
			name:       "multiply by negative number",
			price:      value_object.MustNewPrice(25.0),
			multiplier: -2.0,
			want:       -50.0,
			wantErr:    true, // negative result should error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.price.Multiply(tt.multiplier)
			if (err != nil) != tt.wantErr {
				t.Errorf("Price.Multiply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result.Amount() != tt.want {
				t.Errorf("Price.Multiply() = %v, want %v", result.Amount(), tt.want)
			}
		})
	}
}

func TestPriceMultiplyByInt(t *testing.T) {
	tests := []struct {
		name     string
		price    *value_object.Price
		quantity int
		want     float64
		wantErr  bool
	}{
		{
			name:     "multiply by positive integer",
			price:    value_object.MustNewPrice(25.0),
			quantity: 3,
			want:     75.0,
			wantErr:  false,
		},
		{
			name:     "multiply by zero",
			price:    value_object.MustNewPrice(25.0),
			quantity: 0,
			want:     0.0,
			wantErr:  false,
		},
		{
			name:     "multiply by negative integer",
			price:    value_object.MustNewPrice(25.0),
			quantity: -2,
			want:     0.0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.price.MultiplyByInt(tt.quantity)
			if (err != nil) != tt.wantErr {
				t.Errorf("Price.MultiplyByInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result.Amount() != tt.want {
				t.Errorf("Price.MultiplyByInt() = %v, want %v", result.Amount(), tt.want)
			}
		})
	}
}

func TestPriceDivide(t *testing.T) {
	tests := []struct {
		name    string
		price   *value_object.Price
		divisor float64
		want    float64
		wantErr bool
	}{
		{
			name:    "divide by positive number",
			price:   value_object.MustNewPrice(100.0),
			divisor: 2.0,
			want:    50.0,
			wantErr: false,
		},
		{
			name:    "divide by decimal",
			price:   value_object.MustNewPrice(100.0),
			divisor: 0.5,
			want:    200.0,
			wantErr: false,
		},
		{
			name:    "divide by zero",
			price:   value_object.MustNewPrice(100.0),
			divisor: 0.0,
			want:    0.0,
			wantErr: true,
		},
		{
			name:    "divide nil price",
			price:   nil,
			divisor: 2.0,
			want:    0.0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.price.Divide(tt.divisor)
			if (err != nil) != tt.wantErr {
				t.Errorf("Price.Divide() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result.Amount() != tt.want {
				t.Errorf("Price.Divide() = %v, want %v", result.Amount(), tt.want)
			}
		})
	}
}

func TestPriceDivideByInt(t *testing.T) {
	tests := []struct {
		name    string
		price   *value_object.Price
		divisor int
		want    float64
		wantErr bool
	}{
		{
			name:    "divide by positive integer",
			price:   value_object.MustNewPrice(100.0),
			divisor: 4,
			want:    25.0,
			wantErr: false,
		},
		{
			name:    "divide by zero",
			price:   value_object.MustNewPrice(100.0),
			divisor: 0,
			want:    0.0,
			wantErr: true,
		},
		{
			name:    "divide by one",
			price:   value_object.MustNewPrice(100.0),
			divisor: 1,
			want:    100.0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.price.DivideByInt(tt.divisor)
			if (err != nil) != tt.wantErr {
				t.Errorf("Price.DivideByInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result.Amount() != tt.want {
				t.Errorf("Price.DivideByInt() = %v, want %v", result.Amount(), tt.want)
			}
		})
	}
}

func TestPriceEquals(t *testing.T) {
	tests := []struct {
		name   string
		price1 *value_object.Price
		price2 *value_object.Price
		want   bool
	}{
		{
			name:   "equal prices",
			price1: value_object.MustNewPrice(50.0),
			price2: value_object.MustNewPrice(50.0),
			want:   true,
		},
		{
			name:   "different prices",
			price1: value_object.MustNewPrice(50.0),
			price2: value_object.MustNewPrice(60.0),
			want:   false,
		},
		{
			name:   "both nil",
			price1: nil,
			price2: nil,
			want:   true,
		},
		{
			name:   "one nil",
			price1: value_object.MustNewPrice(50.0),
			price2: nil,
			want:   false,
		},
		{
			name:   "very close prices (within epsilon)",
			price1: value_object.MustNewPrice(50.0),
			price2: value_object.MustNewPrice(50.0000000001),
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.price1.Equals(tt.price2); got != tt.want {
				t.Errorf("value_object.Price.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPriceGreaterThan(t *testing.T) {
	tests := []struct {
		name   string
		price1 *value_object.Price
		price2 *value_object.Price
		want   bool
	}{
		{
			name:   "first price greater",
			price1: value_object.MustNewPrice(60.0),
			price2: value_object.MustNewPrice(50.0),
			want:   true,
		},
		{
			name:   "first price smaller",
			price1: value_object.MustNewPrice(40.0),
			price2: value_object.MustNewPrice(50.0),
			want:   false,
		},
		{
			name:   "equal prices",
			price1: value_object.MustNewPrice(50.0),
			price2: value_object.MustNewPrice(50.0),
			want:   false,
		},
		{
			name:   "first price nil",
			price1: nil,
			price2: value_object.MustNewPrice(50.0),
			want:   false,
		},
		{
			name:   "second price nil",
			price1: value_object.MustNewPrice(50.0),
			price2: nil,
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.price1.GreaterThan(tt.price2); got != tt.want {
				t.Errorf("value_object.Price.GreaterThan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPriceLessThan(t *testing.T) {
	tests := []struct {
		name   string
		price1 *value_object.Price
		price2 *value_object.Price
		want   bool
	}{
		{
			name:   "first price less",
			price1: value_object.MustNewPrice(40.0),
			price2: value_object.MustNewPrice(50.0),
			want:   true,
		},
		{
			name:   "first price greater",
			price1: value_object.MustNewPrice(60.0),
			price2: value_object.MustNewPrice(50.0),
			want:   false,
		},
		{
			name:   "equal prices",
			price1: value_object.MustNewPrice(50.0),
			price2: value_object.MustNewPrice(50.0),
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.price1.LessThan(tt.price2); got != tt.want {
				t.Errorf("value_object.Price.LessThan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPriceString(t *testing.T) {
	tests := []struct {
		name  string
		price *value_object.Price
		want  string
	}{
		{
			name:  "positive price",
			price: value_object.MustNewPrice(50.0),
			want:  "50.00",
		},
		{
			name:  "zero price",
			price: value_object.ZeroPrice(),
			want:  "0.00",
		},
		{
			name:  "decimal price",
			price: value_object.MustNewPrice(99.99),
			want:  "99.99",
		},
		{
			name:  "nil price",
			price: nil,
			want:  "0.00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.price.String(); got != tt.want {
				t.Errorf("value_object.Price.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPriceClone(t *testing.T) {
	original := value_object.MustNewPrice(50.0)
	cloned := original.Clone()

	if cloned == nil {
		t.Error("Clone() returned nil")
		return
	}

	if !original.Equals(cloned) {
		t.Error("Clone() should equal original")
	}

	// Test that they are different objects
	if original == cloned {
		t.Error("Clone() should return different object")
	}

	// Test nil clone
	var nilPrice *value_object.Price
	nilClone := nilPrice.Clone()
	if nilClone != nil {
		t.Error("Clone() of nil should return nil")
	}
}

func TestPriceRound(t *testing.T) {
	tests := []struct {
		name      string
		price     *value_object.Price
		precision int
		want      float64
	}{
		{
			name:      "round to 2 decimals",
			price:     value_object.MustNewPrice(50.126),
			precision: 2,
			want:      50.13,
		},
		{
			name:      "round to 1 decimal",
			price:     value_object.MustNewPrice(50.14),
			precision: 1,
			want:      50.1,
		},
		{
			name:      "round to 0 decimals",
			price:     value_object.MustNewPrice(50.6),
			precision: 0,
			want:      51.0,
		},
		{
			name:      "round nil price",
			price:     nil,
			precision: 2,
			want:      0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.price.Round(tt.precision)
			if result.Amount() != tt.want {
				t.Errorf("value_object.Price.Round() = %v, want %v", result.Amount(), tt.want)
			}
		})
	}
}

func TestPriceToDisplayString(t *testing.T) {
	tests := []struct {
		name     string
		price    *value_object.Price
		currency string
		want     string
	}{
		{
			name:     "with THB currency",
			price:    value_object.MustNewPrice(50.0),
			currency: "THB",
			want:     "THB 50.00",
		},
		{
			name:     "with USD currency",
			price:    value_object.MustNewPrice(99.99),
			currency: "USD",
			want:     "USD 99.99",
		},
		{
			name:     "with empty currency (default THB)",
			price:    value_object.MustNewPrice(50.0),
			currency: "",
			want:     "THB 50.00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.price.ToDisplayString(tt.currency); got != tt.want {
				t.Errorf("value_object.Price.ToDisplayString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPriceJSON(t *testing.T) {
	t.Run("marshal JSON", func(t *testing.T) {
		price := value_object.MustNewPrice(50.0)
		data, err := json.Marshal(price)
		if err != nil {
			t.Errorf("MarshalJSON() error = %v", err)
			return
		}
		want := "50.00"
		if string(data) != want {
			t.Errorf("MarshalJSON() = %v, want %v", string(data), want)
		}
	})

	t.Run("marshal nil price", func(t *testing.T) {
		var price *value_object.Price
		data, err := json.Marshal(price)
		if err != nil {
			t.Errorf("MarshalJSON() error = %v", err)
			return
		}
		want := "null"
		if string(data) != want {
			t.Errorf("MarshalJSON() = %v, want %v", string(data), want)
		}
	})

	t.Run("unmarshal JSON", func(t *testing.T) {
		data := []byte("50.00")
		var price value_object.Price
		err := json.Unmarshal(data, &price)
		if err != nil {
			t.Errorf("UnmarshalJSON() error = %v", err)
			return
		}
		if price.Amount() != 50.0 {
			t.Errorf("UnmarshalJSON() amount = %v, want 50.0", price.Amount())
		}
	})

	t.Run("unmarshal invalid JSON", func(t *testing.T) {
		data := []byte("invalid")
		var price value_object.Price
		err := json.Unmarshal(data, &price)
		if err == nil {
			t.Error("UnmarshalJSON() should error on invalid JSON")
		}
	})

	t.Run("unmarshal negative price", func(t *testing.T) {
		data := []byte("-50.00")
		var price value_object.Price
		err := json.Unmarshal(data, &price)
		if err == nil {
			t.Error("UnmarshalJSON() should error on negative price")
		}
	})
}

// Benchmark tests
func BenchmarkPriceAdd(b *testing.B) {
	price1 := value_object.MustNewPrice(50.0)
	price2 := value_object.MustNewPrice(25.0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = price1.Add(price2)
	}
}

func BenchmarkPriceMultiply(b *testing.B) {
	price := value_object.MustNewPrice(50.0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = price.Multiply(2.0)
	}
}

func BenchmarkPriceEquals(b *testing.B) {
	price1 := value_object.MustNewPrice(50.0)
	price2 := value_object.MustNewPrice(50.0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = price1.Equals(price2)
	}
}
