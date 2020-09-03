package task

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"runtime/debug"
	slog "github.com/m2c/kiplestar/commons/log"
)

type Task func(db *gorm.DB) error

func (task Task) Run(db *gorm.DB) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("panic: %s; calltrace:%s", fmt.Sprint(e), string(debug.Stack()))
		}
	}()
	return task(db)
}

func NewTasks(db *gorm.DB) *Tasks {
	return &Tasks{
		db: db,
	}
}

type Tasks struct {
	db    *gorm.DB
	tasks []Task
}

func (tasks Tasks) With(task ...Task) *Tasks {
	tasks.tasks = append(tasks.tasks, task...)
	return &tasks
}

func (tasks *Tasks) Do() (err error) {
	if len(tasks.tasks) == 0 {
		return nil
	}

	tasks.db = tasks.db.Begin()
	err = tasks.db.Error
	if err != nil {
		return err
	}

	for _, task := range tasks.tasks {
		if runErr := task.Run(tasks.db); runErr != nil {
			// err will bubble up，just handle and rollback in outermost layer
			slog.Errorf("SQL FAILED: %s", runErr.Error())
			tasks.db = tasks.db.Rollback()
			if rollBackErr := tasks.db.Error; rollBackErr != nil {
				slog.Errorf("ROLLBACK FAILED: %s", rollBackErr.Error())
				err = rollBackErr
				return
			}
			return runErr
		}
	}

	tasks.db = tasks.db.Commit()
	if commitErr := tasks.db.Error; commitErr != nil {
		slog.Errorf("TRANSACTION COMMIT FAILED: %s", commitErr.Error())
		return commitErr
	}

	return nil
}
