package transform

import (
	"fmt"
	"time"

	fp "github.com/ahiho/gocandy/filterparser"
	"github.com/ahiho/gocandy/filterparser/validator"
)

func Transform(filter fp.Filter, rules validator.Rules) error {
	for i := range filter.Conditions {
		if err := transformSingle(&filter.Conditions[i], rules); err != nil {
			return fmt.Errorf("condition #%d: %w", i+1, err)
		}
	}

	return nil
}

func transformSingle(cond *fp.Condition, rules validator.Rules) error {
	field, ok := rules.Field(cond.Field)
	if !ok {
		return fmt.Errorf("invalid field: %q", cond.Field)
	}

	if err := interpretTypes(cond, field); err != nil {
		return fmt.Errorf("field: %q: interpret types: %w", cond.Field, err)
	}

	return nil
}

func interpretTypes(cond *fp.Condition, field validator.Field) error {
	if field.Type == fp.TypeTimestamp {
		for i, curr := range cond.Values {
			s, ok := curr.(string)
			if !ok {
				return fmt.Errorf("timestamp: value #%d: not a string", i+1)
			}
			t, err := time.Parse(time.RFC3339, s)
			if err != nil {
				// err.Error() returns the value, we don't need to
				return fmt.Errorf("timestamp: value #%d: %w", i+1, err)
			}
			cond.Values[i] = t
		}
	}

	return nil
}
