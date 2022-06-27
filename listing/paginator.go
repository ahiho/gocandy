package listing

import (
	"encoding/json"
	"fmt"

	"github.com/ahiho/gocandy/gormx/model"
)

// Interface Pagination represents a single pagination request.
type Pagination interface {
	// HasNextPage returns true if there are more pages left, false otherwise.
	// HasNextPage is called after Finish.
	HasNextPage() bool
	// ImplState returns the implementation-specific state of the pagination.
	ImplState() interface{}
	// MustEmbedCommonState returns the state of the pagination which is common across all pagination implementations.
	// Implementations embed CommonState to get this function.
	MustEmbedCommonState() CommonState
	// ModelHook is called for every model in the same order as they are returned by the driver.
	ModelHook(model *model.Common)
	// Finish is called after the client processed all records.
	Finish()
}

// splitRequest splits a Request into a CommonState and marshaled implementation-specific state.
func splitRequest(req Request) (CommonState, []byte, error) {
	ps := CommonState{
		Collection: req.Collection,
		Knobs:      req.Knobs,
	}

	if len(req.PageToken) != 0 {
		pt, err := unmarshalPageToken(req.PageToken)
		if err != nil {
			return ps, nil, fmt.Errorf("invalid page token: %w", err)
		}

		if pt.Common.Knobs != req.Knobs {
			return ps, nil, fmt.Errorf("list parameters changed: %#v -> %#v", req.Knobs, pt.Common.Knobs)
		}

		if pt.Common.Collection != req.Collection {
			return ps, nil, fmt.Errorf("invalid page token")
		}

		return ps, pt.State, nil
	} else {
		return ps, nil, nil
	}
}

// Init initializes a pagination from a Request.
// The common pagination state and the implementation specific state are separated and put into commonState and implState, respectively.
// implState must be the
func Init(req Request, commonState *CommonState, implState interface{}) error {
	if req.Knobs.PageSize <= 0 {
		return fmt.Errorf("invalid page size: %d", req.Knobs.PageSize)
	}

	cs, bytes, err := splitRequest(req)
	if err != nil {
		return err
	}
	*commonState = cs

	if bytes != nil {
		if err := json.Unmarshal(bytes, implState); err != nil {
			// Don't wrap, callers should not depend on marshaling errors
			return fmt.Errorf("unmarshal impl state: %v", err)
		}
	}

	return nil
}

// Finish concludes a pagination request.
// Finish returns a page token that can be passed to Init (as part of Request) to retrieve the next page or nil if there are no more pages.
func Finish(p Pagination) ([]byte, error) {
	p.Finish()

	if !p.HasNextPage() {
		return nil, nil
	}

	nextPageToken, err := nextPageToken(p)
	if err != nil {
		return nil, err
	}

	return nextPageToken, err
}

// nextPageToken returns a page token, combining the implementation-specific state of p and the common state.
func nextPageToken(p Pagination) ([]byte, error) {
	raw, err := json.Marshal(p.ImplState())
	if err != nil {
		return nil, err
	}

	npt := pageToken{
		Common: p.MustEmbedCommonState(),
		State:  raw,
	}

	return npt.Marshal()
}
