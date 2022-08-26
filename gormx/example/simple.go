package example

import (
	"reflect"

	"github.com/ahiho/gocandy/filter/adapter/sql"
	"github.com/ahiho/gocandy/gormx"
	"gorm.io/gorm"
)

// User represents an simple example database schema.
type User struct {
	ID   uint   `gorm:"id" filter:"=;>;>="`
	Name string `gorm:"name" filter:"#"`
}

// UserDAO is an example DAO for user data.
type UserDAO struct {
	db      *gorm.DB
	adaptor *sql.SQLAdaptor
}

// NewUserDAO returns a UserDAO.
func NewUserDAO(db *gorm.DB) (*UserDAO, error) {
	reflection := reflect.ValueOf(&User{})
	adaptor := sql.NewDefaultAdaptorFromStruct(reflection)
	return &UserDAO{
		db:      db,
		adaptor: adaptor,
	}, nil
}

// MakeQuery takes a goven query and performs it against the user database.
func (u *UserDAO) MakeQuery(pageSize int, pageToken string, filter string) ([]User, error) {
	res, _, err := gormx.PaginateTransform(u.db, gormx.PaginateRequest{
		PageSize:  pageSize,
		PageToken: pageToken,
		Filter:    filter,
		Adaptor:   u.adaptor,
	}, func(i User) (_ User, e error) {
		return i, nil
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *UserDAO) MakeQueryWithFilter(filter string) ([]User, error) {
	var user []User

	tx := gormx.Filter(u.db, filter, u.adaptor).Find(&user)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return user, nil
}
