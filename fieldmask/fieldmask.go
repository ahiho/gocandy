// Package fieldmask provides support manipulating field masks.
package fieldmask

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"google.golang.org/genproto/protobuf/field_mask"

	"github.com/ahiho/gocandy/resource"
)

// MaxSize is the hard limit on the number of fields in a single field mask.
const MaxSize = 32

// Mask represents a list of fields in a resource.
type Mask struct {
	// Fields included in the mask, in snake case.
	Fields []string
	// Optional. Resource the mask applies to.
	Resource resource.Resource
}

// New returns a new Mask created from a protobuf mask.
func New(mask *field_mask.FieldMask) (*Mask, error) {
	if len(mask.Paths) > MaxSize {
		return nil, fmt.Errorf("number of fields is %d, maximum allowed is %d", len(mask.Paths), MaxSize)
	}

	fields := make([]string, len(mask.Paths))
	for i, k := range mask.Paths {
		fields[i] = strcase.ToSnake(k)
	}

	return &Mask{
		Fields: fields,
	}, nil
}

// NewResource returns a new Mask created from a protobuf mask against res.
func NewResource(mask *field_mask.FieldMask, res resource.Resource) (*Mask, error) {
	fieldMask, err := New(mask)
	if err != nil {
		return nil, err
	}

	fieldMask.Resource = res
	return fieldMask, nil
}

// RemoveOutputOnly removes all fields which are marked output only by the resource.
// Requires Resource to be set (either explicitly or via NewResource), otherwise panics.
func (f *Mask) RemoveOutputOnly() {
	newFields := make([]string, 0, len(f.Fields))
	for _, k := range f.Fields {
		if !f.Resource.IsFieldOutputOnly(k) {
			newFields = append(newFields, k)
		}
	}
	f.Fields = newFields
}

// Contains returns true if the mask contains the field.
func (f Mask) Contains(field string) bool {
	return f.Index(field) >= 0
}

// Index returns the index of the first instance of field in f, or -1 if field is not present in f.
func (f Mask) Index(field string) int {
	for i, curr := range f.Fields {
		if curr == field {
			return i
		}
	}
	return -1
}
