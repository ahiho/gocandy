package gormx

import (
	"time"

	"gorm.io/gorm"

	"github.com/ahiho/gocandy/gormx/model"
	"github.com/ahiho/gocandy/listing"
)

// sqlPaginationState is the implementation-specific state of an SQL driven pagination.
type sqlPaginationState struct {
	Date time.Time `json:"d"`
	ID   int64     `json:"i"`
}

type sqlPagination struct {
	listing.CommonState
	sqlState *sqlPaginationState

	Count int

	lastModel *model.Common
}

var _ listing.Pagination = (*sqlPagination)(nil)

// NewPagination creates a new SQL driven pagination from a listing.Request.
func NewPagination(req listing.Request) (*sqlPagination, error) {
	pagination := sqlPagination{}
	if err := listing.Init(req, &pagination.CommonState, &pagination.sqlState); err != nil {
		return nil, err
	}

	return &pagination, nil
}

// Query sets up limit, ordering and filtering on an orm.Query to retrieve the next page.
func (l sqlPagination) Query(db *gorm.DB) *gorm.DB {
	pageQuery := db
	if l.sqlState != nil {
		pageQuery.Where("(create_time, id) < (?, ?)", l.sqlState.Date, l.sqlState.ID)
	}
	if l.Knobs.ShowDeleted {
		pageQuery.Unscoped()
	}
	return pageQuery.Limit(l.Knobs.PageSize).Order("create_time DESC").Order("id DESC")
}

// Finish finalizes the pagination
func (l *sqlPagination) Finish() {
	if l.lastModel == nil {
		return
	}

	l.sqlState = &sqlPaginationState{
		Date: l.lastModel.CreateTime,
		ID:   l.lastModel.ID,
	}
}

func (l sqlPagination) ImplState() interface{} {
	return l.sqlState
}

func (l sqlPagination) HasNextPage() bool {
	return l.Count >= l.Knobs.PageSize
}

func (l *sqlPagination) ModelHook(model *model.Common) {
	l.Count++
	l.lastModel = model
}
