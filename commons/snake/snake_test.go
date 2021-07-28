package snake

import (
	"fmt"
	slog "github.com/m2c/kiplestar/commons/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	_ "math/rand"
	"os"
	"runtime"
	"sync"
	"testing"
)

type Transaction struct {
	TransactionNo string
}

func (Transaction) TableName() string {
	return "transaction"
}

/*CREATE TABLE `transaction` (
`bigint` bigint(20) NOT NULL AUTO_INCREMENT,
`transaction_no` varchar(255) NOT NULL,
PRIMARY KEY (`bigint`),
UNIQUE KEY `transaction` (`transaction_no`)
) ENGINE=InnoDB AUTO_INCREMENT=1746697 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
*/
// 1000000 data 1000 coroutine, pressure test
func TestNewNode(t *testing.T) {
	dsn := "root:root@tcp(127.0.0.1:3306)/test1?charset=UTF8"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Print(err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(-1)
	slog.Info(" init do the performance ")
	var wg sync.WaitGroup
	wg.Add(10000)

	for i := 0; i < 10000; i++ {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					var stacktrace string
					for i := 1; ; i++ {
						_, f, l, got := runtime.Caller(i)
						if !got {
							break
						}
						stacktrace += fmt.Sprintf("%s:%d\n", f, l)
						slog.Error(stacktrace)
						os.Exit(0)
					}
				}
			}()

			ExecuteDB1(db, wg.Done)
		}()

	}
	wg.Wait()
}
func ExecuteDB1(db *gorm.DB, done func()) {
	//defer  done()
	j := 0
	defer func() {
		slog.Info("====================================")
		slog.Error(j)
		slog.Info("====================================")
		done()
	}()

	for j = 0; j < 100; j++ {
		generate := fmt.Sprintf("%s", GetSnokeNode().Generate().String())

		t := Transaction{
			TransactionNo: generate,
		}
		if err := db.Create(&t).Error; err != nil {
			slog.Info("===================================")
			slog.Info(generate)
			slog.Info(err.Error())
			slog.Info("===================================")
			slog.Info(fmt.Sprintf("%s", err.Error()))
			os.Exit(0)
		} else {
			slog.Info(" well done")
		}

	}
}
