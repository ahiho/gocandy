package sql

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ahiho/gocandy/filter/parser"
)

// DefaultMatcherWithValidator wraps the default matcher with validation on the value.
func DefaultMatcherWithValidator(validate ValidatorFunc, comps []string) ParseValidateFunc {
	return func(ex *parser.Expression) (*SQLResponse, error) {
		for _, v := range comps {
			if v == ex.Comparator || v == "*" {
				var err error
				if ex.Comparator == parser.HASH.String() {
					values := strings.Split(strings.TrimLeft(strings.TrimRight(ex.Value, ")"), "("), ",")
					for _, v := range values {
						err = validate(v)
						if err != nil {
							return nil, errors.New("invalid value")
						}
					}
				} else {
					err = validate(ex.Value)
					if err != nil {
						return nil, errors.New("invalid value")
					}
				}
				return DefaultMatcher(ex), nil
			}
		}
		return nil, errors.New("field is not allowed")
	}
}

// DefaultMatcher takes an expression and spits out the default SqlResponse.
func DefaultMatcher(ex *parser.Expression) *SQLResponse {
	if ex.Comparator == parser.TokenLookup[parser.PERCENT] {
		fmtValue := fmt.Sprintf("%%%s%%", ex.Value)
		sq := SQLResponse{
			Raw:    fmt.Sprintf("%s LIKE ?", ex.Field),
			Values: []string{fmtValue},
		}
		return &sq
	}
	if ex.Comparator == parser.TokenLookup[parser.HASH] {
		values := strings.Split(strings.TrimLeft(strings.TrimRight(ex.Value, ")"), "("), ",")
		raw := fmt.Sprintf("(%s%s)", strings.Repeat("?, ", len(values)-1), "?")

		sq := SQLResponse{
			Raw:    fmt.Sprintf("%s IN %s", ex.Field, raw),
			Values: values,
		}
		return &sq
	}
	sq := SQLResponse{
		Raw:    fmt.Sprintf("%s%s?", ex.Field, ex.Comparator),
		Values: []string{ex.Value},
	}
	return &sq
}

// NullValidator is a no-op validator on a string, always returns nil error.
func NullValidator(_ string) error {
	return nil
}

// IntegerValidator validates that the input is an integer.
func IntegerValidator(s string) error {
	_, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("value '%s' is not an integer", s)
	}
	return nil
}

// NumericValidator validates that the input is a number.
func NumericValidator(s string) error {
	_, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("value '%s' is not numeric", s)
	}
	return nil
}
