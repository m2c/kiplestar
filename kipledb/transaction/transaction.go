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

	txUnits.db = txUnits.db.Begin()
	err = txUnits.db.Error
	if err != nil {
		return
	}

	for _, task := range txUnits.txUnits {
		if runErr := task.Run(txUnits.db); runErr != nil {
			// err will bubble up，just handle and rollback in outermost layer
			txUnits.db = txUnits.db.Rollback()
			if rollBackErr := txUnits.db.Error; rollBackErr != nil {
				slog.Infof("Rollback Failed: %s", rollBackErr.Error())
				err = rollBackErr
				return
			}
			err = runErr
			return
		}
	}

	txUnits.db = txUnits.db.Commit()
	if commitErr := txUnits.db.Error; commitErr != nil {
		err = commitErr
		return
	}

	return nil
}
