package kipledb

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	slog "github.com/m2c/kiplestar/commons/log"
	server_config "github.com/m2c/kiplestar/config"
	"time"
)

type KipleDB struct {
	db   *gorm.DB
	name string //db name
}

func (slf *KipleDB) DB() *gorm.DB {
	return slf.db
}
func (slf *KipleDB) Name() string {
	return slf.name
}

func (slf *KipleDB) StartDb(config server_config.DataBaseConfig) error {
	if slf.db != nil {
		return errors.New("Db already open")
	}
	slf.name = config.DbName
	driver := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=Local",
		config.User,
		config.Pwd,
		config.Host,
		config.Port,
		config.DataBase)
	var err error
	slf.db, err = gorm.Open("mysql", driver)

	if err != nil {
		slog.Infof("conn Db  error %s", err)
		return err
	}
	slog.Infof("conn Db opened Host %s", config.Host)
	slf.db.DB().SetMaxIdleConns(config.MaxIdleCons)
	slf.db.DB().SetMaxOpenConns(config.MaxOpenCons)
	slf.db.DB().SetConnMaxLifetime(config.MaxLifeTime)
	slf.db.SingularTable(true)
	if server_config.SC.SConfigure.Profile != "prod" {
		slf.db.LogMode(true)
	}

	slf.db.SetLogger(&slog.Slog)
	//slf.db.Callback().Create().Remove("gorm:create")
	//slf.db.Callback().Create().Remove("gorm:update")
	slf.db.Callback().Create().Before("gorm:create").Register("create", func(scope *gorm.Scope) {
		scope.SetColumn("create_time", time.Now())
		scope.SetColumn("update_time", time.Now())
	})
	slf.db.Callback().Update().Before("gorm:update").Register("update", func(scope *gorm.Scope) {
		scope.SetColumn("update_time", time.Now())
	})

	return nil
}

func (slf *KipleDB) StopDb() error {
	if slf.db != nil {
		err := slf.db.Close()
		if err != nil {
			slf.db = nil
		}
		return err
	}
	return errors.New("Db is nil")
}

func (slf *KipleDB) Tx(f func(db *gorm.DB) error) error {
	tx := slf.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}
	if e := f(tx); e != nil {
		tx.Rollback()
		return e
	}
	if e1 := tx.Commit().Error; e1 != nil {
		return e1
	}
	return nil
}
func (slf *KipleDB) ExecuteSql(f func(db *gorm.DB) (interface{}, error)) (interface{}, error) {
	result, ok := f(slf.db)
	if ok != nil {
		return result, ok
	}
	return result, nil
}
