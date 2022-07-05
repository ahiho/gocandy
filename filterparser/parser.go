// Package filterparser provides a parser for filter expressions.
package filterparser

import (
	"fmt"

	"github.com/ahiho/gocandy/filterparser/internal"
)

// Condition represents a single filtering condition.
type Condition struct {
	// Field to compare.
	Field string
	// True to invert the sense of the condition.
	Not bool
	// Check to perform.
	Op Operator
	// Values supplied to the operator.
	Values []interface{}
}

// Filter represents a parsed filter expression.
type Filter struct {
	// Conditions in the filter expression.
	Conditions []Condition
}

// Parse parses a filter expression.
func Parse(input string) (*Filter, error) {
	e := internal.Expression{}
	if err := internal.Parse(input, &e); err != nil {
		return nil, fmt.Errorf("parse: %v", err)
	}

	filter, err := convert(e)
	if err != nil {
		return nil, fmt.Errorf("convert: %v", err)
	}

	return filter, nil
}

func opFromCond(c internal.Condition) Operator {
	if c.In != nil {
		return OpIn
	} else if c.Compare != nil {
		op, ok := OperatorFromString(c.Compare.Operator)
		if !ok {
			panic("invalid operator: " + c.Compare.Operator)
		}
		return op
	} else if c.Between != nil {
		return OpRange
	}

	panic("invalid ast condition")
}

func appendValue(f *Condition, v internal.Value) error {
	// Possible enhancement: error if too many values, useful to keep the number of values in an in condition from getting too high

	if v.Int != nil {
		f.Values = append(f.Values, *v.Int)
	} else if v.String != nil {
		f.Values = append(f.Values, *v.String)
	} else if v.Float != nil {
		f.Values = append(f.Values, *v.Float)
	}

	return nil
}

func fillValue(f *Condition, c internal.Condition) error {
	if c.In != nil {
		for _, curr := range c.In.Values {
			if err := appendValue(f, curr); err != nil {
				return err
			}
		}
	} else if c.Compare != nil {
		// Handle booleans only here
		if c.Compare.Value.Boolean != nil {
			f.Values = []interface{}{bool(*c.Compare.Value.Boolean)}
		} else {
			if err := appendValue(f, c.Compare.Value); err != nil {
				return err
			}
		}
	} else if c.Between != nil {
		if err := appendValue(f, c.Between.Start); err != nil {
			return err
		}

		if err := appendValue(f, c.Between.End); err != nil {
			return err
		}
	}

	return nil
}

func convert(expr internal.Expression) (*Filter, error) {
	filter := Filter{
		Conditions: make([]Condition, len(expr.And)),
	}

	for i, curr := range expr.And {
		c := Condition{
			Not:   curr.Not,
			Field: curr.Symbol,
			Op:    opFromCond(curr),
		}

		if err := fillValue(&c, curr); err != nil {
			return nil, err
		}

		filter.Conditions[i] = c
	}

	return &filter, nil
}
