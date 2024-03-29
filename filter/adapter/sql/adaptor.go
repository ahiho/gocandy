// Package sql_adaptor provides functions to convert a goven query to a valid and safe SQL query.
package sql

import (
	"reflect"
	"regexp"
	"strings"
)

const (
	tagName = "filter"
)

// NewDefaultAdaptorFromStruct returns a new basic SQLAdaptor from the reflection of your database object.
func NewDefaultAdaptorFromStruct(gorm reflect.Value) *SQLAdaptor {
	matchers := map[*regexp.Regexp]ParseValidateFunc{}
	fieldMappings := map[string]string{}
	defaultFields := FieldParseValidatorFromStruct(gorm)
	return NewSQLAdaptor(fieldMappings, defaultFields, matchers)
}

// FieldParseValidatorFromStruct takes the reflection of your database object and returns a map of fieldnames to ParseValidateFuncs.
// Don't panic - reflection is only used once on initialisation.
func FieldParseValidatorFromStruct(gorm reflect.Value) map[string]ParseValidateFunc {
	defaultFields := map[string]ParseValidateFunc{}
	e := gorm.Elem()

	for i := 0; i < e.NumField(); i++ {
		varName := strings.ToLower(e.Type().Field(i).Name)
		varType := e.Type().Field(i).Type
		vType := strings.TrimPrefix(varType.String(), "*")
		comps := strings.Split(e.Type().Field(i).Tag.Get(tagName), ";")
		switch vType {
		case "float32", "float64":
			defaultFields[varName] = DefaultMatcherWithValidator(NumericValidator, comps)
		case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
			defaultFields[varName] = DefaultMatcherWithValidator(IntegerValidator, comps)
		default:
			defaultFields[varName] = DefaultMatcherWithValidator(NullValidator, comps)
		}
	}
	return defaultFields
}
