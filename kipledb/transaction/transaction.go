package transaction

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"runtime/debug"
	slog "github.com/m2c/kiplestar/commons/log"
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

func (txUnits TxUnits) With(task ...TxUnit) *TxUnits {
	txUnits.txUnits = append(txUnits.txUnits, task...)
	return &txUnits
}

func (txUnits *TxUnits) Do() (err error) {
	if len(txUnits.txUnits) == 0 {
		return nil
	}

	txUnits.db = txUnits.db.Begin()
	err = txUnits.db.Error
	if err != nil {
		return err
	}

	for _, task := range txUnits.txUnits {
		if runErr := task.Run(txUnits.db); runErr != nil {
			// err will bubble up，just handle and rollback in outermost layer
			slog.Infof("SQL Run Failed: %s", runErr.Error())
			txUnits.db = txUnits.db.Rollback()
			if rollBackErr := txUnits.db.Error; rollBackErr != nil {
				slog.Infof("Rollback Failed: %s", rollBackErr.Error())
				err = rollBackErr
				return
			}
			return runErr
		}
	}

	txUnits.db = txUnits.db.Commit()
	if commitErr := txUnits.db.Error; commitErr != nil {
		slog.Infof("Transaction Commit Failed: %s", commitErr.Error())
		return commitErr
	}

	return nil
}
