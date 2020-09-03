package task

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"testing"
)

// db structure
type Player struct {
	ID   int64  `gorm:"column:id;AUTO_INCREMENT"`
	Name string `gorm:"column:name"`
}

func (Player) TableName() string {
	return "player"
}

func TestNewTasks(t *testing.T) {
	// init db connection
	// db config need to be edit
	db, err := gorm.Open("mysql", "root:xxx@tcp(127.0.0.1:3306)/xxx?charset=utf8")
	if err != nil {
		t.Fatal("连接数据库失败")
	}
	defer db.Close()
	db.SingularTable(true)

	// example start
	name := "test02"
	testFunc1 := func(db *gorm.DB) error {
		player := Player{
			Name: name,
		}
		err := db.Create(&player).Error
		if err != nil {
			return err
		}
		return nil
	}

	testFunc2 := func(db *gorm.DB) error {
		pp := Player{}

		err := db.Where("name=?", name).First(&pp).Error
		if err != nil {
			fmt.Printf("record not found")
			return err
		}

		// 数据联动展示，此func依赖上一个func的执行结果
		// Data linkage display, this func depends on the execution result of the previous func
		t.Logf("get %s: %v", name, pp)
		player := Player{
			Name: fmt.Sprintf("repeat_%s", pp.Name),
		}
		err = db.Create(&player).Error
		if err != nil {
			return err
		}
		return nil
	}
	tx := NewTasks(db)
	if err := tx.With(testFunc1, testFunc2).Do(); err != nil {
		t.Fatal(err)
	}
	t.Log("success")
}
