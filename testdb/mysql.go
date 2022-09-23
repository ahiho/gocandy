package testdb

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	mysqldb "github.com/go-sql-driver/mysql"
	objpool "github.com/jolestar/go-commons-pool/v2"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	p    = "root"
	host = "127.0.0.1"

	once      sync.Once
	dtPool    *dockertest.Pool
	container *dockertest.Resource
	port      string
	dtError   error

	defaultDB *gorm.DB
	dbIndex   int32
	dbPool    *objpool.ObjectPool
	mu        sync.Mutex
)

func GetMysqlDB() (*gorm.DB, error) {
	once.Do(func() {
		createMysqlDockerTestContainer()
	})

	if dtError != nil {
		return nil, dtError
	}

	mu.Lock()
	_ = container.Expire(300)
	defer mu.Unlock()

	db, e := dbPool.BorrowObject(context.Background())
	if e != nil {
		return nil, e
	}
	return db.(*gorm.DB), nil
}

func ReleaseMysqlDB(db *gorm.DB) {
	_ = dropAllMysqlTables(db)
	_ = dbPool.ReturnObject(context.Background(), db)
	if dbPool.GetNumActive() == 0 {
		// auto terminate container after 5 seconds
		_ = container.Expire(5)
	}
}

func createMysqlDockerTestContainer() {
	dtPool, dtError = dockertest.NewPool("")
	dtPool.MaxWait = time.Second * 25
	if dtError != nil {
		return
	}
	container, dtError = dtPool.RunWithOptions(&dockertest.RunOptions{
		Repository:   "mysql",
		Tag:          "8.0",
		ExposedPorts: []string{"3306"},
		Env: []string{
			fmt.Sprintf("MYSQL_ROOT_PASSWORD=%s", p),
		},
	}, func(hc *docker.HostConfig) {
		hc.AutoRemove = true
		hc.RestartPolicy = docker.NeverRestart()
	})
	if dtError != nil {
		return
	}
	port = container.GetPort("3306/tcp")
	dtError = dtPool.Retry(func() error {
		var err error
		dsn := fmt.Sprintf("root:%s@tcp(%s:%s)/information_schema", p, host, port)
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return err
		}
		defaultDB = db
		return nil
	})
	if dtError != nil {
		_ = container.Close()
		return
	}

	factory := objpool.NewPooledObjectFactorySimple(func(ctx context.Context) (interface{}, error) {
		idx := atomic.AddInt32(&dbIndex, 1)

		dbName := fmt.Sprintf("database_%v", idx)
		err := defaultDB.Exec("CREATE DATABASE IF NOT EXISTS " + dbName).Error
		if err != nil {
			return nil, err
		}
		dsn := fmt.Sprintf("root:%s@(%s:%s)/%s?parseTime=true", p, host, port, dbName)
		newDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, err
		}

		return newDB, nil
	})

	dbPool = objpool.NewObjectPoolWithDefaultConfig(context.Background(), factory)
}

func dropAllMysqlTables(db *gorm.DB) error {
	dialector := db.Dialector.(*mysql.Dialector)
	cfg, err := mysqldb.ParseDSN(dialector.DSN)
	if err != nil {
		return err
	}

	tables := []string{}
	err = db.Raw(
		"select CONCAT('`', table_schema, '`.`', table_name, '`') from information_schema.TABLES where TABLE_SCHEMA = ?", cfg.DBName,
	).Pluck("tbn", &tables).Error

	if err != nil {
		return err
	}
	if len(tables) == 0 {
		return nil
	}

	tx := db.Begin()
	tx.Exec(`SET FOREIGN_KEY_CHECKS = 0;`)
	tx.Exec("DROP TABLE " + strings.Join(tables, ", "))
	tx.Exec(`SET FOREIGN_KEY_CHECKS = 1;`)
	return tx.Commit().Error
}
