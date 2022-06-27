package listing

import (
	"encoding/json"
	"errors"
	"fmt"
)

var ErrPageToken = errors.New("invalid page token")

// CommonState represents the part of the pagination state that is common across all pagination implementations.
type CommonState struct {
	// Collection to list.
	Collection string `json:"c"`
	// Parameters of the pagination.
	Knobs Knobs `json:"k"`
}

func (c CommonState) MustEmbedCommonState() CommonState { return c }

// pageToken represents all information in a page token.
// Marshal marshals a pageToken into a stream of bytes; unmarshalPageToken does the reverse.
type pageToken struct {
	// State understood by all implementations of pagination.
	Common CommonState `json:"p"`
	// Implementation specific pagination state
	State json.RawMessage `json:"s"`
}

func unmarshalPageToken(data []byte) (*pageToken, error) {
	pt := pageToken{}
	if err := json.Unmarshal(data, &pt); err != nil {
		return nil, fmt.Errorf("json.Unmarshal(): %w", err)
	}

	return &pt, nil
}

func (p pageToken) Marshal() ([]byte, error) {
	bytes, err := json.Marshal(&p)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
