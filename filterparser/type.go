package filterparser

import (
	"reflect"
)

// Type represents the type of a field or a constant value in a filter expression.
//
// Comparisons can only be made if both sides of the comparison are of the same type.
// The type of a field determines the default set of comparison operators allowed for it.
//
// A zero value for Type is invalid, passing it to functions in this library causes undefined behavior or panic.
type Type uint

const (
	TypeString Type = iota + 1
	TypeInt
	TypeUInt
	TypeFloat
	TypeBool
	TypeTimestamp
)

var kindToType = map[reflect.Kind]Type{
	reflect.Bool:    TypeBool,
	reflect.Int:     TypeInt,
	reflect.Int64:   TypeInt,
	reflect.Int32:   TypeInt,
	reflect.Int16:   TypeInt,
	reflect.Int8:    TypeInt,
	reflect.Uint:    TypeUInt,
	reflect.Uint64:  TypeUInt,
	reflect.Uint32:  TypeUInt,
	reflect.Uint16:  TypeUInt,
	reflect.Uint8:   TypeUInt,
	reflect.String:  TypeString,
	reflect.Float64: TypeFloat,
	reflect.Float32: TypeFloat,
}

var stringToType = map[string]Type{
	"bool":   TypeBool,
	"int":    TypeInt,
	"uint":   TypeUInt,
	"string": TypeString,
	"float":  TypeFloat,

	"timestamp": TypeTimestamp,
}

// TypeFromKind returns a Type based on a reflect.Kind.
func TypeFromKind(kind reflect.Kind) (t Type, ok bool) {
	t, ok = kindToType[kind]
	return t, ok
}

// TypeFromString returns a Type based on a string.
//
// Recognized strings: "bool", "int", "uint", "string", "float", "timestamp".
func TypeFromString(str string) (t Type, ok bool) {
	t, ok = stringToType[str]
	return t, ok
}
