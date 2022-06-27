// Package resource declares the generic resource interface.
package resource

type Resource interface {
	IsFieldOutputOnly(field string) bool
}
