package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	fp "github.com/ahiho/gocandy/filterparser"
)

// Field represents restrictions on a single field.
type Field struct {
	Type      fp.Type
	Container fp.Container
	Operators Operators
}

type Operators []fp.Operator

// Contains returns true if Operators contains op.
func (o Operators) Contains(op fp.Operator) bool {
	for _, curr := range o {
		if curr == op {
			return true
		}
	}

	return false
}

// Rules defines restrictions for filter expressions.
type Rules struct {
	Fields map[string]Field
}

// Field returns the Field corresponding to a potentially nested field name.
//
// If key contains "." (dot), the dot and everything after it is ignored when looking up the field. The Field is only returned if it is a filterparser.ContainerMap.
// If key does not contain a dot, Field is similar to r.Fields[key] but with an additional check to ensures it is not a filterparser.ContainerMap.
func (r Rules) Field(key string) (Field, bool) {
	needMap := false
	if s := strings.Split(key, "."); len(s) > 1 {
		key = s[0]
		needMap = true
	}

	field, ok := r.Fields[key]
	if !ok {
		return field, false
	}

	// If the input is a nested field name, the top level field must be a map.
	// If the input is NOT a nested field name, the top level field must NOT be a map.
	if isMap := (field.Container == fp.ContainerMap); needMap != isMap {
		return field, false
	}

	return field, ok
}

func (r Rules) Validate(f fp.Filter) error {
	for i := range f.Conditions {
		if err := r.validateCondition(&f.Conditions[i]); err != nil {
			return fmt.Errorf("condition #%d: %v", i+1, err)
		}
	}

	return nil
}

func (f Field) Ops() Operators {
	if len(f.Operators) == 0 {
		// No operators specified, return the default set of operators for the type
		switch f.Type {
		case fp.TypeString:
			return []fp.Operator{fp.OpEqual, fp.OpNotEqual, fp.OpIn, fp.OpContains}
		case fp.TypeInt, fp.TypeUInt, fp.TypeTimestamp, fp.TypeFloat:
			return []fp.Operator{fp.OpEqual, fp.OpNotEqual, fp.OpIn, fp.OpGreater, fp.OpLess, fp.OpGreaterOrEqual, fp.OpLessOrEqual, fp.OpRange}
		case fp.TypeBool:
			return []fp.Operator{fp.OpEqual, fp.OpNotEqual}
		}
	}

	return f.Operators
}

func (r Rules) validateConditionAgainstField(c fp.Condition, field Field) error {
	// Check if op is valid for field
	if !field.Ops().Contains(c.Op) {
		return fmt.Errorf("operator %q invalid", c.Op)
	}

	// Check if only one type of value is supplied
	if err := checkSameType(c); err != nil {
		return err
	}

	if err := validateValue(c, field); err != nil {
		return fmt.Errorf("invalid value: %v", err)
	}

	return nil
}

func isKind(v interface{}, want reflect.Kind) error {
	if kind := reflect.TypeOf(v).Kind(); kind != want {
		return fmt.Errorf("%s required, got %s", want, kind)
	}

	return nil
}

func validateValue(c fp.Condition, field Field) error {
	// Assume that all values are of the same type
	switch field.Type {
	case fp.TypeString:
		return isKind(c.Values[0], reflect.String)

	case fp.TypeUInt:
		if v, ok := c.Values[0].(int64); !ok || v < 0 {
			return errors.New("uint required")
		}

	case fp.TypeInt:
		return isKind(c.Values[0], reflect.Int64)

	case fp.TypeBool:
		return isKind(c.Values[0], reflect.Bool)

	case fp.TypeFloat:
		return isKind(c.Values[0], reflect.Float64)

	case fp.TypeTimestamp:
		if _, ok := c.Values[0].(time.Time); !ok {
			return errors.New("timestamp required")
		}

	default:
		panic("unknown field type")
	}

	return nil
}

func (r Rules) validateCondition(cond *fp.Condition) error {
	field, ok := r.Field(cond.Field)
	if !ok {
		return fmt.Errorf("invalid field: %q", cond.Field)
	}

	if err := r.validateConditionAgainstField(*cond, field); err != nil {
		return fmt.Errorf("field %q: validate conditions: %v", cond.Field, err)
	}

	return nil
}

func checkSameType(cond fp.Condition) error {
	if len(cond.Values) < 2 {
		// No values in the slice or only a single value cannot have a different type than itself
		return nil
	}

	t := reflect.TypeOf(cond.Values[0])
	for i, curr := range cond.Values[1:] {
		if reflect.TypeOf(curr) != t {
			return fmt.Errorf("value %d is of different type", i+1)
		}
	}

	return nil
}
