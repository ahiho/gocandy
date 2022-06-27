// Package listing implements generic pagination and page tokens.
package listing

// Request contains common parameters for list requests.
type Request struct {
	// Knobs controlling the result set.
	Knobs Knobs
	// Collection being enumerated.
	Collection string
	// Listing (pagination) state. Empty for the first page or a value returned from a previous call to the listing function for subsequent pages.
	PageToken []byte
}

// Knobs represents parameters controlling filtering and ordering of the results of list requests.
type Knobs struct {
	// True to include deleted resources in the list, false to omit them.
	ShowDeleted bool
	// Maximum number of resources to return
	PageSize int
	// Filtering condition
	Filter string
	// Ordering clause
	OrderBy string
}

func (l *Knobs) SetPageSize(input, def, max int) {
	if input <= 0 {
		input = def
	} else if input > max {
		input = max
	}
	l.PageSize = input
}
