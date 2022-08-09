package example

import (
	"context"
	"reflect"

	"github.com/ahiho/filter/adapter/sql"
	"gorm.io/gorm"
)

// User represents an simple example database schema.
type User struct {
	ID   uint   `gorm:"id" filter:"=;>;>="`
	Name string `gorm:"name" filter:"#"`
}

// UserDAO is an example DAO for user data.
type UserDAO struct {
	db           *gorm.DB
	queryAdaptor *sql.SQLAdaptor
}

// NewUserDAO returns a UserDAO.
func NewUserDAO(db *gorm.DB) (*UserDAO, error) {
	reflection := reflect.ValueOf(&User{})
	adaptor := sql.NewDefaultAdaptorFromStruct(reflection)
	return &UserDAO{
		db:           db,
		queryAdaptor: adaptor,
	}, nil
}

// CreateUser commits the provided user to the database.
func (u *UserDAO) CreateUser(user *User) error {
	ctx := context.Background()
	tx := u.db.Begin().WithContext(ctx)
	err := tx.Create(user).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// MakeQuery takes a goven query and performs it against the user database.
func (u *UserDAO) MakeQuery(q string) ([]User, error) {
	var users []User
	ctx := context.Background()
	query := u.db.WithContext(ctx)
	queryResp, err := u.queryAdaptor.Parse(q)
	if err != nil {
		return nil, err
	}
	query = query.Model(User{}).Where(queryResp.Raw, sql.StringSliceToInterfaceSlice(queryResp.Values)...)
	err = query.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// func main() {
// 	var users []User
// 	dsn := "root:root@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"
// 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		panic("")
// 	}

// 	dao, err := NewUserDAO(db)
// 	ctx := context.Background()
// 	query := dao.db.WithContext(ctx)

// 	queryResp, err := dao.queryAdaptor.Parse("name#\"(duckhue02, duckhue10)\"")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	query = query.Model(User{}).Where(queryResp.Raw, sql_adaptor.StringSliceToInterfaceSlice(queryResp.Values)...)
// 	err = query.Find(&users).Error
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	fmt.Println(users)

// }
