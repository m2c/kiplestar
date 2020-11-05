package transaction

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	slog "github.com/m2c/kiplestar/commons/log"
	"runtime/debug"
)

type TxUnit func(db *gorm.DB) error

func (tu TxUnit) Run(db *gorm.DB) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("panic: %s; calltrace:%s", fmt.Sprint(e), string(debug.Stack()))
		}
	}()
	return tu(db)
}

func NewTxUnits(db *gorm.DB) *TxUnits {
	return &TxUnits{
		db: db,
	}
}

type TxUnits struct {
	db      *gorm.DB
	txUnits []TxUnit
}

func (txUnits *TxUnits) With(task ...TxUnit) *TxUnits {
	txUnits.txUnits = append(txUnits.txUnits, task...)
	return txUnits
}

func (txUnits *TxUnits) Do() (err error) {
	if len(txUnits.txUnits) == 0 {
		return
	}

	handleDb := txUnits.db.Begin()

	for _, task := range txUnits.txUnits {
		if runErr := task.Run(txUnits.db); runErr != nil {
			// err will bubble up，just handle and rollback in outermost layer
			if rollBackErr := handleDb.Rollback().Error; rollBackErr != nil {
				slog.Infof("Rollback Failed: %s", rollBackErr.Error())
				err = rollBackErr
				return
			}
			err = runErr
			return
		}
	}

	if commitErr := handleDb.Commit().Error; commitErr != nil {
		err = commitErr
		return
	}

	return nil
}
