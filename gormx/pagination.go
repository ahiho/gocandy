package gormx

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type PaginateRequest struct {
	PageSize int

	// For now we only support offset
	// Example: offset:100 (b2Zmc2V0OjEwMA)
	PageToken string
}

const (
	DefaultPageSize = 10

	offsetPrefix = "offset:"
	empty        = ""
)

var (
	ErrInvalidPageToken = errors.New("invalid page token")
)

func Paginate[T any](db *gorm.DB, pr PaginateRequest) ([]T, string, error) {
	pageSize, offset, err := pr.getQueryParams()
	if err != nil {
		return nil, empty, err
	}

	rs := []T{}
	tx := db.Offset(offset).Limit(pageSize).Find(&rs)
	if tx.Error != nil {
		return nil, empty, tx.Error
	}
	return rs, nextToken(offset, pageSize, len(rs)), nil
}

func PaginateTransform[T any, R any](db *gorm.DB, pr PaginateRequest, fn func(i T) (o R, e error)) ([]R, string, error) {
	pageSize, offset, err := pr.getQueryParams()
	if err != nil {
		return nil, empty, err
	}

	rs := []T{}
	tx := db.Offset(offset).Limit(pageSize).Find(&rs)
	if tx.Error != nil {
		return nil, empty, tx.Error
	}
	r := []R{}
	for _, i := range rs {
		o, e := fn(i)
		if e != nil {
			return nil, empty, e
		}
		r = append(r, o)
	}

	return r, nextToken(offset, pageSize, len(rs)), nil
}

func (p PaginateRequest) getQueryParams() (offset int, pageSize int, err error) {
	pageSize = p.PageSize
	if p.PageSize < 1 {
		pageSize = DefaultPageSize
	}
	offset = 0
	if p.PageToken == "" {
		return pageSize, 0, nil
	}
	pt, e := base64.URLEncoding.DecodeString(p.PageToken)
	if e != nil {
		return 0, 0, e
	}
	token := string(pt)
	if !strings.HasPrefix(token, offsetPrefix) {
		return 0, 0, ErrInvalidPageToken
	}
	ofs := token[len(offsetPrefix):]
	offset, e = strconv.Atoi(ofs)
	if e != nil || offset < 0 {
		return 0, 0, ErrInvalidPageToken
	}
	return pageSize, offset, nil
}

func nextToken(offset int, pageSize int, count int) string {
	if count < pageSize {
		return empty
	}
	return base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf("%v%v", offsetPrefix, offset+count)))
}
