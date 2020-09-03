package transaction

import (
	"errors"
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

type Sport struct {
	ID       int64  `gorm:"column:id;AUTO_INCREMENT"`
	PlayerId int64  `gorm:"column:player_id"`
	Name     string `gorm:"column:name"`
}

func (Sport) TableName() string {
	return "sport"
}

func TestNewTxUnits(t *testing.T) {
	// init db connection
	// db config need to be edit
	db, err := gorm.Open("mysql", "root:xiaohu@tcp(127.0.0.1:3306)/the-last-bastion?charset=utf8")
	if err != nil {
		t.Fatal("DB Connection Failed")
	}
	defer db.Close()
	db.SingularTable(true)

	// example start
	var pid int64 // Data can be used by all tx funcs when data put out func. You can put structure Player/Sport also.
	testFunc1 := func(db *gorm.DB) error {
		player := Player{
			Name: "test name",
		}
		err := db.Create(&player).Error
		if err != nil {
			// You can change err before return
			return errors.New(fmt.Sprintf("Create Player error: %s", err.Error()))
		}
		pid = player.ID

		return nil
	}

	testFunc2 := func(db *gorm.DB) error {
		sp := Sport{
			Name:     "test sport name",
			PlayerId: pid, // pid is uesed here.
		}

		err = db.Create(&sp).Error
		if err != nil {
			return err
		}
		return nil
	}

	tu := NewTxUnits(db)
	// when not many funcs
	if err := tu.With(testFunc1, testFunc2).Do(); err != nil {
		t.Fatal(err)
	}

	// when many funcs
	/*
		tu = tu.With(testFunc1)
		tu = tu.With(testFunc2)
		if err := tu.Do(); err != nil {
			t.Fatal(err)
		}
	*/

	t.Log("success")
}
