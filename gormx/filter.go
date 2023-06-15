package gormx

import (
	filter "github.com/ahiho/gocandy/filter/adapter/sql"
	"gorm.io/gorm"
)

func Filter(db *gorm.DB, filterReq string, adaptor *filter.SQLAdaptor) (*gorm.DB, error) {
	queryResp := &filter.SQLResponse{}
	var err error
	if filterReq != "" {
		queryResp, err = adaptor.Parse(filterReq)
		if err != nil {
			return nil, err
		}
	}

	return db.Where(queryResp.Raw, filter.StringSliceToInterfaceSlice(queryResp.Values)...), nil
}
